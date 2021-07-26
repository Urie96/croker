package model

type AddJobRequest struct {
	Command  string `json:"command" form:"command" binding:"required"`
	CronSpec string `json:"cronSpec" form:"cronSpec"`
}

type AddJobResponse struct {
	BaseResponse
	JobID string `json:"jobId"`
}

type GetJobRequest struct {
}

type GetJobResponse struct {
	BaseResponse
	JobInfos []*JobInfo `json:"jobInfos"`
}
type StopJobResponse struct {
	BaseResponse
	JobID string `json:"jobId"`
}

type JobInfo struct {
	ID        string `json:"id"`
	Command   string `json:"command"`
	CronSpec  string `json:"cronSpec"`
	IsRunning bool   `json:"isRunning"`
	CreatedAt int64  `json:"createdAt"`
	StartedAt int64  `json:"started"`
	StoppedAt int64  `json:"stoppedAt"`
}

type GetJobLogRequest struct {
	Follow bool `json:"follow" form:"follow"`
}

type ServerVersionResponse struct {
	BaseResponse
	Version string `json:"version"`
}

type GetServerDetailResponse struct {
	BaseResponse
	Status  string `json:"status"`
	Version string `json:"version"`
}
