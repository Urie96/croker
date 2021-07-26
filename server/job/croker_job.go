package job

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

type crokerJob struct {
	done chan struct{}
}

func newCrokerJob() *crokerJob {
	return &crokerJob{
		done: make(chan struct{}),
	}
}

func (*crokerJob) Start() error {
	return nil
}

func (c *crokerJob) Stop() error {
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTERM)
	return nil
}

func (*crokerJob) SetWriter(w io.Writer) {
	os.Stdout = w.(*os.File)
	os.Stderr = os.Stdout
	fmt.Println(12)
}

func (c *crokerJob) Done() <-chan struct{} {
	return c.done
}
