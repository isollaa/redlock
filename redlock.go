package redlock

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

// const refer to https://github.com/go-redsync/redsync/blob/52511c81d0a572cb8b849331f9a3ecb69a5651ea/redsync.go#L10
const (
	minRetryDelayMilliSec = 50
	maxRetryDelayMilliSec = 250
)

// Block contains an action block
type Block func() error

// Redlock contains redsync client with its configurable options
type Redlock struct {
	Client        *redsync.Redsync
	ExpirySeconds int
}

// New create an instance of Redlock to be used to obtain a mutual exclusion lock.
func New(conf Config) Redlock {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		DB:       conf.Redis.DB,
		Password: conf.Redis.Password,
	})

	return Redlock{
		Client:        redsync.New(goredis.NewPool(client)),
		ExpirySeconds: conf.ExpirySeconds,
	}
}

// WithMutexes performs action with mutexes
func (r *Redlock) WithMutexes(keys []string, block Block) (err error) {
	mtxs := r.NewMutexes(keys)
	if err = mtxs.Lock(); err != nil {
		return
	}

	defer func() {
		e := mtxs.Unlock()
		if err == nil {
			err = e
		}
	}()

	return block()
}

// NewMutexes returns list of new distributed mutex with given key
func (r *Redlock) NewMutexes(keys []string) (res Mutexes) {
	var opts []redsync.Option
	if r.ExpirySeconds > 0 {
		exprMilliSec := r.ExpirySeconds * 1000
		opts = []redsync.Option{
			redsync.WithExpiry(time.Duration(exprMilliSec) * time.Millisecond),
			redsync.WithTries(int(math.Ceil(float64(exprMilliSec) / float64(maxRetryDelayMilliSec)))),
			redsync.WithRetryDelay(time.Duration(rand.Intn(maxRetryDelayMilliSec-minRetryDelayMilliSec)+minRetryDelayMilliSec) * time.Millisecond),
		}
	}

	unqKeys := make(map[string]struct{})
	for _, v := range keys {
		if _, ok := unqKeys[v]; !ok {
			res = append(res, r.Client.NewMutex(v, opts...))
		}
		unqKeys[v] = struct{}{}
	}

	return
}
