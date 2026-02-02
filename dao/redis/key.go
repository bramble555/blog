package redis

// redis key, 注意使用命名，方便查询和拆分

const (
	KeyPrefix = "gvb:"

	KeySTrTokenLogoutPF = "token:logout:"

	// 首页最新帖子(只选择最新10个帖子)
	KeyZSetHomeLatestArticleSN = "home:latest:article:sn"

	// 文章浏览数
	KeyHashArticleLookCount = "article:look"
	// 文章点赞数
	KeyHashArticleDiggCount = "article:digg"

	// 评论点赞数
	KeyHashCommentDiggCount = "comment:digg"
)

func getKeyName(key string) string {
	return KeyPrefix + key
}
