package model

type Comment struct {
	Id        string `json:"id"`
	ReplyTo   string `json:"reply_to"`
	Timestamp string `json:"timestamp"`
	Page      string `json:"page"`
	Content   string `json:"content"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Notify    bool   `json:"notify"`
}
