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
	Value *RValue
}

type RValue struct {
	Password string
	Mail	string
}

func (rd *RData) MarshalBinary() ([]byte, error) {
	return json.Marshal(rd)
}
func (rv *RValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(rv)
}

type RedisClientOpts struct {
	Address	string
	Port	string
	DB 		int
}
func NewRedisClient(opts *RedisClientOpts) (c *RedisClient, err error) {
	if c == nil {
		c = &RedisClient{}
	}
	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if opts.Address == ""  {
		err = errors.Errorf("Address Required")
		return
	}
	if opts.Port == "" {
		err = errors.Errorf("Port Required")
	}
	c.address, c.port, c.db = opts.Address, opts.Port, opts.DB
	c.client = redis.NewClient(&redis.Options{
		Addr: c.address + ":" + c.port,
		DB: c.db,
	})
	return
}


func (rc *RedisClient) Setnx(data *RData) {
	boolCmd := rc.client.SetNX(rc.ctx, data.Key, data.Value,0)
	flag , err := boolCmd.Result()
	trap(err)
	if !flag {
		fmt.Println("Key Exists Already.")
	} else {
		fmt.Println("Success on creating key.")
	}

}


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
	rc.Setnx(&RData{"users:jaehyunlee",
		&RValue{"password", "bigdata304@gmail.com" }})
}

func trap (err error) {
	if err == nil {
		return
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}