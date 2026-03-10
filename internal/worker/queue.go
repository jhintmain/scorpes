package worker

type Job func()
type JobQueue chan Job
