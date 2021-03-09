package util

import "sync/atomic"

var value uint32 = 0

/**
return 0 ~ ox7fffffff
本地自增序列号
*/
func GetSeq() uint32 {
	var newInt uint32
	newInt = atomic.AddUint32(&value, 1)
	return newInt & 0x7fffffff
}