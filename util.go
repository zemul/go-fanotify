package go_fanotify

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"syscall"
)

func isMountPoint(path string) bool {
	var st1, st2 syscall.Stat_t
	syscall.Stat(path, &st1)
	syscall.Stat(filepath.Dir(path), &st2)
	return st1.Dev != st2.Dev
}

func getPathFromFD(fd int32) string {
	path, _ := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	return path
}

func getMountPoint(path string) string {
	var stat unix.Statfs_t
	if err := unix.Statfs(path, &stat); err != nil {
		return path
	}

	mountPoint := path
	for {
		parent := filepath.Dir(mountPoint)
		if parent == mountPoint {
			break
		}

		var parentStat unix.Statfs_t
		if err := unix.Statfs(parent, &parentStat); err != nil {
			break
		}

		if stat.Fsid != parentStat.Fsid {
			break
		}
		mountPoint = parent
	}

	return mountPoint
}
