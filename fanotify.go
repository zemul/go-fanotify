package go_fanotify

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	FAN_OPEN          = unix.FAN_OPEN
	FAN_ACCESS        = unix.FAN_ACCESS
	FAN_MODIFY        = unix.FAN_MODIFY
	FAN_CLOSE_WRITE   = unix.FAN_CLOSE_WRITE
	FAN_CLOSE_NOWRITE = unix.FAN_CLOSE_NOWRITE
	FAN_CREATE        = unix.FAN_CREATE
	FAN_DELETE        = unix.FAN_DELETE
	FAN_MOVED_FROM    = unix.FAN_MOVED_FROM
	FAN_MOVED_TO      = unix.FAN_MOVED_TO
	FAN_MOVE_SELF     = unix.FAN_MOVE_SELF
	FAN_OPEN_EXEC     = unix.FAN_OPEN_EXEC
	FAN_OPEN_PERM     = unix.FAN_OPEN_PERM
)

type Event struct {
	Path string `json:"path"`
	PID  int32  `json:"pid"`
	Mask uint64 `json:"mask"`
	Err  error  `json:"-"`
}

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
func (n *Notifier) AddWatch(paths []string, mask uint64) error {
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

		err = unix.FanotifyMark(n.fd,
			unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
			mask,
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
	ch := make(chan Event, 4096)
	go func() {
		defer close(ch)
		buf := make([]byte, 16*1024)

		for {
			readLen, err := unix.Read(n.fd, buf)
			if err != nil {
				ch <- Event{Err: err}
				return
			}

			offset := 0
			for offset < readLen {
				meta := (*unix.FanotifyEventMetadata)(unsafe.Pointer(&buf[offset]))

				if meta.Event_len == 0 {
					break
				}

				if meta.Vers != unix.FANOTIFY_METADATA_VERSION {
					ch <- Event{Err: fmt.Errorf("metadata version mismatch")}
					break
				}

				if meta.Mask&unix.FAN_Q_OVERFLOW != 0 {
					ch <- Event{Err: fmt.Errorf("fanotify queue overflow")}
					offset += int(meta.Event_len)
					continue
				}

				// 构造事件
				event := Event{
					PID:  meta.Pid,
					Mask: meta.Mask,
				}

				if meta.Fd >= 0 {
					event.Path = getPathFromFD(meta.Fd)
					unix.Close(int(meta.Fd))
				}

				for _, prefix := range n.path {
					if isUnderTargetDir(event.Path, prefix) {
						ch <- event
						break
					}
				}

				offset += int(meta.Event_len)
			}

		}
	}()
	return ch
}

func (n *Notifier) Close() error {
	return unix.Close(n.fd)
}
