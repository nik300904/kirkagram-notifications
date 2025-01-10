package psgr

import (
	"database/sql"
	"errors"
	"fmt"
	"kirkagram-notification/internal/storage"
)

type LikeStorage struct {
	db *sql.DB
}

func NewLikeStorage(db *sql.DB) *LikeStorage {
	return &LikeStorage{db: db}
}

func (l *LikeStorage) GetUsernameByUserID(userID int) (string, error) {
	const op = "storage.psgr.like.GetUsernameByID"

	var username string

	err := l.db.QueryRow(`
		SELECT DISTINCT u.username
		FROM "like" l
		JOIN "users" u ON l.user_id = u.id
		WHERE l.user_id = $1`,
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

func (l *LikeStorage) GetUserIDToSendByPostID(postID int) (int, error) {
	const op = "storage.psgr.like.GetUserIDToSendByID"

	var userID int

	err := l.db.QueryRow(`
		SELECT DISTINCT u.id
		FROM "like" l
		JOIN "post" p ON l.post_id = p.id
		JOIN "users" u ON p.user_id = u.id
		WHERE l.post_id = $1
		`,
		postID,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrLikeNotFound
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
