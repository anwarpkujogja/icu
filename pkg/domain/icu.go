package domain

import (
	"context"
	"time"
)

// Entities

type Pasien struct {
	NoRM         string `json:"no_rm" db:"no_rm"`
	NamaPasien   string `json:"nama_pasien" db:"nama_pasien"`
	JenisKelamin string `json:"jenis_kelamin" db:"jenis_kelamin"`
	TanggalLahir string `json:"tanggal_lahir" db:"tanggal_lahir"`
}

type Kunjungan struct {
	KodeReg     string `json:"kode_reg" db:"kode_reg"`
	NoRM        string `json:"no_rm" db:"no_rm"`
	NoKunjungan string `json:"no_kunjungan" db:"no_kunjungan"`
	PoliRuang   string `json:"poli_ruang" db:"poli_ruang"`
	Kamar       string `json:"kamar" db:"kamar"`
	Bed         string `json:"bed" db:"bed"`
}

// SearchResult combines Pasien and Kunjungan for the /search response
type SearchResult struct {
	NoRM         string `json:"no_rm"`
	NamaPasien   string `json:"nama_pasien"`
	JenisKelamin string `json:"jenis_kelamin"`
	TanggalLahir string `json:"tanggal_lahir"`
	PoliRuang    string `json:"poli_ruang"`
	Kamar        string `json:"kamar"`
	Bed          string `json:"bed"`
	NoKunjungan  string `json:"no_kunjungan"`
}

// ResultSubmission represents the payload for POST /result
type PatientData struct {
	PID    string `json:"pid"`
	VID    string `json:"vid"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type ResultSubmission struct {
	Patient map[string]interface{} `json:"patient"`
	Result  map[string]interface{} `json:"result"` // Varied format
}

type AppLog struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Endpoint  string    `json:"endpoint" db:"endpoint"`
	Method    string    `json:"method" db:"method"`
	Status    int       `json:"status" db:"status"`
	Message   string    `json:"message" db:"message"`
}

type LabResult struct {
	ID          int64                  `json:"id" db:"id"`
	PatientData map[string]interface{} `json:"patient_data" db:"patient_data"`
	ResultData  map[string]interface{} `json:"result_data" db:"result_data"`
}

type AdmissionRequest struct {
	NoRM         string `json:"no_rm"`
	NamaPasien   string `json:"nama_pasien"`
	JenisKelamin string `json:"jenis_kelamin"`
	TanggalLahir string `json:"tanggal_lahir"`
	PoliRuang    string `json:"poli_ruang"`
	Kamar        string `json:"kamar"`
	Bed          string `json:"bed"`
	NoKunjungan  string `json:"no_kunjungan"`
}

// Interfaces

type ICUUseCase interface {
	SearchPasien(ctx context.Context, kodeReg string) (SearchResult, error)
	SubmitResult(ctx context.Context, data ResultSubmission) error
	GetLogs(ctx context.Context, date string, limit int) ([]AppLog, error)
	SaveLog(ctx context.Context, log AppLog) error
	GetReport(ctx context.Context) ([]LabResult, error)
	RegisterAdmission(ctx context.Context, req AdmissionRequest) (string, error)
	GetPatients(ctx context.Context, limit int) ([]Pasien, error)
}

type ICURepository interface {
	GetPasienByKodeReg(ctx context.Context, kodeReg string) (SearchResult, error)
	SaveResult(ctx context.Context, data ResultSubmission) error
	GetLogs(ctx context.Context, date string, limit int) ([]AppLog, error)
	SaveLog(ctx context.Context, log AppLog) error
	GetReport(ctx context.Context) ([]LabResult, error)
	RegisterAdmission(ctx context.Context, req AdmissionRequest) (string, error)
	GetPatients(ctx context.Context, limit int) ([]Pasien, error)
}
