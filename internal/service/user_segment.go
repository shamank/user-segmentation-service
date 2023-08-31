package service

import (
	"errors"
	"fmt"
	"github.com/shamank/user-segmentation-service/internal/repository"
	"time"
)

const (
	grantOperation  = "grant"
	revokeOperation = "revoke"
)

type UserSegmentService struct {
	userService            User
	segmentService         Segment
	userSegmentHistoryRepo repository.UserSegmentHistory
	csvManager             CSVer
}

func NewUserSegmentService(userService User, segmentService Segment, userSegmentHistoryRepo repository.UserSegmentHistory, CSVManager CSVer) *UserSegmentService {

	return &UserSegmentService{
		userService:            userService,
		segmentService:         segmentService,
		userSegmentHistoryRepo: userSegmentHistoryRepo,
		csvManager:             CSVManager,
	}
}

func (s *UserSegmentService) AddUserToSegments(userId int, slugs []string, expire *time.Duration) error {

	segmentsIds := make([]int, 0)

	for _, slug := range slugs {
		slug := slug
		segment, err := s.segmentService.GetSegmentBySlug(slug)
		if err != nil {
			// сделать обработку ошибок
			return err
		}
		segmentsIds = append(segmentsIds, segment.Id)
	}

	if len(segmentsIds) != len(slugs) {
		return errors.New("no match")
	}

	var endAt *time.Time

	if expire != nil {
		tmp := time.Now().Add(*expire * time.Second)
		endAt = &tmp
	}

	if err := s.userService.AddUserToSegments(userId, segmentsIds, endAt); err != nil {
		return err
	}

	if err := s.userSegmentHistoryRepo.SaveHistory(userId, segmentsIds, grantOperation); err != nil {
		return err
	}

	return nil

}

func (s *UserSegmentService) RemoveUserFromSegments(userId int, slugs []string) error {
	segmentsIds := make([]int, 0)

	for _, slug := range slugs {
		slug := slug
		segment, err := s.segmentService.GetSegmentBySlug(slug)
		if err != nil {
			// сделать обработку ошибок
			return err
		}
		segmentsIds = append(segmentsIds, segment.Id)
	}

	if len(segmentsIds) != len(slugs) {
		return errors.New("no match")
	}

	if err := s.userService.RemoveUserFromSegments(userId, segmentsIds); err != nil {
		return err
	}

	if err := s.userSegmentHistoryRepo.SaveHistory(userId, segmentsIds, revokeOperation); err != nil {
		return err
	}

	return nil
}

func (s *UserSegmentService) GetUserSegmentHistory(userID int, startDate time.Time, endDate time.Time) (string, error) {

	userHistory, err := s.userSegmentHistoryRepo.GetHistoryByUserID(userID, startDate, endDate)
	if err != nil {
		return "", err
	}

	userHistoryAsIntefaces := make([][]interface{}, len(userHistory))

	for i, history := range userHistory {
		i, history := i, history

		userHistoryAsIntefaces[i] = []interface{}{
			history.UserID,
			history.SegmentID,
			history.Operation,
			history.CreatedAt,
		}
	}

	file, err := s.csvManager.CreateFile(userHistoryAsIntefaces)
	if err != nil {
		return "", err
	}

	return file, nil
}

func (s *UserSegmentService) SetSegmentToRandomUsers(segmentSlug string, percentage int) error {

	users, err := s.userService.GetRandomUsers(percentage)
	if err != nil {
		return err
	}

	fmt.Println("USERS: ", users)

	for _, user := range users {
		fmt.Println(user.Id)
		err := s.AddUserToSegments(user.Id, []string{segmentSlug}, nil)
		if err != nil {
			return err
		}
	}

	return nil

}
