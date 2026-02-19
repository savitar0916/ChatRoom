package domain

type Message struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}
