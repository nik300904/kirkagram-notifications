package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	internalConfig "kirkagram-notification/internal/config"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

var (
	ErrNoFollowers      = errors.New("no followers found")
	ErrLikeNotFound     = errors.New("Like not found")
	ErrFollowerNotFound = errors.New("Follower not found")
)

func New(cfg *internalConfig.Config) *sql.DB {
	// info := ConnectionInfo{
	// 	Host:     "localhost",
	// 	Port:     5432,
	// 	Username: "myuser",
	// 	DBName:   "mydatabase",
	// 	SSLMode:  "disable",
	// 	Password: "mypassword",
	// }

	connStr := "host=localhost port=5467 user=user password=12345 dbname=kirkagram sslmode=disable"

	// db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
	// 	info.Host, info.Port, info.Username, info.Password, info.DBName, info.SSLMode))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("failed to connect to database: %v\n", err)
		panic(err)
	}

	return db
}
