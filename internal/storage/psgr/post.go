package psgr

import (
	"database/sql"
	"errors"
	"fmt"
	"kirkagram-notification/internal/storage"
)

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{db: db}
}

func (l *PostStorage) GetUsernameByID(userID int) (string, error) {
	const op = "storage.psgr.post.GetUsernameByID"

	var username string

	err := l.db.QueryRow(`
		SELECT DISTINCT u.username 
		FROM "post" p 
		JOIN "users" u ON p.user_id = u.id 
		WHERE p.user_id = $1
		`,
		userID,
	).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrLikeNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return username, nil
}

func (l *PostStorage) GetUsersIDToSendByID(userID int) ([]int, error) {
	const op = "storage.psgr.post.GetUserIDToSendByID"

	rows, err := l.db.Query(
		`
		SELECT DISTINCT f.follower_id
		FROM "post" p
		JOIN "follow" f ON p.user_id = f.following_id
		WHERE p.user_id = $1
		`,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]int, 0), storage.ErrLikeNotFound
		}

		return make([]int, 0), fmt.Errorf("%s: %w", op, err)
	}

	var userIDList []int

	for rows.Next() {
		var ID int

		if err = rows.Scan(&ID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		userIDList = append(userIDList, ID)
	}

	if len(userIDList) == 0 {
		return nil, storage.ErrNoFollowers
	}

	return userIDList, nil
}
