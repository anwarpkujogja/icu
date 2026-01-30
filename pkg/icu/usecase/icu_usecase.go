package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hamilton/icu-app/pkg/domain"
)

type icuUseCase struct {
	repo           domain.ICURepository
	contextTimeout time.Duration
}

func NewICUUseCase(r domain.ICURepository, timeout time.Duration) domain.ICUUseCase {
	return &icuUseCase{
		repo:           r,
		contextTimeout: timeout,
	}
}

func (u *icuUseCase) SearchPasien(c context.Context, kodeReg string) (domain.SearchResult, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.GetPasienByKodeReg(ctx, kodeReg)
}

func (u *icuUseCase) SubmitResult(c context.Context, data domain.ResultSubmission) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if len(data.Patient) == 0 {
		return fmt.Errorf("patient data cannot be empty")
	}
	if len(data.Result) == 0 {
		return fmt.Errorf("result data cannot be empty")
	}

	return u.repo.SaveResult(ctx, data)
}

func (u *icuUseCase) GetLogs(c context.Context, date string, limit int) ([]domain.AppLog, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.GetLogs(ctx, date, limit)
}

func (u *icuUseCase) SaveLog(c context.Context, log domain.AppLog) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.SaveLog(ctx, log)
}

func (u *icuUseCase) GetReport(c context.Context) ([]domain.LabResult, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.GetReport(ctx)
}

func (u *icuUseCase) RegisterAdmission(c context.Context, req domain.AdmissionRequest) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.RegisterAdmission(ctx, req)
}

func (u *icuUseCase) GetPatients(c context.Context, limit int) ([]domain.Pasien, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.repo.GetPatients(ctx, limit)
}
