package session

import (
	"encoding/gob"
	"goroutines/subscription_project/data"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func Initialize() *scs.SessionManager {

	gob.Register(data.User{})

	session := scs.New()

	session.Store = redisstore.New(getRegisPool())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session

}

func getRegisPool() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", strings.TrimSpace(os.Getenv("REDIS")))
		},
	}

	return redisPool

}
