package store

import (
	"context"

	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

type Code uint32

const (
	EnrollmentDB = 2
)

const (
	ErrNoRedisServerAddress Code = 0 + iota
	ErrNoRedisPort
	ErrTypeAssertion
	ErrNoConnWithRedis
	ErrKeyNotExists
	ErrKeyExistsAlready
	ErrInternal
)

type StoreError struct {
	Code	Code
	Msg		string
}

func (se *StoreError) Error() string {
	return fmt.Sprintf("Code : %d Message : %s",se.Code,se.Msg)
}

func Error (c Code, msg string) *StoreError {
	return &StoreError{
		Code: c,
		Msg:  msg,
	}
}

type RedisValue interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data[]byte) error
}

type UserStore interface {
	Save(key string, value RedisValue) error
	Find(key string) (interface{},error)
}

type RedisStore struct {
	mtx		*sync.RWMutex
	cli		*redis.Client
}

type RedisClientOpts struct {
	Address	string
	Port	string
	DB 		int
}

func NewRedisUserStore (opts *RedisClientOpts) (*RedisStore, error) {
	if opts.Address == "" {
		return nil, Error(ErrNoRedisServerAddress,"Address Required")
	}

	if opts.Port == "" {
		return nil, Error(ErrNoRedisPort, "Port Required")
	}

	rs := &RedisStore{
		mtx: &sync.RWMutex{},
		cli: redis.NewClient(&redis.Options{
			Addr: opts.Address + ":" + opts.Port,
			DB: opts.DB,
		}),
	}
	return rs, nil
}

func (rs *RedisStore) Save (key string, value RedisValue) error {
	ctx := context.Background()
	if err := rs.Ping(ctx); err != nil {
		fmt.Println(err)
		return Error(ErrNoConnWithRedis, "Not Connected to redis")
	}

	rs.mtx.Lock()
	defer rs.mtx.Unlock()

	boolCmd := rs.cli.SetNX(ctx, key, value,0)

	flag, err := boolCmd.Result()

	if err != nil {
		return err
	}

	if !flag {
		return Error(ErrKeyExistsAlready, "Key exists already")
	} else {
		return nil
	}
}

func (rs *RedisStore) Find (key string) (interface{} ,error) {
	ctx := context.Background()
	if err := rs.Ping(ctx); err != nil {
		return nil, Error(ErrNoConnWithRedis,"Not connected to redis")
	}

	rs.mtx.RLock()
	defer rs.mtx.RUnlock()
	cmd := rs.cli.Get(ctx, key)

	if _, err := cmd.Result(); err != nil {
		return nil, Error(ErrKeyNotExists, "Key not exists")
	}
	return cmd, nil
}


func (rs *RedisStore) Ping(ctx context.Context) (err error) {
	pong, err := rs.cli.Ping(ctx).Result()
	if err != nil {
		return Error(ErrNoConnWithRedis,"Not connected to redis")
	}
	if pong != "PONG" {
		return Error(ErrInternal, "Internal Error")
	}
	return nil
}