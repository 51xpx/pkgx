package locker

import (
	"testing"
	
	"github.com/go-redis/redis"
)

func TestRedisLocker_Lock(t *testing.T) {
	opt := &redis.Options{
		Addr: "81.69.233.94:30017",
		Password: "",
		DB: 0,
	}
	
	client := redis.NewClient(opt)
	
	locker := NewRedisLocker(client, "gcore:test:%d", 1)
	flag := locker.Lock()
	if flag{
		t.Log("lock success", locker.Key())
	}else{
		t.Log("lock fai", locker.Key())
	}
	uflag := locker.Unlock()
	t.Log("unlock", uflag)
}
