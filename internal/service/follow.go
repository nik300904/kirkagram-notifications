package service

import "log/slog"

type FollowService interface {
	GetUsernameByID(followerID int) (string, error)
}

type Follow struct {
	client FollowService
	log    *slog.Logger
}

func NewFollowService(client FollowService, log *slog.Logger) *Follow {
	return &Follow{
		client: client,
		log:    log,
	}
}

func (f *Follow) GetUsernameByID(followerID int) (string, error) {
	return f.client.GetUsernameByID(followerID)
}
