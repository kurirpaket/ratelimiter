package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kurirpaket/throttle"
	"time"
)

type Key string

const StatKey = Key("throttle_stat")
const StatExpiration = time.Minute * 5

const RateKey = Key("throttle_rate")
const RateExpiration = time.Minute * 5

type redisStorage struct {
	cl *redis.Client
}

func NewRedisStorage(addr, username, password string, db int) (throttle.Store, error) {
	cl := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	err := cl.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &redisStorage{cl: cl}, nil
}

func (r *redisStorage) SetStat(ctx context.Context, key string, stat throttle.Stat) error {
	bt, err := json.Marshal(stat)
	if err != nil {
		return err
	}

	return r.cl.Set(ctx, fmt.Sprintf("%s_%s", StatKey, key), string(bt), StatExpiration).Err()
}

func (r *redisStorage) SetRate(ctx context.Context, key string, rate throttle.Rate) error {
	bt, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	return r.cl.Set(ctx, fmt.Sprintf("%s_%s", RateKey, key), string(bt), RateExpiration).Err()
}

func (r *redisStorage) GetRate(ctx context.Context, key string) (*throttle.Rate, error) {
	var rate throttle.Rate

	encRate, err := r.cl.Get(ctx, fmt.Sprintf("%s_%s", RateKey, key)).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}

	if encRate == "" {
		return nil, errors.New("not found")
	}

	err = json.Unmarshal([]byte(encRate), &rate)
	if err != nil {
		return nil, err
	}

	return &rate, nil
}

func (r *redisStorage) GetStat(ctx context.Context, key string) (*throttle.Stat, error) {
	var stat throttle.Stat

	encStat, err := r.cl.Get(ctx, fmt.Sprintf("%s_%s", StatKey, key)).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}

	if encStat == "" {
		return nil, errors.New("not found")
	}

	err = json.Unmarshal([]byte(encStat), &stat)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

func (r *redisStorage) Reset(ctx context.Context, key string) error {
	var currentStat throttle.Stat
	currentStat.ResetAt = time.Now()
	currentStat.TotalHit = 1

	err := r.SetStat(ctx, key, currentStat)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStorage) Increment(ctx context.Context, key string) error {
	currentStat, err := r.GetStat(ctx, key)
	if err != nil {
		return err
	}

	currentStat.TotalHit++

	err = r.SetStat(ctx, key, *currentStat)
	if err != nil {
		return err
	}

	return nil
}
