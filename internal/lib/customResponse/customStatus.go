package customResponse

type CustomStatus struct {
	Status int `json:"status"`
}

func NewStatus(status int) *CustomStatus {
	return &CustomStatus{Status: status}
}
