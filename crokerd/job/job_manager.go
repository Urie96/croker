package job

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/urie96/croker/crokerd/consts"
	"github.com/urie96/croker/utils"
)

var manager = NewManager()

func init() {
	panicIfErr(manager.load())

	c := make(chan os.Signal, 1)
	os.Interrupt.Signal()
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		manager.sync()
		manager.stopAll() // 避免僵尸进程
		os.Exit(0)
	}()
}

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
	return &JobManager{
		jobs:     map[string]Job{},
		jobinfos: map[string]*JobInfo{},
	}
}

func Manager() *JobManager {
	return manager
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
	context.Background()
	return job, id, nil
}

func (j *JobManager) Infos() map[string]*JobInfo {
	infos := j.jobinfos
	for id, info := range infos {
		info.IsRunning = toBool(j.jobs[id].Done())
	}
	return infos
}

func (j *JobManager) add(id, cronspec, command string) {
	var job Job
	if cronspec == "" {
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
	file, err := os.OpenFile(consts.JobLogPath+"/"+id, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
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
	logPath := consts.JobLogPath + "/" + id
	return os.OpenFile(logPath, os.O_RDONLY, 0777)
}

func (j *JobManager) FollowLog(ctx context.Context, prefix string) (io.ReadCloser, error) {
	job, id, err := j.GetByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	if !toBool(job.Done()) {
		return j.Log(prefix)
	}
	logPath := consts.JobLogPath + "/" + id
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
	logPath := consts.JobLogPath + "/" + id
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

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (j *JobManager) load() error {
	if len(j.jobs) > 0 {
		return consts.JobNotEmpty
	}
	file, err := os.OpenFile(consts.JobsSyncPath, os.O_RDONLY, 0777)
	if err != nil {
		return nil
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	infos := map[string]*JobInfo{}
	_ = json.Unmarshal(b, &infos)
	for id, info := range infos {
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

func (j *JobManager) sync() {
	file, err := os.OpenFile(consts.JobsSyncPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	panicIfErr(err)
	defer file.Close()
	infos := j.Infos()
	data, err := json.Marshal(infos)
	panicIfErr(err)
	_, err = file.WriteString(string(data))
	panicIfErr(err)
}

func (j *JobManager) stopAll() {
	for _, job := range j.jobs {
		job.Stop()
	}
}
