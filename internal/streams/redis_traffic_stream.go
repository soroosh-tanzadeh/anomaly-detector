package streams

import (
	"context"
	"fmt"
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

func (r *RedisTrafficStream) Add(ctx context.Context, traffic float64) error {
	_, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		MaxLen: 3600 * 8,
		Values: map[string]string{
			"traffic": fmt.Sprintf("%.18f", traffic),
		},
	}).Result()
	return err
}
func (r *RedisTrafficStream) Range(ctx context.Context, from time.Time, to time.Time) ([]float64, error) {
	messages, err := r.client.XRange(context.Background(), r.stream, strconv.Itoa(int(from.UnixMilli())), strconv.Itoa(int(to.UnixMilli()))).Result()
	if err != nil {
		return []float64{}, err
	}
	traffics := []float64{}
	for _, message := range messages {
		traffic, _ := strconv.ParseFloat(message.Values["traffic"].(string), 64)
		traffics = append(traffics, traffic)
	}
	return traffics, nil
}
