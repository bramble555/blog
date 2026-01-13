package model

type CommentModel struct {
	MODEL
	Content         string `json:"content"`                  // 评论内容
	DiggCount       int64  `json:"digg_count"`               // 点赞量
	CommentCount    int64  `json:"comment_count"`            // 子评论量
	ParentCommentSN int64  `json:"parent_comment_sn,string"` // 父评论ID,为-1，表示为根评论
	ArticleSN       int64  `json:"article_sn,string"`        // 评论的文章ID
	UserSN          int64  `json:"user_sn,string"`           // 评论的用户ID
}
type ParamPostComment struct {
	ArticleSN       int64  `json:"article_sn,string" binding:"required"` // 文章ID
	ParentCommentSN int64  `json:"parent_comment_sn,string"`             // 父级评论ID
	Content         string `json:"content" binding:"required"`           // 评论内容
}
type ParamCommentList struct {
	ArticleSN int64 `form:"article_sn" binding:"required"` // 文章ID
}
type ResponseCommentList struct {
	MODEL
	Content         string                `json:"content"`
	ParentCommentSN int64                 `json:"parent_comment_SN"`
	ArticleSN       int64                 `json:"article_sn,string"`
	DiggCount       int64                 `json:"digg_count"`             // 点赞数
	CommentCount    int64                 `json:"comment_count"`          // 子评论数量
	SubComments     []ResponseCommentList `json:"sub_comments,omitempty"` // 子评论列表，嵌套结构
	*UserDetail
}
