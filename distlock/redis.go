package distlock

import (
	"errors"
	"time"

	"github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
)

type RedisLock struct {
	redisClient *redis.Client
	timeout     time.Duration
	resLocks    map[string]*lock.Locker
}

func NewRedisLock(redisAddr string, timeout time.Duration) *RedisLock {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	return &RedisLock{
		redisClient: redisClient,
		timeout:     timeout,
		resLocks:    make(map[string]*lock.Locker),
	}
}

func (l *RedisLock) Lock(resName string) error {
	if l.resLocks[resName] != nil {
		ok, err := l.resLocks[resName].Lock()
		if err != nil {
			return err
		}
		if ok != true {
			return errors.New("Failed to renew lock")
		}
	}
	opts := &lock.Options{
		LockTimeout: l.timeout,
	}
	lock, err := lock.Obtain(l.redisClient, resName, opts)
	if err != nil {
		return err
	}
	if lock == nil {
		return errors.New("Failed to obtain lock")
	}
	l.resLocks[resName] = lock
	return nil
}

func (l *RedisLock) Unlock(resName string) error {
	if l.resLocks[resName] == nil {
		return errors.New("Lock not found")
	}
	l.resLocks[resName].Unlock()
	l.resLocks[resName] = nil
	return nil
}
