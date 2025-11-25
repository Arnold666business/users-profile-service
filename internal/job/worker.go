package job

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"
	"users-profile-service/internal/models"
	"users-profile-service/internal/user/unblock"

	"go.uber.org/zap"
)

type CheckUsersForUnblockJob struct {
	logger               *zap.SugaredLogger
	unblockUserProcessor *unblock.UnblockUserProcessor
}

func Build(logger *zap.SugaredLogger,
	unblockUserProcessor *unblock.UnblockUserProcessor) *CheckUsersForUnblockJob {
	return &CheckUsersForUnblockJob{
		logger:               logger,
		unblockUserProcessor: unblockUserProcessor,
	}
}

func (job *CheckUsersForUnblockJob) Run(ctx context.Context) {
	interval, err := strconv.Atoi(os.Getenv("UNBLOCK_JOB_INTERVAL_SECOND"))
	if err != nil {
		job.logger.Fatal("unblock job interval parse failed", zap.Error(err))
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	workerPool := 3
	limit := 1000
	jobQueue := make(chan *models.UserBlockStatus, limit)

	var wg sync.WaitGroup
	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go job.worker(ctx, &wg, jobQueue)
	}

	for {
		select {
		case <-ticker.C:
			unblockStatusEntities, err := job.unblockUserProcessor.FindForUnblock(ctx, limit)
			if err != nil {
				job.logger.Errorw("unblock job fetch failed", zap.Error(err))
				continue
			}

			for _, unblockStatus := range unblockStatusEntities {
				select {
				case jobQueue <- unblockStatus:
				case <-ctx.Done():
					return
				}
			}

		case <-ctx.Done():
			close(jobQueue)
			wg.Wait()
			return
		}
	}
}

func (job *CheckUsersForUnblockJob) worker(ctx context.Context, wg *sync.WaitGroup, jobQueue chan *models.UserBlockStatus) {
	defer wg.Done()
	for {
		select {
		case blockUserStatus, ok := <-jobQueue:
			if !ok {
				return
			}
			err := job.unblockUserProcessor.Process(ctx, blockUserStatus)
			if err != nil {
				job.logger.Errorw("unblock job processing failed", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}
