package service

import "log/slog"

type PostService interface {
	GetUsernameByID(userID int) (string, error)
	GetUsersIDToSendByID(userID int) ([]int, error)
}

type Post struct {
	client PostService
	log    *slog.Logger
}

func NewPostService(client PostService, log *slog.Logger) *Post {
	return &Post{client: client, log: log}
}

func (p *Post) GetUsernameByID(userID int) (string, error) {
	return p.client.GetUsernameByID(userID)
}

func (p *Post) GetUsersIDToSendByID(userID int) ([]int, error) {
	return p.client.GetUsersIDToSendByID(userID)
}
