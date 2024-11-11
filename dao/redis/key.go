package redis

// redis key, 注意使用命名，方便查询和拆分

const (
	KeyPrefix          = "gvb:"
	KeyZSetArticleDig  = "article:dig" // 帖子及发帖时间
	KeyZSetCommentDigg = "article:comment"
	KeyZSetPostVotedPF = "post:voted:" // 记录用户及投票类型；参数是post_id
	KeySetToken        = "token"
)

func getKeyName(key string) string {
	return KeyPrefix + key
}
