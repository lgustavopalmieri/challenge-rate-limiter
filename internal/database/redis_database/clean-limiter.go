package redis_database

import "context"

func (rd *RedisRepositoryDb) CleanLimiter(ctx context.Context, ip, token string) error {
	key := rd.MakeKey(ip, token)
	_, err := rd.redisClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
