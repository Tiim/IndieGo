package model

type GenericComment struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	ReplyTo   string `json:"replyTo"`
	FromEmail string `json:"-"`
	Timestamp string `json:"timestamp"`
	Page      string `json:"page"`
	Url       string `json:"url"`
	Content   string `json:"content"`
	Name      string `json:"name"`
}
