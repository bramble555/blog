package model

type CommentModel struct {
	*MODEL
	Content         string         `json:"content"`                         // 评论内容
	DiggCount       uint           `json:"digg_count"`                      // 点赞数
	CommentCount    uint           `json:"comment_count"`                   // 子评论数量
	ParentCommentID int            `json:"parent_comment_id"`               // 父级评论ID
	ArticleID       uint           `json:"article_id,string"`               // 文章ID
	UserID          uint           `json:"user_id,string"`                  // 评论的用户ID
}
type ParamPostComment struct {
	ArticleID       uint   `json:"article_id,string" binding:"required"` // 文章ID
	ParentCommentID int    `json:"parent_comment_id,string"`             // 父级评论ID
	Content         string `json:"content" binding:"required"`           // 评论内容
}
type ParamCommentList struct {
	ArticleID uint `form:"article_id" binding:"required"` // 文章ID
}
type ResponseCommentList struct {
	*MODEL
	Content         string                `json:"content"`
	ParentCommentID int                   `json:"parent_comment_id"`
	ArticleID       uint                  `json:"article_id,string"`
	DiggCount       uint                  `json:"digg_count"`             // 点赞数
	CommentCount    uint                  `json:"comment_count"`          // 子评论数量
	SubComments     []ResponseCommentList `json:"sub_comments,omitempty"` // 子评论列表，嵌套结构
	*UserDetail
}
