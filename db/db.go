package db

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

var conn redis.Conn
var err error

func init() {
	conn, err = redis.Dial("tcp", "redis-19300.c1.us-east1-2.gce.cloud.redislabs.com:19300", redis.DialPassword("8Nl34nx5IKKMTGWYJj063HFtkRgtn55J"))
	if err != nil {
		log.Fatal(err)
	}
}

// SetSomething test write to redis
func SetSomething() {
	_, err = conn.Do("HMSET", "album:2", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Electric Ladyland added!")
	return
}
