package auth

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey = []byte("Some Secrete value")
	ctx    = context.Background()
	rdb    *redis.Client
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func InitializeRedis(redisClient *redis.Client) {
	rdb = redisClient
}

func Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, username, hashedPassword, 0).Err()
	if err != nil {
		log.Println(err)
	}
	return nil
}

func Authenticate(username, password string) (string, error) {
	storedPassword, err := rdb.Get(ctx, username).Result()
	if err != nil {
		return "There is some error ", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return "", err
	}
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
