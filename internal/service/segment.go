package service

import (
	"github.com/shamank/user-segmentation-service/internal/entity"
	"github.com/shamank/user-segmentation-service/internal/repository"
	"log/slog"
)

type SegmentService struct {
	repo   repository.Segment
	logger *slog.Logger
}

func NewSegmentService(repo repository.Segment, log *slog.Logger) *SegmentService {
	return &SegmentService{
		repo:   repo,
		logger: log,
	}
}

func (s *SegmentService) CreateSegment(slug string) (int, error) {

	Id, err := s.repo.CreateSegment(slug)

	return Id, err
}

func (s *SegmentService) DeleteSegment(slug string) error {

	return s.repo.DeleteSegment(slug)
}

func (s *SegmentService) GetSegmentBySlug(slug string) (entity.Segment, error) {

	return s.repo.GetSegmentBySlug(slug)
}

func (s *SegmentService) GetSegmentsBySlug(slugs []string) ([]entity.Segment, error) {
	return s.repo.GetSegmentsBySlug(slugs)
}
