// +build darwin dragonfly freebsd netbsd openbsd

package logs

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

func TryToDropFilePageCache(fd int, offset int64, length int64) {
}
