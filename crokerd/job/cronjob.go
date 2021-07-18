package job

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/urie96/croker/crokerd/consts"
)

type CronJob struct {
	cronSpec string
	cronJob  *cron.Cron
	cmdJob   *CmdJob
	w        io.Writer
	done     chan struct{}
}

func NewCronJob(cronspec, command string) *CronJob {
	c := &CronJob{
		cronSpec: cronspec,
		cmdJob:   NewCmdJob(command),
		w:        os.Stdout,
	}
	return c
}

func (c *CronJob) SetWriter(w io.Writer) {
	c.cmdJob.SetWriter(w)
	c.w = w
}

func (c *CronJob) Start() error {
	if toBool(c.Done()) {
		return consts.HasStart
	}
	c.done = make(chan struct{})
	c.cronJob = cron.New()
	_, err := c.cronJob.AddFunc(c.cronSpec, func() {
		io.WriteString(c.w, fmt.Sprintf("\n%s  :\n", time.Now().Format("2006-01-02 15:04:05")))
		if err := c.cmdJob.Start(); err != nil {
			io.WriteString(c.w, err.Error())
		}
	})
	if err != nil {
		return err
	}
	c.cronJob.Start()
	return nil
}

func (c *CronJob) Stop() error {
	if !toBool(c.Done()) {
		return consts.HasStop
	}
	c.cronJob.Stop()
	close(c.done)
	return nil
}

func (c *CronJob) Done() <-chan struct{} {
	if c.done == nil {
		return closedchan
	}
	return c.done
}
