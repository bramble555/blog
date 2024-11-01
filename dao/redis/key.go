package redis

// redis key, 注意使用命名，方便查询和拆分

const (
	KeyPrefix          = "gvb:"
	KeyZSetPostTime    = "post:time"   // 帖子及发帖时间
	KeyZSetPostScore   = "post:score"  // 帖子及投票分数
	KeyZSetPostVotedPF = "post:voted:" // 记录用户及投票类型；参数是post_id
	KeySetCommuntiyPF  = "community:"  // 记录社区的post列表；参数是communtiy_id
	KeySetToken        = "token"
)

func getKeyName(key string) string {
	return KeyPrefix + key
}
