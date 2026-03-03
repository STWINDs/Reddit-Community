package redis

// redis key

//redis key 命名时切割命名空间 便于拆分
const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"
	KeyPostScoreZSet       = "post:score"
	KeyPostVotedZSetPrefix = "post:voted:" // zset:记录用户及投票类型；params为post id
	KeyCommunitySetPF      = "community:"  //set; 保存每个分区下帖子的id
)

//加前缀
func getRedisKey(Key string) string {
	return KeyPrefix + Key
}
