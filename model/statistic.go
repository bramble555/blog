package model

type DataSumResponse struct {
	UserCount      int64 `json:"user_count"`
	ArticleCount   int64 `json:"article_count"`
	MessageCount   int64 `json:"message_count"`
	ChatGroupCount int64 `json:"chat_group_count"`
}
