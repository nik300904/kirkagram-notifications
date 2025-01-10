package models

type Posts struct {
	UserID   int    `json:"user_id"`
	Caption  string `json:"caption"`
	ImageURL string `json:"image_url"`
}
