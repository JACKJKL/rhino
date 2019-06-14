package channel

import (
	"fmt"

	"github.com/okpub/rhino/process"
)

var (
	OverfullErr = fmt.Errorf("channel overfull")
)

type Option func(*Options)

type MessageQueue interface {
	process.Process
	Options() Options
	Post(interface{}) error
}