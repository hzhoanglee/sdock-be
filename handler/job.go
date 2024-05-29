package handler

import (
	"app/database"
	"app/model"
	"time"
)

var jobQueue = make(chan model.Job, 100) // Buffered channel with capacity 100

func init() {
	//go ScanJob()
	//go ExecuteJobs()
}

func ScanJob() {
	for {
		time.Sleep(1 * time.Second)
		jobs := []model.Job{}
		db := database.DB

		// Retrieve pending jobs from the database
		if err := db.Where("status = 0").Find(&jobs).Error; err != nil {
			// Handle error
			continue
		}

		// Enqueue jobs to the FIFO queue
		for _, job := range jobs {
			jobQueue <- job
		}
	}
}

func ExecuteJobs() {
	for job := range jobQueue {
		executeJob(job)
	}
}

func executeJob(job model.Job) {
	// Implement job execution logic here
	// You can use the job.Task field to determine the task to be executed
	// Update the job status and retry count in the database after execution
	// ...
}
