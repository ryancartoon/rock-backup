package mq

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	Client *redis.Client
	Queue  string
	Ctx    context.Context
}

func NewRedisQueue(rdb *redis.Client, queue string) *RedisQueue {
	return &RedisQueue{
		Client: rdb,
		Queue:  queue,
		Ctx:    context.Background(),
	}
}

func (rq *RedisQueue) Enqueue(message string) error {
	err := rq.Client.RPush(rq.Ctx, rq.Queue, message).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %v", err)
	}
	fmt.Printf("Enqueued message: %s\n", message)
	return nil
}

func (rq *RedisQueue) Dequeue() (string, error) {
	message, err := rq.Client.LPop(rq.Ctx, rq.Queue).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no messages in queue")
	} else if err != nil {
		return "", fmt.Errorf("failed to dequeue message: %v", err)
	}
	fmt.Printf("Dequeued message: %s\n", message)
	return message, nil
}
