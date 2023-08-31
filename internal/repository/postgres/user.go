package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shamank/user-segmentation-service/internal/entity"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {

	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetRandomUsers(percentage int) ([]entity.User, error) {

	fmt.Println("percentage: ", percentage)

	query := `select id, username from users order by random() limit (select  round(count(id) * $1 / 100) from users)`

	rows, err := r.db.Query(query, percentage)
	if err != nil {
		return nil, err
	}

	fmt.Println("ROWS: ", rows)

	users := make([]entity.User, 0)

	for rows.Next() {
		var user entity.User

		if err := rows.Scan(
			&user.Id,
			&user.Username,
		); err != nil {
			return nil, err
		}
		fmt.Println("USER: ", user)
		users = append(users, user)
	}
	return users, nil

}

func (r *UserRepo) CreateUser(username string) (int, error) {

	query := "INSERT INTO users(username) VALUES ($1) RETURNING id"

	row := r.db.QueryRow(query, username)

	var Id int

	if err := row.Scan(&Id); err != nil {
		// TODO: сделать обработку ошибок
		return 0, err
	}

	return Id, nil
}

func (r *UserRepo) GetUserSegments(userId int) ([]string, error) {

	query := `SELECT s.slug FROM users u
    INNER JOIN user_segment us on u.id = us.user_id
    INNER JOIN segments s on us.segment_id = s.id
    WHERE u.id = $1 and s.deleted_at is null and us.end_at < now()`

	var segments []string

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var value string

		err := rows.Scan(&value)
		if err != nil {
			//  TODO: сделать обработку ошибок
			return nil, err
		}

		segments = append(segments, value)
	}

	return segments, nil
}

func (r *UserRepo) AddUserToSegments(userId int, segments []int, endAt *time.Time) error {

	segmentsIds, subQuery := makeQueryForSlug(segments, 2)

	query := fmt.Sprintf(`insert into user_segment(user_id, end_at, segment_id)
									select $1, $2, segments.id from segments
									    where id in (%s) and segments.deleted_at is null`, subQuery)

	result, err := r.db.Exec(query, append([]interface{}{userId, endAt}, segmentsIds...)...)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return errors.New("no rows affected!")
	}

	return nil
}

func (r *UserRepo) RemoveUserFromSegments(userId int, segments []int) error {

	segmentsIds, subQuery := makeQueryForSlug(segments, 1)

	query := fmt.Sprintf(`delete from user_segment where user_id = $1 and segment_id in (%s)`, subQuery)

	result, err := r.db.Exec(query, append([]interface{}{userId}, segmentsIds...)...)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return errors.New("no rows affected!")
	}

	return nil
}
