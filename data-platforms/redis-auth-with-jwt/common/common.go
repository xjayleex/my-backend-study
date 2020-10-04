package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
)
var (
	ErrKeyNotExists = errors.New("Key Not exist.")
	ErrKeyExistsAlready = errors.New("Key exists already.")
	ErrNotConnected = errors.New("Not Connected.")
	ErrNilContext = errors.New("Context is nil value.")
	ErrNoAddress = errors.New("Address Required.")
	ErrNoPort = errors.New("Port Required.")
	ErrUnknown = errors.New("Unknown error occurs")
)


type User struct {
	Username			string
	HashedPassword		string
	Mail				string
}


func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func NewUser(username string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash pw: %w", err)
	}

	user := &User {
		Username: username,
		HashedPassword: string(hashedPassword),
	}
	return user, nil
}

func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

type UserStore interface {
	SignUp(user *User) error
	Find(username string) (*User, error)
}

type RedisUserStore	struct {
	mutex       sync.RWMutex
	userInRedis *RedisClient
}

func NewRedisUserStore(opts *RedisClientOpts) *RedisUserStore {
	if opts == nil {
		fmt.Println("Redis Options Needed.")
		os.Exit(1)
	}
	return &RedisUserStore{
		userInRedis: NewRedisClient(opts),
	}
}

func (store *RedisUserStore) Find (username string) (*User, error ) {
	return nil, nil
}


func (store *RedisUserStore) SignUp (user *User ) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	err := store.userInRedis.SetNX(&RData{
		Key:   user.Username,
		Value: user,
	})
	return err
}

func main() {
	rs := NewRedisUserStore(&RedisClientOpts{
		Address: "localhost",
		Port:    "6379",
		DB:      1,
	})

	err := rs.SignUp(&User{
		Username:       "jay1",
		HashedPassword: "1234",
		Mail:           "bigdata304@gmail.com",
	})
	fmt.Println(err)
}