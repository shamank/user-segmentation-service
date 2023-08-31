package service

import (
	"github.com/shamank/user-segmentation-service/internal/entity"
	"github.com/shamank/user-segmentation-service/internal/repository"
	"log/slog"
	"time"
)

type UserService struct {
	repo   repository.User
	logger *slog.Logger
}

func NewUserService(repo repository.User, logger *slog.Logger) *UserService {

	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) GetRandomUsers(assignPercentage int) ([]entity.User, error) {

	return s.repo.GetRandomUsers(assignPercentage)
}

func (s *UserService) CreateUser(username string) (int, error) {

	return s.repo.CreateUser(username)
}

func (s *UserService) GetUserSegments(userId int) ([]string, error) {

	return s.repo.GetUserSegments(userId)

}

func (s *UserService) AddUserToSegments(userId int, segmentsIds []int, endAt *time.Time) error {
	return s.repo.AddUserToSegments(userId, segmentsIds, endAt)
}

func (s *UserService) RemoveUserFromSegments(userId int, segmentsIds []int) error {
	return s.repo.RemoveUserFromSegments(userId, segmentsIds)
}
