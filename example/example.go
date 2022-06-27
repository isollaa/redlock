package main

import "redlock"

func main() {
	// setup config
	var conf redlock.Config

	// setup redis connection
	conf.Redis.Host = "localhost"
	conf.Redis.Port = 6379
	conf.Redis.Password = ""
	conf.Redis.DB = 0

	// setup lock expiration in seconds [optional]
	conf.ExpirySeconds = 30

	// create an instance of Redlock to be used to obtain a mutual exclusion lock.
	rl := redlock.New(conf)

	// declare list of key that requires the lock
	keys := []string{"key1", "key2", "key3"}

	err := rl.WithMutexes(keys, func() error {
		// do your work that requires the lock here . . .

		return nil
	})

	// err contains error caused by lock / unlock process or your work
	if err != nil {
		panic(err)
	}
}
