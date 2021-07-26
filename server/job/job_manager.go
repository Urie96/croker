package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/urie96/croker/server/consts"
	"github.com/urie96/croker/utils"
)

const CROKER_ID = "main"

type JobInfo struct {
	Command   string
	CronSpec  string
	IsRunning bool
	CreatedAt int64
	StartedAt int64
	StoppedAt int64
}

type Job interface {
	Start() error
	Stop() error
	SetWriter(io.Writer)
	Done() <-chan struct{}
}

type JobManager struct {
	jobs     map[string]Job
	jobinfos map[string]*JobInfo
}

func NewManager() *JobManager {
	j := JobManager{
		jobs:     map[string]Job{},
		jobinfos: map[string]*JobInfo{},
	}
	if err := j.load(); err != nil {
		fmt.Println(err)
	}
	if err := j.add(CROKER_ID, "", "croker"); err == nil {
		j.Start(CROKER_ID)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		j.Close()
		os.Exit(0)
	}()
	return &j
}

func (j *JobManager) GetByPrefix(prefix string) (Job, string, error) {
	if prefix == "" {
		return nil, "", errNoSuchJob(prefix)
	}
	var job Job
	var id string
	for k, v := range j.jobs {
		if strings.HasPrefix(k, prefix) {
			if job != nil {
				return nil, "", errMultipleIDs(prefix)
			}
			id, job = k, v
		}
	}
	if job == nil {
		return nil, "", errNoSuchJob(prefix)
	}
	return job, id, nil
}

func (j *JobManager) Infos() map[string]*JobInfo {
	infos := j.jobinfos
	for id, info := range infos {
		info.IsRunning = toBool(j.jobs[id].Done())
	}
	return infos
}

func (j *JobManager) add(id, cronspec, command string) error {
	if _, exist := j.jobs[id]; exist {
		return errors.New("job exist")
	}
	var job Job
	if id == CROKER_ID {
		job = newCrokerJob()
	} else if cronspec == "" {
		job = NewCmdJob(command)
	} else {
		job = NewCronJob(cronspec, command)
	}
	j.jobs[id] = job
	j.jobinfos[id] = &JobInfo{
		Command:   command,
		CronSpec:  cronspec,
		CreatedAt: time.Now().Unix(),
	}
	return nil
}

func (j *JobManager) Add(cronspec, command string) (string, error) {
	id := genID()
	j.add(id, cronspec, command)
	err := j.Start(id)
	return id, err
}

func (j *JobManager) Start(prefix string) error {
	job, id, err := j.GetByPrefix(prefix)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(consts.JobLogDir+"/"+id, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	job.SetWriter(file)
	if err = job.Start(); err != nil {
		return err
	}
	j.jobinfos[id].StartedAt = time.Now().Unix()
	go func() {
		<-job.Done()
		j.jobinfos[id].StoppedAt = time.Now().Unix()
	}()
	return nil
}

func (j *JobManager) Stop(prefix string) error {
	job, _, err := j.GetByPrefix(prefix)
	if err != nil {
		return err
	}
	return job.Stop()
}

func (j *JobManager) Log(prefix string) (io.ReadCloser, error) {
	_, id, err := j.GetByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	logPath := consts.JobLogDir + "/" + id
	return os.Open(logPath)
}

func (j *JobManager) FollowLog(ctx context.Context, prefix string) (io.ReadCloser, error) {
	job, id, err := j.GetByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	if !toBool(job.Done()) {
		return j.Log(prefix)
	}
	logPath := consts.JobLogDir + "/" + id
	file, err := utils.FollowFile(logPath)
	if err != nil {
		return nil, err
	}
	go func() {
		select {
		case <-ctx.Done():
		case <-job.Done():
		}
		file.Close()
	}()
	return file, nil
}

func (j *JobManager) Remove(prefix string) error {
	job, id, err := j.GetByPrefix(prefix)
	if err != nil {
		return err
	}
	if toBool(job.Done()) {
		return errRemoveRunningJob(id)
	}
	delete(j.jobs, id)
	delete(j.jobinfos, id)
	logPath := consts.JobLogDir + "/" + id
	os.Remove(logPath)
	return nil
}

var genID = func() func() string {
	charset := []byte("0123456789abcdef")
	randsrc := rand.NewSource(time.Now().Unix())
	return func() string {
		num := randsrc.Int63()
		res := make([]byte, 12)
		for i := 0; i < len(res); i++ {
			res[i] = charset[num&15]
			num = num >> 4
		}
		return string(res)
	}
}()

func (j *JobManager) load() error {
	data, err := ioutil.ReadFile(consts.JobsSyncPath)
	if err != nil {
		return nil
	}
	infos := map[string]*JobInfo{}
	_ = json.Unmarshal(data, &infos)
	for id, info := range infos {
		if _, exist := j.jobs[id]; exist {
			continue
		}
		j.add(id, info.CronSpec, info.Command)
		j.jobinfos[id].CreatedAt = info.CreatedAt
		if info.IsRunning {
			err := j.Start(id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (j *JobManager) Close() {
	data, _ := json.Marshal(j.Infos())
	ioutil.WriteFile(consts.JobsSyncPath, data, 0777)
	for _, job := range j.jobs {
		job.Stop()
	}
}
