package http

import (
	"io"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/urie96/croker/crokerd/job"
	"github.com/urie96/croker/model"
	"github.com/urie96/croker/utils"
)

func Run() error {
	r := gin.Default()
	r.GET("/jobs", func(c *gin.Context) {
		var resp model.GetJobResponse
		var err error
		defer func() {
			resp.SetError(err)
			c.JSON(200, resp)
		}()
		infos := job.Manager().Infos()
		jobInfos := make([]*model.JobInfo, 0, len(infos))
		for id, info := range job.Manager().Infos() {
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
		resp.JobInfos = jobInfos
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
			resp.JobID, err = job.Manager().Add(req.CronSpec, req.Command)
		})

		_job.DELETE("/:id", func(c *gin.Context) {
			var resp model.BaseResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			err = job.Manager().Remove(c.Param("id"))
		})

		_job.PUT("/:id/start", func(c *gin.Context) {
			var resp model.BaseResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			err = job.Manager().Start(c.Param("id"))
		})

		_job.PUT("/:id/stop", func(c *gin.Context) {
			var resp model.StopJobResponse
			var err error
			defer func() {
				resp.SetError(err)
				c.JSON(200, resp)
			}()
			m := job.Manager()
			_, jobId, _ := m.GetByPrefix(c.Param("id"))
			resp.JobID = jobId
			err = m.Stop(jobId)
		})

		_job.GET("/:id/log", func(c *gin.Context) {
			var req model.GetJobLogRequest
			var err error
			_ = c.Bind(&req)
			id := c.Param("id")
			var r io.ReadCloser
			m := job.Manager()
			if req.Follow {
				r, err = m.FollowLog(c.Request.Context(), id)
			} else {
				r, err = m.Log(id)
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

	r.GET("/stop", func(c *gin.Context) {
		c.Status(200)
		os.Interrupt.Signal()
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	})

	r.GET("/version", func(c *gin.Context) {
		resp := model.ServerVersionResponse{
			Version: "0.0.1",
		}
		c.JSON(200, &resp)
	})

	return r.Run(":1996")
}
