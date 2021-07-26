package http

import (
	"io"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/urie96/croker/model"
	"github.com/urie96/croker/server/job"
	"github.com/urie96/croker/utils"
	"github.com/urie96/croker/version"
)

func Handler() *gin.Engine {
	jobManager := job.NewManager()
	gin.DefaultWriter = os.Stdout // 重新更新输出位置
	gin.DefaultErrorWriter = os.Stderr
	r := gin.Default()
	r.GET("/detail", func(c *gin.Context) {
		c.JSON(200, model.GetServerDetailResponse{
			Status:  "running",
			Version: version.Version,
		})
	})
	r.GET("/jobs", func(c *gin.Context) {
		var resp model.GetJobResponse
		infos := jobManager.Infos()
		jobInfos := make([]*model.JobInfo, 0, len(infos))
		for id, info := range jobManager.Infos() {
			jobInfos = append(jobInfos, &model.JobInfo{
				ID:        id,
				Command:   info.Command,
				CronSpec:  info.CronSpec,
				IsRunning: info.IsRunning,
				CreatedAt: info.CreatedAt,
				StartedAt: info.StartedAt,
				StoppedAt: info.StoppedAt,
			})
		}
		sort.Slice(jobInfos, func(i, j int) bool {
			return jobInfos[i].CreatedAt < jobInfos[j].CreatedAt
		})
		resp.JobInfos = jobInfos
		c.JSON(200, resp)
	})
	_job := r.Group("/job")
	{
		_job.POST("/", func(c *gin.Context) {
			var req model.AddJobRequest
			var resp model.AddJobResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			if err = c.Bind(&req); err != nil {
				return
			}
			resp.JobID, err = jobManager.Add(req.CronSpec, req.Command)
		})

		_job.DELETE("/:id", func(c *gin.Context) {
			var resp model.BaseResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			err = jobManager.Remove(c.Param("id"))
		})

		_job.PUT("/:id/start", func(c *gin.Context) {
			var resp model.BaseResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			err = jobManager.Start(c.Param("id"))
		})

		_job.PUT("/:id/stop", func(c *gin.Context) {
			var resp model.StopJobResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			_, jobId, _ := jobManager.GetByPrefix(c.Param("id"))
			resp.JobID = jobId
			err = jobManager.Stop(jobId)
		})

		_job.GET("/:id/log", func(c *gin.Context) {
			var req model.GetJobLogRequest
			var err error
			_ = c.Bind(&req)
			id := c.Param("id")
			var r io.ReadCloser
			if req.Follow {
				r, err = jobManager.FollowLog(c.Request.Context(), id)
			} else {
				r, err = jobManager.Log(id)
			}
			if err != nil {
				c.String(400, err.Error())
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Flush()
			w := utils.NewWriteFlusher(c.Writer)
			io.Copy(w, r)
		})
	}
	return r
}
