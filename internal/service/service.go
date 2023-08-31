package service

import (
	"github.com/shamank/user-segmentation-service/internal/entity"
	"github.com/shamank/user-segmentation-service/internal/repository"
	"log/slog"
	"time"
)

type Segment interface {
	CreateSegment(slug string) (int, error)
	DeleteSegment(slug string) error

	GetSegmentBySlug(slug string) (entity.Segment, error)
	GetSegmentsBySlug(slugs []string) ([]entity.Segment, error)
}

type User interface {
	CreateUser(username string) (int, error)

	AddUserToSegments(userId int, segmentsIds []int, endAt *time.Time) error
	RemoveUserFromSegments(userId int, segmentsIds []int) error

	GetUserSegments(userId int) ([]string, error)
}

type UserSegment interface {
	AddUserToSegments(userId int, slugs []string, expire *time.Duration) error
	RemoveUserFromSegments(userId int, slugs []string) error
	GetUserSegmentHistory(userID int, startDate time.Time, endDate time.Time) (string, error)
}

type Services struct {
	Segment     Segment
	User        User
	UserSegment UserSegment
}

type CSVer interface {
	CreateFile(date [][]interface{}) (string, error)
}

type Dependencies struct {
	Logger     *slog.Logger
	Repo       *repository.Repositories
	CSVManager CSVer
}

func NewService(deps Dependencies) *Services {

	segmentService := NewSegmentService(deps.Repo.Segment, deps.Logger)
	userService := NewUserService(deps.Repo.User, deps.Logger)
	userSegmentService := NewUserSegmentService(
		userService,
		segmentService,
		deps.Repo.UserSegmentHistory,
		deps.CSVManager)

	return &Services{
		Segment:     segmentService,
		User:        userService,
		UserSegment: userSegmentService,
	}
}
