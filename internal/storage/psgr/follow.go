package psgr

import (
	"database/sql"
	"errors"
	"fmt"
	"kirkagram-notification/internal/storage"
)

type FollowStorage struct {
	db *sql.DB
}

func NewFollowStorage(db *sql.DB) *FollowStorage {
	return &FollowStorage{db: db}
}

func (f *FollowStorage) GetUsernameByID(followerID int) (string, error) {
	const op = "storage.psgr.follow.GetUsernameByID"

	var username string

	err := f.db.QueryRow(`
		SELECT DISTINCT u.username
		FROM "users" u
		WHERE u.id = $1`,
		followerID,
	).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrFollowerNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return username, nil
}
