package repository

import (
	"database/sql"
	"github.com/shamank/user-segmentation-service/internal/entity"
	"github.com/shamank/user-segmentation-service/internal/repository/postgres"
	"time"
)

type Segment interface {
	CreateSegment(slug string) (int, error)
	DeleteSegment(slug string) error
	GetSegmentBySlug(slug string) (entity.Segment, error)
	GetSegmentsBySlug(slugs []string) ([]entity.Segment, error)

	//AddUser(userID int, slugs []string)
	//DeleteUser(userID int, slugs []string)
	//GetUser()
}

type User interface {
	CreateUser(username string) (int, error)
	GetRandomUsers(percentage int) ([]entity.User, error)

	GetUserSegments(userId int) ([]string, error)
	AddUserToSegments(userID int, segments []int, endAt *time.Time) error
	RemoveUserFromSegments(userId int, segments []int) error
}

type UserSegmentHistory interface {
	SaveHistory(userID int, segments []int, operation string) error
	GetHistoryByUserID(userID int, startDate time.Time, endDate time.Time) ([]entity.UserSegmentHistory, error)
}

type Repositories struct {
	db *sql.DB

	Segment            Segment
	User               User
	UserSegmentHistory UserSegmentHistory
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Segment:            postgres.NewSegmentRepo(db),
		User:               postgres.NewUserRepo(db),
		UserSegmentHistory: postgres.NewUserSegmentHistoryRepo(db),
	}
}
