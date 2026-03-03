package redis

import (
	"BLUEBELL/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// AddPostIDToRedis 将新创建的贴子ID添加到Redis的Sorted Set中
// 按发布时间和分数分别存储到两个ZSet中
func AddPostIDToRedis(postID, communityID int64) error {
	// 使用当前时间戳作为分数（精确到秒）
	now := float64(time.Now().Unix())
	postIDStr := strconv.FormatInt(postID, 10)

	pipeline := client.Pipeline()

	// 添加到按时间排序的ZSet
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: postIDStr,
	})

	// 添加到按分数排序的ZSet（初始分数为0）
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postIDStr,
	})

	// 添加到社区ZSet
	cKey := getRedisKey(KeyCommunitySetPF + strconv.FormatInt(communityID, 10))
	pipeline.SAdd(cKey, postIDStr)

	_, err := pipeline.Exec()
	return err
}

func GetPostIDsInOrder(p *models.PostListParam) ([]string, error) {
	//从发布时间ZSet/分数ZSet来获取id
	//根据用户请求的redis参数确认查询的redis Key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.Orderscore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZREVRANGE查询
	return client.ZRevRange(key, start, end).Result()

}

func GetCommunityPostIDsInOrder(p *models.CommunityPostListParam) ([]string, error) {
	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset按之前的逻辑取数据

	// 1. 获取社区帖子的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.FormatInt(p.CommunityID, 10))

	// 2. 获取排序的key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.Orderscore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 3. 生成临时 zset 的 key
	tempKey := orderKey + ":" + strconv.FormatInt(p.CommunityID, 10)

	// 4. 检查临时 zset 是否已存在（或者为了实时性，直接生成，并设置过期时间）
	if client.Exists(tempKey).Val() == 0 {
		// 不存在，则生成
		pipeline := client.Pipeline()
		pipeline.ZInterStore(tempKey, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipeline.Expire(tempKey, 60*time.Second) // 缓存 60 秒
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	// 5. 从临时 zset 中获取数据
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(tempKey, start, end).Result()
}

// 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// data = make([]int64, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	// 	// 查找key中分数是1的元素的数量->统计每篇帖子赞成票的数量
	// 	v := client.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }
	// return data, err

	//用pipeline  一次发送多条命令 减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return data, err
}
