package service

import (
	"context"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"post-api/idp/constants"
	"post-api/story/models/request"
	"post-api/story/repository"
)

type reportService struct {
	repository repository.ReportRepository
}

type ReportService interface {
	ReportPost(ctx context.Context, report request.Report) *golaerror.Error
}

func (service reportService) ReportPost(ctx context.Context, report request.Report) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "ReportService").WithField("method", "ReportPost")
	err := service.repository.ReportPost(ctx, report)
	if err != nil {
		logger.Errorf("unable to report post %v", err)
		return &constants.InternalServerError
	}
	logger.Info("successfully reported post")

	return nil
}

func NewReportService(reportRepository repository.ReportRepository) ReportService {
	return reportService{
		repository: reportRepository,
	}
}
