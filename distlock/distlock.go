package distlock

type DistLock interface {
	Lock(resName string) error
	Unlock(resName string) error
}
