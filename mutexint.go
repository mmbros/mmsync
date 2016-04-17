package mmsync

// A MutexInt represents an object that can be locked and unlocked by int.
type MutexInt interface {
	Lock(int)
	Unlock(int)
}
