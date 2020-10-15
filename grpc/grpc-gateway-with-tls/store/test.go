package store

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	rs, err := NewRedisUserStore(&RedisClientOpts{
		Address: "localhost",
		Port:    "6379",
		DB:      2,
	})

	if err != nil {
		fmt.Println(err)
	}

	//us := UserStore(rs)

	us := UserStore(rs)

	err = us.Save(NewUser("jaehyun","bigdata304@gmail.com"))
	if err != nil {
		fmt.Println(err)
	}
	key := "bigdata304@gmail.com"
	value, err := us.Find(key)
	asserted, ok := value.(*redis.StringCmd)
	if !ok {
		fmt.Println("Error")
	}
	unmarshaled, err := asserted.Bytes()
	if err != nil {
		fmt.Println(err)
	}
	user := &User{}
	err = json.Unmarshal(unmarshaled, user)
	if err == nil {
		fmt.Println(user.Mail,user.Username)
	}

}
