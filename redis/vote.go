package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数值
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

//directive的不同情况 1/0/-1
//

func VoteForPost(userID, postID string, value float64) error {
	//投票功能
	// 1.判断投票限制
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2.更新帖子分数
	// 查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()

	//如果此次投票的值和之前保存的值一致，就提醒不允许重复投票
	if value == ov {
		return ErrVoteRepested
	}

	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value)
	_, err := client.ZIncrBy(getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID).Result()
	if err != nil {
		return err
	}

	if value == 0 {
		_, err = client.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result()
	} else {
		// 3.记录用户为该帖子投票的数据
		_, err = client.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value,  //当前用户投的赞同/反对
			Member: userID, //
		}).Result()
	}
	return err
}
