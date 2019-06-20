package channel

import (
	"github.com/okpub/rhino/errors"
	"github.com/okpub/rhino/process"
)

//class mailbox(使用消息通道代替)
type Mailbox struct {
	process.UntypeProcess
	//message
	buffer      chan interface{} //消息通道
	pendingNum  int              //默认通道
	nonblocking bool             //模式(默认阻塞)
}

func (this *Mailbox) Init(opts ...Option) *Mailbox {
	for _, o := range opts {
		o(this)
	}
	if this.buffer == nil {
		this.buffer = make(chan interface{}, this.pendingNum)
	}
	return this
}

func (this Mailbox) Copy(opts ...Option) *Mailbox {
	return this.Init(opts...)
}

//process
func (this *Mailbox) Start() (err error) {
	this.OnStarted()
	this.Schedule(this.run)
	return
}

func (this *Mailbox) Close() (err error) {
	defer func() { err = errors.Catch(recover()) }()
	close(this.buffer)
	return
}

func (this *Mailbox) run() {
	var (
		body  interface{}
		err   error
		debug = false
	)
	defer func() {
		if debug {
			if err = errors.Catch(recover()); err != nil {
				this.ThrowFailure(err, body)
			}
		}
		this.Close()
		this.PostStop()
	}()
	this.PreStart()
	//run
	for body = range this.buffer {
		//First of all statistics (processing failure will also record)
		this.OnReceived(body)
		//process message
		this.DispatchMessage(body)
	}
}

func (this *Mailbox) Post(v interface{}) (err error) {
	this.OnPosted(v)
	return errors.Try(func() error {
		return this.sendMessage(v)
	}, func(err error) {
		this.OnDiscarded(err, v)
	})
}

//private
func (this *Mailbox) sendMessage(v interface{}) (err error) {
	if this.nonblocking {
		select {
		case this.buffer <- v:
		default:
			err = OverfullErr
		}
	} else {
		this.buffer <- v
	}
	return
}
