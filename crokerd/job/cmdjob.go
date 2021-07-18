package job

import (
	"io"
	"os"
	"os/exec"

	"github.com/urie96/croker/crokerd/consts"
)

type CmdJob struct {
	command string
	cmd     *exec.Cmd
	w       io.Writer
	done    chan struct{}
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

func NewCmdJob(command string) *CmdJob {
	return &CmdJob{
		command: command,
		w:       os.Stdout,
	}
}

func (s *CmdJob) SetWriter(w io.Writer) {
	if w != nil {
		s.w = w
	}
}

func (c *CmdJob) Start() error {
	if toBool(c.Done()) {
		return consts.HasStart
	}
	go func() {
		c.done = make(chan struct{})
		err := c.execCmd()
		if err != nil {
			io.WriteString(c.w, err.Error())
		}
		close(c.done)
	}()
	return nil
}

func (c *CmdJob) Stop() error {
	if !toBool(c.Done()) {
		return consts.HasStop
	}
	err := c.cmd.Process.Kill()
	close(c.done)
	return err
}

func (c *CmdJob) execCmd() error {
	cmd := exec.Command("sh", "-c", c.command)
	cmd.Dir, _ = os.UserHomeDir()
	c.cmd = cmd

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = io.Copy(c.w, io.MultiReader(stdout, stderr))
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func (c *CmdJob) Done() <-chan struct{} {

	if c.done == nil {
		return closedchan
	}
	return c.done
}

func toBool(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}
