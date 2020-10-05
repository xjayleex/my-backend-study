package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"time"
)
var (
	ErrKeyNotExists = errors.New("Key Not exist.")
	ErrKeyExistsAlready = errors.New("Key exists already.")
	ErrNotConnected = errors.New("Not Connected.")
	ErrNilContext = errors.New("Context is nil value.")
	ErrNoAddress = errors.New("Address Required.")
	ErrNoPort = errors.New("Port Required.")
	ErrUnknown = errors.New("Unknown error occurs")
	ErrUnexpectedToken = errors.New("Unexpected token signing method.")
	ErrInvalidClaims = errors.New("Invalid token claims")
	ErrPasswordHash = errors.New("cannot hash pw")
	ErrNilUserObject = errors.New("Nil user error")
	ErrNoAuthServer = errors.New("No Auth server")
	ErrIncorrectInfo = errors.New("Incorrect mail or password")
)

const (
	UserStoreDB 	= 1
	TokenStoreDB	= 2
)

type User struct {
	Mail				string	`json:"mail"`
	Username			string	`json:"username"`
	HashedPassword		string	`json:"password"`
}


func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func NewUser(mail string, username string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHash
	}

	user := &User {
		Username: username,
		Mail: mail,
		HashedPassword: string(hashedPassword),
	}
	return user, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))

	return err == nil
}

type UserStore interface {
	SignUp(user *User) error
	Find(mail string) (*User, error)
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

func (store *RedisUserStore) Find (mail string) (*User, error ) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	cmd, err := store.userInRedis.Get(mail)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, ErrNilUserObject
	}
	val, err := cmd.Bytes()
	if err != nil {
		return nil, errors.New("Marshal/Unmarshal Error")
	}
	user := &User{}
	err = json.Unmarshal(val,user)
	return user, err
}


func (store *RedisUserStore) SignUp (user *User ) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	err := store.userInRedis.SetNX(&RData{
		Key:   user.Mail,
		Value: user,
	})
	return err
}

type JWTManager struct {
	secretKey		string
	tokenDuration	time.Duration
}

type UserClaims struct { // Sensitive한 정보는 토큰에 담지 말 것...
	jwt.StandardClaims
	Mail string `json:"mail"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager{
	return &JWTManager{secretKey: secretKey, tokenDuration: tokenDuration}
}

func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Mail: user.Mail,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error){
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedToken
			}
			return []byte(manager.secretKey), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	return claims, nil
}