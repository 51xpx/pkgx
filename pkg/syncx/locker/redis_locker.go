package locker

import (
	"fmt"
	"time"
	
	"github.com/go-redis/redis"
)

// RedisLocker Redis 锁
type RedisLocker struct {
	client   *redis.Client
	key      string
	retries  int
	interval time.Duration
	timeout  time.Duration
}

func NewRedisLocker(client *redis.Client, key string, keyFmtArgs ...interface{}) *RedisLocker {
	return &RedisLocker{
		client:   client,
		key:      fmt.Sprintf(key, keyFmtArgs...),
		retries:  5,
		interval: 100 * time.Millisecond,
		timeout:  2 * time.Second,
	}
}

func (l *RedisLocker) WithConfig(retries int, interval, timeout time.Duration) *RedisLocker {
	
	if retries <= 1 {
		retries = 1
	}
	
	if min := 5 * time.Millisecond; interval < min {
		interval = min
	}
	
	l.retries = retries
	l.interval = interval
	l.timeout = timeout
	
	return l
}

func (l *RedisLocker) Key() string {
	return l.key
}

func (l *RedisLocker) Touch(duration time.Duration) error {
	return l.client.Expire(l.key, duration).Err()
}

func (l *RedisLocker) TouchAt(t time.Time) error {
	return l.client.ExpireAt(l.key, t).Err()
}

func (l *RedisLocker) Lock() bool {
	for i := 0; i <= l.retries; i++ {
		ok, err := l.client.SetNX(l.key, 1, l.timeout).Result()
		if err == nil && ok {
			return true
		}
		time.Sleep(l.interval)
	}
	return false
}

func (l *RedisLocker) Unlock() bool {
	for i := 0; i <= l.retries; i++ {
		if err := l.client.Del(l.key).Err(); err == nil {
			return true
		}
		time.Sleep(l.interval)
	}
	return false
}
