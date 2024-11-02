package model

type CommentModel struct {
	*MODEL
	Content            string          `json:"content"`           // 评论内容
	DiggCount          uint            `json:"digg_count"`        // 点赞数
	CommentCount       uint            `json:"comment_count"`     // 子评论数量
	ParentCommentID    *int            `json:"parent_comment_id"` // 父级评论ID
	SubComments        []*CommentModel `json:"sub_comments"`      // 子评论列表
	ParentCommentModel *CommentModel   `json:"parent_comment"`    // 父级评论
	ArticleID          uint            `json:"article_id,string"` // 文章ID
	UserID             uint            `json:"user_id,string"`    // 评论的用户ID
}
