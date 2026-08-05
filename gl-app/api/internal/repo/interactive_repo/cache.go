package interactive_repo

import (
	"context"
	"fmt"
	"strconv"

	"gl-app/api/internal/constants"

	redisv9 "github.com/redis/go-redis/v9"
)

type LikeStateCache interface {
	OverlayUserLikedStates(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType, fallback map[int64]bool) (map[int64]bool, error)
	MutateIfPresent(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked bool) (LikeMutation, bool, error)
	Mutate(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked, initialLiked bool, initialCount, initialVersion int64) (LikeMutation, error)
}

type likeStateCacheImpl struct{ redis redisv9.UniversalClient }

func NewLikeStateCache(client redisv9.UniversalClient) LikeStateCache {
	return &likeStateCacheImpl{redis: client}
}

// OverlayUserLikedStates 优先采用 Redis 中尚未异步落库的最新状态，缓存缺失项沿用数据库结果。
func (c *likeStateCacheImpl) OverlayUserLikedStates(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType, fallback map[int64]bool) (map[int64]bool, error) {
	if len(resourceIDs) == 0 { return fallback, nil }
	keys := make([]string, 0, len(resourceIDs))
	for _, resourceID := range resourceIDs {
		keys = append(keys, fmt.Sprintf("interaction:like:%s:%d:user:%d", bizType, resourceID, userID))
	}
	values, err := c.redis.MGet(ctx, keys...).Result()
	if err != nil { return fallback, err }
	for index, value := range values {
		if value == nil { continue }
		fallback[resourceIDs[index]] = fmt.Sprint(value) == "1"
	}
	return fallback, nil
}

var mutateLikeScript = redisv9.NewScript(`
if ARGV[5] == '0' and (redis.call('EXISTS', KEYS[1]) == 0 or redis.call('EXISTS', KEYS[2]) == 0 or redis.call('EXISTS', KEYS[3]) == 0) then
  return {-1, 0, 0}
end
if redis.call('EXISTS', KEYS[1]) == 0 then redis.call('SET', KEYS[1], ARGV[1]) end
if redis.call('EXISTS', KEYS[2]) == 0 then redis.call('SET', KEYS[2], ARGV[2]) end
if redis.call('EXISTS', KEYS[3]) == 0 then redis.call('SET', KEYS[3], ARGV[4]) end
local current = redis.call('GET', KEYS[1])
if current == ARGV[3] then
  local version = redis.call('GET', KEYS[3]) or '0'
  return {0, tonumber(redis.call('GET', KEYS[2]) or '0'), tonumber(version)}
end
redis.call('SET', KEYS[1], ARGV[3])
local count
if ARGV[3] == '1' then count = redis.call('INCR', KEYS[2])
else
  count = tonumber(redis.call('GET', KEYS[2]) or '0')
  if count > 0 then count = redis.call('DECR', KEYS[2]) end
end
local version = redis.call('INCR', KEYS[3])
redis.call('EXPIRE', KEYS[1], 2592000)
redis.call('EXPIRE', KEYS[2], 2592000)
redis.call('EXPIRE', KEYS[3], 2592000)
return {1, count, version}
`)

func (c *likeStateCacheImpl) MutateIfPresent(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked bool) (LikeMutation, bool, error) {
	values, err := c.runMutation(ctx, userID, resourceID, bizType, liked, false, 0, 0, false)
	if err != nil {
		return LikeMutation{}, false, err
	}
	if values[0] == -1 {
		return LikeMutation{}, false, nil
	}
	return valuesToMutation(values, liked)
}

func (c *likeStateCacheImpl) Mutate(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked, initialLiked bool, initialCount, initialVersion int64) (LikeMutation, error) {
	values, err := c.runMutation(ctx, userID, resourceID, bizType, liked, initialLiked, initialCount, initialVersion, true)
	if err != nil {
		return LikeMutation{}, err
	}
	mutation, _, err := valuesToMutation(values, liked)
	return mutation, err
}

func (c *likeStateCacheImpl) runMutation(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked, initialLiked bool, initialCount, initialVersion int64, initialize bool) ([]int64, error) {
	prefix := fmt.Sprintf("interaction:like:%s:%d", bizType, resourceID)
	stateKey := fmt.Sprintf("%s:user:%d", prefix, userID)
	values, err := mutateLikeScript.Run(ctx, c.redis, []string{stateKey, prefix + ":count", stateKey + ":version"},
		boolString(initialLiked), strconv.FormatInt(initialCount, 10), boolString(liked), strconv.FormatInt(initialVersion, 10), boolString(initialize)).Int64Slice()
	if err != nil || len(values) != 3 {
		if err == nil {
			err = fmt.Errorf("Redis点赞脚本返回值不完整")
		}
		return nil, err
	}
	return values, nil
}

func valuesToMutation(values []int64, liked bool) (LikeMutation, bool, error) {
	if len(values) != 3 {
		return LikeMutation{}, false, fmt.Errorf("Redis点赞脚本返回值不完整")
	}
	return LikeMutation{Changed: values[0] == 1, Liked: liked, Count: values[1], Version: values[2]}, true, nil
}

func boolString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
