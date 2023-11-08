package repository

import (
	"context"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/story/models/request"
)

type reportRepository struct {
	db *sqlx.DB
}

type ReportRepository interface {
	ReportPost(ctx context.Context, report request.Report) error
}

const (
	ReportPost = "INSERT INTO post_reports(post_id, reported_by, reason) values ($1, $2, $3)"
)

func (repository reportRepository) ReportPost(ctx context.Context, report request.Report) error {
	logger := logging.GetLogger(ctx).WithField("class", "ReportRepository").WithField("method", "ReportPost")

	_, err := repository.db.ExecContext(ctx, ReportPost, report.PostID, report.UserID, report.Reason)
	if err != nil {
		logger.Errorf("unable to report post %v", err)
		return err
	}

	logger.Info("successfully reported post")
	return nil
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return reportRepository{db: db}
}
