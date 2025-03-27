package go_fanotify

type Event struct {
	Path      string `json:"path"`
	PID       int32  `json:"pid"`
	Timestamp int64  `json:"timestamp"`
	Opened    bool   `json:"opened"`   // FAN_OPEN
	Accessed  bool   `json:"accessed"` // FAN_ACCESS
	Modified  bool   `json:"modified"` // FAN_MODIFY
	Closed    bool   `json:"closed"`   // FAN_CLOSE_WRITE/FAN_CLOSE_NOWRITE
	Created   bool   `json:"created"`  // FAN_CREATE
	Deleted   bool   `json:"deleted"`  // FAN_DELETE
	Moved     bool   `json:"moved"`    // FAN_MOVE

	Executable bool `json:"executable"` // FAN_OPEN_EXEC
	Permitted  bool `json:"permitted"`  // FAN_OPEN_PERM

	Err error `json:"-"`
}
