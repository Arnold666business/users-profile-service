package users

import (
	"context"
	"errors"
	"users-profile-service/api/generation/users/v1/service"
	"users-profile-service/internal/models"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/user/management"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	logger                *zap.SugaredLogger
	userManagementService *management.UserManager
	service.UnimplementedUserProfileServiceServer
}

func Build(Logger *zap.SugaredLogger, userManagementService *management.UserManager) *UserService {
	return &UserService{logger: Logger, userManagementService: userManagementService}
}

func (u *UserService) GetUserById(ctx context.Context, in *service.GetUserRequest) (*service.GetUserResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userAgg, err := u.userManagementService.GetUser(ctx, in.GetId())
	if err != nil {
		u.logger.Error("failed to get user",
			zap.Int64("user_id", in.GetId()),
			zap.Error(err),
		)
		if errors.Is(err, repository.NotFoundUserError) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return mapUserAggregateToProto(userAgg), nil
}

func mapUserAggregateToProto(u *models.UserProfile) *service.GetUserResponse {
	return &service.GetUserResponse{
		Id:               u.Id,
		Login:            u.Login,
		Email:            u.Email,
		EmailAccess:      u.EmailAccess,
		Role:             int32(u.Role),
		IsDeleted:        u.IsDeleted,
		DeletedAt:        timestamppb.New(u.DeletedAt),
		BlockTypeId:      int32(u.BlockTypeId),
		ForeverFlag:      u.ForeverFlag,
		UnblockDate:      timestamppb.New(u.UnBlockDate),
		BlockTitle:       u.BlockTitle,
		BlockDescription: u.BlockDescription,
	}
}
