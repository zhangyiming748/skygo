package redis

import (
	"skygo_detection/service"

	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

const redisPrefix = "_PHCR"

var redisClient redis.Cmdable

var redisLock sync.Mutex

func NewRedis(db int) redis.Cmdable {
	if redisClient == nil {
		redisConfig := service.LoadConfig().Redis
		if len(redisConfig.Addr) > 1 {
			if redisConfig.Auth == "" {
				redisClient = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:    redisConfig.Addr,
					PoolSize: redisConfig.PoolSize,
				})
			} else {
				redisClient = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:    redisConfig.Addr,
					Password: redisConfig.Auth,
					PoolSize: redisConfig.PoolSize,
				})
			}
		} else {
			redisClient = redis.NewClient(&redis.Options{
				Addr:        redisConfig.Addr[0],
				Password:    redisConfig.Auth,
				DB:          db,
				ReadTimeout: redisConfig.Timeout * time.Millisecond,
			})
		}
	}

	return redisClient
}

type Redis_service struct{}

func (s *Redis_service) Get(key string) *redis.StringCmd {
	return NewRedis(0).Get(redisPrefix + key)
}

func (s *Redis_service) Set(key string, value interface{}, t time.Duration) *redis.StatusCmd {
	return NewRedis(0).Set(redisPrefix+key, value, t)
}

func (s *Redis_service) Expire(key string, t time.Duration) *redis.BoolCmd {
	return NewRedis(0).Expire(redisPrefix+key, t)
}

func (s *Redis_service) Exist(key string) *redis.IntCmd {
	return NewRedis(0).Exists(redisPrefix + key)
}

func (s *Redis_service) ZAdd(key string, score float64, member interface{}) *redis.IntCmd {
	return NewRedis(0).ZAdd(redisPrefix+key, &redis.Z{Score: score, Member: member})
}

func (s *Redis_service) ZCount(key string, min, max string) *redis.IntCmd {
	return NewRedis(0).ZCount(redisPrefix+key, min, max)
}

func (s *Redis_service) ZScore(key, member string) *redis.FloatCmd {
	return NewRedis(0).ZScore(redisPrefix+key, member)
}

// 从有序集合中，获取分数范围min max的所有元素
// 不用提供LIMIT相关参数，默认全部查询
func (s *Redis_service) ZRangeScoreAll(key, min string, max string) *redis.StringSliceCmd {
	rangeBy := redis.ZRangeBy{
		Min: min,
		Max: max,
	}
	return NewRedis(0).ZRangeByScore(redisPrefix+key, &rangeBy)
}

func (s *Redis_service) HGetAll(key string) *redis.StringStringMapCmd {
	return NewRedis(0).HGetAll(redisPrefix + key)
}

func (s *Redis_service) HGet(key, field string) *redis.StringCmd {
	return NewRedis(0).HGet(redisPrefix+key, field)
}

func (s *Redis_service) HSet(key, field string, value interface{}) *redis.IntCmd {
	return NewRedis(0).HSet(redisPrefix+key, field, value)
}

func (s *Redis_service) HMSet(key string, fields map[string]interface{}) *redis.BoolCmd {
	return NewRedis(0).HMSet(redisPrefix+key, fields)
}

func (s *Redis_service) HDel(key, field string) *redis.IntCmd {
	return NewRedis(0).HDel(redisPrefix+key, field)
}

// 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (s *Redis_service) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return NewRedis(0).BLPop(timeout, keys...)
}

// 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (s *Redis_service) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return NewRedis(0).BRPop(timeout, keys...)
}

// 从列表中弹出一个值，将弹出的元素插入到另外一个列表中并返回它； 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (s *Redis_service) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return NewRedis(0).BRPopLPush(source, destination, timeout)
}

// 通过索引获取列表中的元素
func (s *Redis_service) LIndex(key string, index int64) *redis.StringCmd {
	return NewRedis(0).LIndex(key, index)
}

// 在列表的元素前或者后插入元素
func (s *Redis_service) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return NewRedis(0).LInsert(key, op, pivot, value)
}

// 在列表的元素前插入元素
func (s *Redis_service) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	return NewRedis(0).LInsertBefore(key, pivot, value)
}

// 在列表的元素后插入元素
func (s *Redis_service) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	return NewRedis(0).LInsertAfter(key, pivot, value)
}

// 获取列表长度
func (s *Redis_service) LLen(key string) *redis.IntCmd {
	return NewRedis(0).LLen(key)
}

// 移出并获取列表的第一个元素
func (s *Redis_service) LPop(key string) *redis.StringCmd {
	return NewRedis(0).LPop(key)
}

// 将一个或多个值插入到列表头部
func (s *Redis_service) LPush(key string, values ...interface{}) *redis.IntCmd {
	return NewRedis(0).LPush(key, values...)
}

// 将一个值插入到已存在的列表头部
func (s *Redis_service) LPushX(key string, value interface{}) *redis.IntCmd {
	return NewRedis(0).LPushX(key, value)
}

// 获取列表指定范围内的元素
func (s *Redis_service) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return NewRedis(0).LRange(key, start, stop)
}

// 移除列表元素
func (s *Redis_service) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return NewRedis(0).LRem(key, count, value)
}

// 移除列表元素
func (s *Redis_service) LRemoveAll(key string) {
	for i := int64(0); i < s.LLen(key).Val(); i++ {
		s.LPop(key)
	}
}

// 通过索引设置列表元素的值
func (s *Redis_service) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return NewRedis(0).LSet(key, index, value)
}

// 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
func (s *Redis_service) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return NewRedis(0).LTrim(key, start, stop)
}

// 移除列表的最后一个元素，返回值为移除的元素。
func (s *Redis_service) RPop(key string) *redis.StringCmd {
	return NewRedis(0).RPop(key)
}

// 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
func (s *Redis_service) RPopLPush(source, destination string) *redis.StringCmd {
	return NewRedis(0).RPopLPush(source, destination)
}

// 在列表中添加一个或多个值
func (s *Redis_service) RPush(key string, values ...interface{}) *redis.IntCmd {
	return NewRedis(0).RPush(key, values...)
}

// 为已存在的列表添加值
func (s *Redis_service) RPushX(key string, value interface{}) *redis.IntCmd {
	return NewRedis(0).RPushX(key, value)
}
