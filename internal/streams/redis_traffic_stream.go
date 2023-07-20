package streams

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTrafficStream struct {
	client *redis.Client
	stream string
}

func NewRedisTrafficStream(stream string, client *redis.Client) *RedisTrafficStream {
	return &RedisTrafficStream{
		client: client,
		stream: stream,
	}
}

func (r *RedisTrafficStream) Add(traffic int64) error {
	_, err := r.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: r.stream,
		MaxLen: 3600 * 8,
		Values: map[string]string{
			"traffic": strconv.Itoa(int(traffic)),
		},
	}).Result()
	return err
}
func (r *RedisTrafficStream) Range(from time.Time, to time.Time) ([]int64, error) {
	messages, err := r.client.XRange(context.Background(), r.stream, strconv.Itoa(int(from.UnixMilli())), strconv.Itoa(int(to.UnixMilli()))).Result()
	if err != nil {
		return []int64{}, err
	}
	traffics := []int64{}
	for _, message := range messages {
		traffics = append(traffics, message.Values["traffic"].(int64))
	}
	return traffics, nil
}
