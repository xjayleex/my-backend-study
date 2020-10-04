package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"os"
)
type RedisClient struct {
	client		*redis.Client
	address 	string
	port		string
	db			int
	ctx 		context.Context
}

type RData struct {
	Key string
	Value interface{}
}

type RValue struct {
	Password string	`json:"password"`
	Mail	string `json:"mail"`
}

func (rd *RData) MarshalBinary() ([]byte, error) {
	return json.Marshal(rd)
}

func (rd *RData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data,rd)
}


func (rv *RValue) String() string {
	return fmt.Sprintf("%s %s",rv.Password,rv.Mail)
}

type RedisClientOpts struct {
	Address	string
	Port	string
	DB 		int
}
func NewRedisClient(opts *RedisClientOpts) (c *RedisClient) {
	if c == nil {
		c = &RedisClient{}
	}
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if opts.Address == ""  {
		err := errors.Errorf("Address Required")
		trap(err)
		return
	}
	if opts.Port == "" {
		err := errors.Errorf("Port Required")
		trap(err)
		return
	}
	c.address, c.port, c.db = opts.Address, opts.Port, opts.DB
	c.client = redis.NewClient(&redis.Options{
		Addr: c.address + ":" + c.port,
		DB: c.db,
	})
	return
}


func (rc *RedisClient) SetNX(data *RData) error {
	if err := rc.Ping() ; err != nil {
		return ErrNotConnected
	}
	boolCmd := rc.client.SetNX(rc.ctx, data.Key, data.Value,0)
	flag , err := boolCmd.Result()
	if err != nil {
		return err
	}
	if !flag {
		return ErrKeyExistsAlready
	} else {
		return nil
	}
}

func (rc *RedisClient) Get(key string) (*redis.StringCmd, error) {
	if err := rc.Ping() ; err != nil {
		return nil, ErrNotConnected
	}
	cmd := rc.client.Get(rc.ctx, key)

	if _, err := cmd.Result(); err != nil {
		return nil, ErrKeyNotExists
	}
	return cmd, nil
}


func (rc *RedisClient) Ping() (err error) {
	if rc.ctx == nil {
		return ErrNilContext
	}
	pong, err := rc.client.Ping(rc.ctx).Result()
	if err != nil {
		return ErrNotConnected
	}
	if pong != "PONG" {
		return ErrNotConnected
	}
	return nil
}

func trap (err error) {
	if err == nil {
		return
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
func main() {
	rc, err := NewRedisClient(&RedisClientOpts{
		Address: "localhost",
		Port: "6379",
		DB: 1,
	})
	trap(err)
	pong, err := rc.client.Ping(rc.ctx).Result()
	trap(err)
	if pong != "PONG" {
		fmt.Println("Not Connected.")
		os.Exit(1)
	}
	rc.SetNX(&RData{"users:jaehyunlee",
		&RValue{"password", "bigdata304@gmail.com" }})
	rc.Get("users:jaehyunlee")
	rv, err := rc.GetBytes("users:jaehyunlee")
	trap(err)
	fmt.Println(rv.String())

}*/

func (rc *RedisClient) GetBytes(key string) (*RValue,error) {
	val, err := rc.client.Get(rc.ctx, key).Bytes()
	if err != nil {
		fmt.Println("Key doesn't exist.")
		return nil ,err
	}
	rv := &RValue{}
	json.Unmarshal(val,rv)
	return rv, nil
}
