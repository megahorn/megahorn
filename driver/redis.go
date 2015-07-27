package driver

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type RedisDriver struct {
	redis     redis.Conn
	key       string
	command   string
	wg        sync.WaitGroup
	sendQueue chan string
	quit      chan bool
}

func (r *RedisDriver) Configure(config map[string]string) (err error) {
	network, ok := config["network"]
	if !ok {
		network = "tcp"
	}

	address, ok := config["address"]
	if !ok {
		address = "localhost:6379"
	}

	key, ok := config["key"]
	if !ok {
		return errors.New("key is required")
	}
	r.key = key

	command, ok := config["command"]
	if !ok {
		command = "APPEND"
	}
	r.command = command

	r.redis, err = redis.Dial(network, address)
	if err != nil {
		return
	}

	password, ok := config["password"]
	if ok {
		if _, err = r.redis.Do("AUTH", password); err != nil {
			r.redis.Close()
			return err
		}
	}

	r.sendQueue = make(chan string, 32)
	r.quit = make(chan bool, 1)

	go func() {
		for {
			select {
			case data := <-r.sendQueue:
				r.redis.Do(r.command, r.key, data)
				r.wg.Done()
			case <-r.quit:
				return
			}
		}
	}()

	return
}

func (r *RedisDriver) Write(data []byte) (int, error) {
	r.wg.Add(1)
	r.sendQueue <- string(data)
	return len(data), nil
}

func (r *RedisDriver) Close() error {
	r.wg.Wait()
	r.redis.Flush()
	r.quit <- true
	return r.redis.Close()
}
