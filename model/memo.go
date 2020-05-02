package model

type Memo struct {
	MemoId  string `json:"memo_id"`
	UserId  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
