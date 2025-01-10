package service

import "log/slog"

type LikeService interface {
	GetUsernameByUserID(userID int) (string, error)
	GetUserIDToSendByPostID(userID int) (int, error)
}

type Like struct {
	client LikeService
	log    *slog.Logger
}

func NewLikeService(client LikeService, log *slog.Logger) *Like {
	return &Like{
		client: client,
		log:    log,
	}
}

func (l *Like) GetUsernameByUserID(userID int) (string, error) {
	return l.client.GetUsernameByUserID(userID)
}

func (l *Like) GetUserIDToSendByPostID(userID int) (int, error) {
	return l.client.GetUserIDToSendByPostID(userID)
}
