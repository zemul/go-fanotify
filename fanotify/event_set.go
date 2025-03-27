package fanotify

import (
	"golang.org/x/sys/unix"
	"strings"
)

type EventType struct {
	mask uint64
	name string
}

var (
	eventFileWriteComplete = EventType{unix.FAN_CLOSE_WRITE, "FileWriteComplete"}
	eventFileModified      = EventType{unix.FAN_MODIFY, "FileModified"}
	eventFileOpened        = EventType{unix.FAN_OPEN, "FileOpened"}
	eventFileAccessed      = EventType{unix.FAN_ACCESS, "FileAccessed"}
	eventDirCreated        = EventType{unix.FAN_CREATE, "DirCreated"}
	eventDirMoved          = EventType{unix.FAN_MOVED_TO, "DirMoved"}
	eventDirDeleted        = EventType{unix.FAN_DELETE, "DirDeleted"}
)

func (e EventType) Mask() uint64   { return e.mask }
func (e EventType) String() string { return e.name }

var (
	FileWriteComplete = eventFileWriteComplete
	FileModified      = eventFileModified
	FileOpened        = eventFileOpened
	FileAccessed      = eventFileAccessed
	DirCreated        = eventDirCreated
	DirMoved          = eventDirMoved
	DirDeleted        = eventDirDeleted
)

type EventSet struct {
	mask  uint64
	names []string
}

func NewEventSet(events ...EventType) EventSet {
	es := EventSet{}
	for _, e := range events {
		es.mask |= e.mask
		es.names = append(es.names, e.name)
	}
	return es
}

func (es EventSet) Mask() uint64   { return es.mask }
func (es EventSet) String() string { return strings.Join(es.names, "|") }

var (
	WriteComplete = NewEventSet(FileWriteComplete, FileModified)
	AllOps        = NewEventSet(FileWriteComplete, FileModified, FileOpened, DirCreated)
)
