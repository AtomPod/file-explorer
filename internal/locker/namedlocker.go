package locker

//NamedLocker 锁接口
type NamedLocker interface {
	Lock(string)
	UnLock(string)
}
