//go:build no_wire

package main

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Server struct {
	redisCmd redis.Cmdable
	db       *gorm.DB
}

func NewServer(redisCmd redis.Cmdable, db *gorm.DB) *Server {
	return &Server{redisCmd: redisCmd, db: db}
}

func (s *Server) Run() error {
	return nil
}

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{Addr: "localhost:6379"})
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {
	redisCmd := InitRedis()
	db := InitDB()
	server := NewServer(redisCmd, db)
	if err := server.Run(); err != nil {
		log.Println(err)
	}
}
