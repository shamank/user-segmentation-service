package postgres

import (
	"database/sql"
	"fmt"
	"github.com/shamank/user-segmentation-service/internal/entity"
	"time"
)

type UserSegmentHistoryRepo struct {
	db *sql.DB
}

func NewUserSegmentHistoryRepo(db *sql.DB) *UserSegmentHistoryRepo {
	return &UserSegmentHistoryRepo{
		db: db,
	}
}

func (r *UserSegmentHistoryRepo) SaveHistory(userId int, segments []int, operation string) error {
	segmentsIds, subQuery := makeQueryForSlug(segments, 2)

	query := fmt.Sprintf(`insert into user_segment_history(user_id, operation, segment_id)
						select $1, $2, segments.id from segments
						where id in (%s) and segments.deleted_at is null`,
		subQuery)

	_, err := r.db.Exec(query, append([]interface{}{userId, operation}, segmentsIds...)...)

	return err
}

func (r *UserSegmentHistoryRepo) GetHistoryByUserID(userID int, startDate time.Time, endDate time.Time) ([]entity.UserSegmentHistory, error) {

	query := `select * from user_segment_history where user_id = $1 and created_at between $2 and $3`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]entity.UserSegmentHistory, 0)

	for rows.Next() {
		var history entity.UserSegmentHistory
		if err := rows.Scan(
			&history.ID,
			&history.UserID,
			&history.SegmentID,
			&history.Operation,
			&history.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, history)
	}

	return result, nil
}
