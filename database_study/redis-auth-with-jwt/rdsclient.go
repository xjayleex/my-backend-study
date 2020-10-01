package main
import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"os"
)
type RedisClient struct {
	client		*redis.Client
	address 	string
	port		string
}

type RedisClientOpts struct {
	Address	string
	Port	string
}
func NewRedisClient(opts *RedisClientOpts) (c *RedisClient, err error) {
	if opts.Address == "" {
		err = errors.Errorf("Address Required")
		return
	}
	if opts.Port == "" {
		err = errors.Errorf("Port Required")
	}
	c.address = opts.Address
	c.client = redis.NewClient(&redis.Options{
		Addr: c.address + ":" + c.port,
	})
	return
}


func main() {
	c, err := NewRedisClient(&RedisClientOpts{
		Address: "dn2",
		Port: "6379",
	})
	trap(err)

	pong, err := c.client.Ping().Result()
	trap(err)
	fmt.Println(pong)
}

func trap (err error) {
	if err == nil {
		return
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}