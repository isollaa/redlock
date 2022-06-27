package redlock_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "redlock"

	"github.com/stretchr/testify/assert"
)

func setConfig() Config {
	var conf Config
	conf.Redis.Host = "localhost"
	conf.Redis.Port = 6379
	conf.Redis.Password = ""
	conf.Redis.DB = 0

	return conf
}

func TestWithMutexes(t *testing.T) {
	redSync := New(setConfig())
	tests := map[string]struct {
		totalTest    int
		execDuration int
		expirySecond int
		wantErr      error
	}{
		"SUCCESS": { // success scenario represent there is no fail simulation on its scenario
			totalTest:    3,
			execDuration: 1,
			expirySecond: 1,
			wantErr:      nil,
		},
		"FAIL": { // fail scenario represent there is fail simulation on its scenario
			totalTest:    5,
			execDuration: 1,
			expirySecond: 1,
			wantErr:      fmt.Errorf("redsync: failed to acquire lock"),
		},
	}

	for k, v := range tests {
		t.Run(fmt.Sprintf("%s SCENARIO", k), func(t *testing.T) {
			redSync.ExpirySeconds = v.expirySecond
			var err error
			var wg sync.WaitGroup

			for i := 0; i < v.totalTest; i++ {
				wg.Add(1)
				go func(t *testing.T, i int, wg *sync.WaitGroup) {
					defer wg.Done()
					t.Run(fmt.Sprintf("SIMULATION %d", i), func(*testing.T) {
						keys := []string{}
						for i := 0; i < v.totalTest; i++ {
							keys = append(keys, fmt.Sprintf("lock-key-test-%d", i))
						}
						if e := redSync.WithMutexes(keys, func() error {
							time.Sleep(time.Duration(v.execDuration) * time.Millisecond)
							return nil
						}); e != nil {
							err = e
						}
					})
				}(t, i, &wg)
			}

			wg.Wait()
			assert.Equal(t, v.wantErr, err)
		})
	}
}
