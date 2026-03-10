package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/hooneun/scorpes/internal/job"
	"github.com/hooneun/scorpes/internal/worker"
	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cron     *cron.Cron
	jobQueue worker.JobQueue
	cfg      *config.Config
	db       *db.Queries
}

func NewCronScheduler(queue worker.JobQueue, cfg *config.Config, db *db.Queries) *CronScheduler {
	c := cron.New(
		// seconds 지원
		cron.WithSeconds(),
		cron.WithChain(
			// 이전 job 끝나기 전에 실행 방지
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	return &CronScheduler{
		cron:     c,
		jobQueue: queue,
		cfg:      cfg,
		db:       db,
	}
}

type CronInfo struct {
	ID          string
	ExecuteTime time.Time
}

func (s *CronScheduler) Start() {
	ctx := context.Background()
	// 매 분 마다 실행
	/**
	┌──────── second
	│ ┌────── minute
	│ │ ┌──── hour
	│ │ │ ┌── day
	│ │ │ │ ┌ month
	│ │ │ │ │ ┌ weekday
	│ │ │ │ │ │
	* * * * * *
	*/

	_, err := s.cron.AddFunc("0 */1 * * * *", func() {
		s.jobQueue <- func() {
			targets, err := s.db.ListTargets(ctx)
			if err != nil {
				log.Fatal(err)
			}

			var cronList []CronInfo
			for _, target := range targets {
				cronList = append(cronList, CronInfo{
					ID:          target.ID.String(),
					ExecuteTime: time.Now(),
				})
			}

			job.HealthCheck()
			log.Println("cron job executed")
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	s.cron.Start()
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
}
