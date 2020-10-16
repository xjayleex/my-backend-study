package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"sync"
)

const (
	EnrollmentDB = 2
)

type User struct {
	Mail				string	`json:"mail"`
	Username			string	`json:"name"`
}

func NewUser(username string, mail string) *User {
	return &User {
		Username: username,
		Mail: mail,
	}
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

type UserStore interface {
	Save(e interface{}) error
	Find(key string) (interface{},error)
}

type RedisStore struct {
	mtx		sync.RWMutex
	cli		*redis.Client
}

type RedisClientOpts struct {
	Address	string
	Port	string
	DB 		int
}

func NewRedisUserStore (opts *RedisClientOpts) (*RedisStore, error) {
	if opts.Address == "" {
		return nil, errors.New("Address Required.")
	}

	if opts.Port == "" {
		return nil, errors.New("Port Required")
	}

	rs := &RedisStore{
		mtx: sync.RWMutex{},
		cli: redis.NewClient(&redis.Options{
			Addr: opts.Address + ":" + opts.Port,
			DB: opts.DB,
		}),
	}
	return rs, nil
}

func (rs *RedisStore) Save (e interface{}) error {
	user, ok := e.(*User)
	if !ok {
		return errors.New("Type mismatch error")
	}
	ctx := context.Background()
	if err := rs.Ping(ctx); err != nil {
		fmt.Println(err)
		return errors.New("Not connected to redis")
	}

	rs.mtx.Lock()
	defer rs.mtx.Unlock()

	boolCmd := rs.cli.SetNX(ctx,user.Mail,user,0)

	flag, err := boolCmd.Result()

	if err != nil {
		return err
	}

	if !flag {
		return errors.New("Mail addr is exists already")
	} else {
		return nil
	}
}

func (rs *RedisStore) Find (key string) (interface{} ,error) {
	ctx := context.Background()
	if err := rs.Ping(ctx); err != nil {
		return nil, errors.New("Not connected to redis")
	}

	rs.mtx.RLock()
	defer rs.mtx.RUnlock()
	cmd := rs.cli.Get(ctx, key)

	if _, err := cmd.Result(); err != nil {
		return nil, errors.New("Key not exists")
	}
	return cmd, nil
}


func (rs *RedisStore) Ping(ctx context.Context) (err error) {
	pong, err := rs.cli.Ping(ctx).Result()
	if err != nil {
		return errors.New("Not connected to Redis")
	}
	if pong != "PONG" {
		return errors.New("Internal Error")
	}
	return nil
}