package model

type Email string
type Notify bool

type Comment struct {
	Id                string `json:"id"`
	ReplyTo           string `json:"reply_to"`
	Timestamp         string `json:"timestamp"`
	Page              string `json:"page"`
	Content           string `json:"content"`
	Name              string `json:"name"`
	Email             Email  `json:"email"`
	Notify            Notify `json:"notify"`
	UnsubscribeSecret string `json:"-"`
}

func (Email) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}

func (Notify) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}
