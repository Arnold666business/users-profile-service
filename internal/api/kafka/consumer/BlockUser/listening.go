package BlockUser

import (
	"context"
	"encoding/json"
	"fmt"
	"users-profile-service/internal/user/block"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func (c *BlockUser) StartListening(ctx context.Context) {
	c.logger.Info("starting kafka consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context cancelled, stopping consumer")
			return

		default:
			if err := c.processNextMessage(ctx); err != nil {
				c.logger.Error("kafka read error",
					zap.Error(err),
				)
			}
		}
	}
}

func (c *BlockUser) processNextMessage(ctx context.Context) error {
	msg, err := c.conn.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	c.logger.Debug("message received",
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("key", string(msg.Key)),
		zap.Int("value_size", len(msg.Value)),
	)

	blockInfo, err := c.parseMessage(msg.Value)
	if err != nil {
		c.logger.Error("invalid message format, skipping",
			zap.Error(err),
			zap.ByteString("raw_value", msg.Value),
		)
		c.commitMessage(ctx, msg)
		return nil
	}
	if err := c.processor.Process(ctx, block.BlockRequest{
		UserId:      blockInfo.UserID,
		BlockTypeId: blockInfo.BlockTypeID,
		ForeverFlag: blockInfo.Forever,
	}); err != nil {
		c.logger.Error("failed to process block",
			zap.Error(err),
			zap.Any("block_info", blockInfo),
		)
		return fmt.Errorf("processing failed: %w", err)
	}

	c.commitMessage(ctx, msg)
	return nil
}

func (c *BlockUser) parseMessage(data []byte) (*BlockedInfo, error) {
	var info BlockedInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return &info, nil
}

func (c *BlockUser) commitMessage(ctx context.Context, msg kafka.Message) {
	if err := c.conn.CommitMessages(ctx, msg); err != nil {
		c.logger.Error("failed to commit message",
			zap.Error(err),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
		)
	} else {
		c.logger.Debug("message committed",
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
		)
	}
}

func (c *BlockUser) StopListening() error {
	return c.conn.Close()
}
