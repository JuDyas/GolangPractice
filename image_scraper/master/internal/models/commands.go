package models

type CommandType string

const (
	CmdStart  CommandType = "start"
	CmdStop   CommandType = "stop"
	CmdPause  CommandType = "pause"
	CmdResume CommandType = "resume"
	CmdDone   CommandType = "done"
)

type CommandMessage struct {
	Command CommandType `json:"command"`
	URL     string      `json:"url,omitempty"`
}
