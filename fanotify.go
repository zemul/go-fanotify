package go_fanotify

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	EventFileWriteComplete = unix.FAN_CLOSE_WRITE // 文件写入完成
	EventFileModified      = unix.FAN_MODIFY      // 文件内容修改
	EventFileOpened        = unix.FAN_OPEN        // 文件打开
	EventFileAccessed      = unix.FAN_ACCESS      // 文件读取

	// 目录操作事件
	EventDirCreated = unix.FAN_CREATE   // 目录/文件创建
	EventDirMoved   = unix.FAN_MOVED_TO // 目录/文件移动至
	EventDirDeleted = unix.FAN_DELETE   // 目录/文件删除

	// 组合事件
	EventAllWrites = EventFileWriteComplete | EventFileModified
	EventAllOps    = EventFileWriteComplete | EventFileModified | EventFileOpened | EventDirCreated
)

type Notifier struct {
	mounted map[string]bool
	path    []string
	fd      int
}

func New() (*Notifier, error) {
	fd, err := unix.FanotifyInit(
		unix.FAN_CLASS_NOTIF|unix.FAN_CLOEXEC,
		unix.O_RDONLY|unix.O_LARGEFILE,
	)
	if err != nil {
		return nil, fmt.Errorf("fanotify_init failed: %w", err)
	}
	return &Notifier{fd: fd, mounted: make(map[string]bool)}, nil
}

// AddWatch Adding a monitoring path
func (n *Notifier) AddWatch(paths []string, events EventSet) error {
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		n.path = append(n.path, absPath)
		mountPoint := getMountPoint(absPath)
		if n.mounted[mountPoint] {
			continue
		}

		err = unix.FanotifyMark(fd,
			unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
			events.Mask(),
			unix.AT_FDCWD,
			mountPoint)
		if err != nil {
			return fmt.Errorf("fanotify_mark failed: %w", err)
		}
		n.mounted[mountPoint] = true
	}
	return nil
}

func (n *Notifier) ReadEvents() <-chan Event {
	ch := make(chan Event, 100)
	go func() {
		defer close(ch)
		buf := make([]byte, 4096)

		for {
			n, err := unix.Read(n.fd, buf)
			if err != nil {
				ch <- Event{Err: err}
				return
			}

			meta := (*unix.FanotifyEventMetadata)(unsafe.Pointer(&buf[0]))
			if meta.Vers != unix.FANOTIFY_METADATA_VERSION {
				ch <- Event{Err: fmt.Errorf("metadata version mismatch")}
				continue
			}
			event := parseMask(meta.Mask)
			event.PID = meta.Pid
			if meta.Fd >= 0 {
				event.Path = getPathFromFD(meta.Fd)
				unix.Close(meta.Fd)
			}

			ch <- event
		}
	}()
	return ch
}

func (n *Notifier) Close() error {
	return unix.Close(n.fd)
}

func parseMask(mask uint64) Event {
	return Event{
		Opened:     mask&unix.FAN_OPEN != 0,
		Accessed:   mask&unix.FAN_ACCESS != 0,
		Modified:   mask&unix.FAN_MODIFY != 0,
		Closed:     mask&(unix.FAN_CLOSE_WRITE|unix.FAN_CLOSE_NOWRITE) != 0,
		Created:    mask&unix.FAN_CREATE != 0,
		Deleted:    mask&unix.FAN_DELETE != 0,
		Moved:      mask&(unix.FAN_MOVED_FROM|unix.FAN_MOVED_TO) != 0,
		Executable: mask&unix.FAN_OPEN_EXEC != 0,
		Permitted:  mask&unix.FAN_OPEN_PERM != 0,
	}
}
