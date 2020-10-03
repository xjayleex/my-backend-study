package common

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
)

var (
	ErrKeyNotExists = errors.New("Key Not exist.")
	ErrNotConnected = errors.New("Not Connected.")
	ErrNilContext = errors.New("Context is nil value.")
	ErrNoAddress = errors.New("Address Required.")
	ErrNoPort = errors.New("Port Required.")
)

type User struct {
	Username			string
	HashedPassword		string
	Mail				string
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
	Save(user *User)
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
}


func (store *RedisUserStore) Save (user *User ) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	cmd, err := store.userInRedis.Get(user.Username)
	// err != nil -> key does not exist -> save 가능
	// err == nil -> save 불가능
	if err != nil && errors.Is(err, ErrKeyNotExists){

	} else { // error doesnt exist || error is not keynotexist err

	}

}