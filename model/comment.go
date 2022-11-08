package model

type GenericComment struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	ReplyTo   string `json:"replyTo"`
	FromEmail string `json:"fromEmail"`
	Timestamp string `json:"timestamp"`
	Page      string `json:"page"`
	Content   string `json:"content"`
	Name      string `json:"name"`
}
