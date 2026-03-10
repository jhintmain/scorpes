package worker

import (
	"log"

	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
)

/**
JobQueue 소비
*/

type Pool struct {
	JobQueue JobQueue
	Workers  int
	Cfg      *config.Config
	Db       *db.Queries
}

func NewPool(workerCount, queueSize int, cfg *config.Config, db *db.Queries) *Pool {
	return &Pool{
		JobQueue: make(JobQueue, queueSize),
		Workers:  workerCount,
		Cfg:      cfg,
		Db:       db,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.Workers; i++ {
		go func(id int) {
			for job := range p.JobQueue {
				log.Printf("worker %d execute job\n", id)
				job()
			}
		}(i)
	}
}
