package postgres

import (
	"database/sql"
	"fmt"
	"github.com/shamank/user-segmentation-service/internal/entity"
)

type SegmentRepo struct {
	db *sql.DB
}

func NewSegmentRepo(db *sql.DB) *SegmentRepo {
	return &SegmentRepo{
		db: db,
	}
}

func (r *SegmentRepo) CreateSegment(slug string) (int, error) {

	query := "INSERT INTO segments(slug) VALUES ($1) RETURNING id"

	res := r.db.QueryRow(query, slug)

	var Id int

	if err := res.Scan(&Id); err != nil {
		return 0, err
	}

	return Id, nil
}

func (r *SegmentRepo) DeleteSegment(slug string) error {

	query := "UPDATE segments SET deleted_at = now() WHERE slug = $1 and deleted_at is null"

	_, err := r.db.Exec(query, slug)
	if err != nil {
		// TODO: сделать обработку ошибок
		return err
	}

	return nil
}

func (r *SegmentRepo) GetSegmentBySlug(slug string) (entity.Segment, error) {
	query := "select * from segments where slug = $1 and deleted_at is null"

	var segment entity.Segment

	if err := r.db.QueryRow(query, slug).Scan(
		&segment.Id,
		&segment.Slug,
		&segment.CreatedAt,
		&segment.DeletedAt,
	); err != nil {
		// TODO: реализовать обработку
		return entity.Segment{}, err
	}

	return segment, nil
}

func (r *SegmentRepo) GetSegmentsBySlug(slugs []string) ([]entity.Segment, error) {

	subQuery := "("

	args := make([]interface{}, len(slugs))

	for i, val := range slugs {
		args[i] = val

		subQuery += fmt.Sprintf("$%d", i+1)
		if i+1 != len(slugs) {
			subQuery += ", "
		} else {
			subQuery += ")"
		}
	}

	query := fmt.Sprintf("select * from segments where slug in %s and deleted_at is null", subQuery)

	segmentSlice := make([]entity.Segment, 0)

	rows, err := r.db.Query(query, args...)
	if err != nil {

		return nil, err
	}

	for rows.Next() {
		var segment entity.Segment
		if err := rows.Scan(
			&segment.Id,
			&segment.Slug,
			&segment.CreatedAt,
			&segment.DeletedAt); err != nil {
			return nil, err
		}
		segmentSlice = append(segmentSlice, segment)
	}

	return segmentSlice, nil
}
