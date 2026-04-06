package repository

import (
	"database/sql"
)

type CronRepository struct {
	DB *sql.DB
}

func NewCronRepository(db *sql.DB) *CronRepository {
	return &CronRepository{DB: db}
}

func (r *CronRepository) CreateLog(job, status, message string) error {
	query := `
	INSERT INTO cron_jobs (job_name, status, message) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(
		query,
		job,
		status,
		message,
	)
	return err
}
