package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hamilton/icu-app/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	Conn *sqlx.DB
}

func NewPostgresRepository(conn *sqlx.DB) domain.ICURepository {
	return &PostgresRepository{Conn: conn}
}

func (p *PostgresRepository) GetPasienByKodeReg(ctx context.Context, kodeReg string) (domain.SearchResult, error) {
	query := `
		SELECT 
			p.no_rm, p.nama_pasien, p.jenis_kelamin, p.tanggal_lahir,
			k.poli_ruang, k.kamar, k.bed, k.no_kunjungan
		FROM kunjungan k
		JOIN pasiens p ON k.no_rm = p.no_rm
		WHERE k.kode_reg = $1
	`
	var result domain.SearchResult
	// sqlx allows scanning directly to struct
	// Note: We need to handle cases where struct fields mismatch DB columns if not tagged carefully,
	// but standard scan works if names match or use 'db' tag.
	// We didn't use 'db' tags on SearchResult, so let's use standard QueryRow or StructScan if tags were present.
	// Let's rely on standard Scan for safety or add tags in domain if needed.
	// I added db tags to entities but SearchResult doesn't have them in my previous step?
	// Wait, I didn't add db tags to SearchResult in domain file.
	// Let's use QueryRow and Scan manually or assume loose matching (which sqlx does).

	// Actually, I should probably use `Get` or `StructScan`.
	// Let's use `Get` but I need to make sure the query columns match lowercase snake_case of struct fields or db tags.
	// In Clean Arch, repository implementation knows about DB structure.

	row := p.Conn.QueryRowContext(ctx, query, kodeReg)
	err := row.Scan(
		&result.NoRM, &result.NamaPasien, &result.JenisKelamin, &result.TanggalLahir,
		&result.PoliRuang, &result.Kamar, &result.Bed, &result.NoKunjungan,
	)

	return result, err
}

func (p *PostgresRepository) SaveResult(ctx context.Context, data domain.ResultSubmission) error {
	query := `INSERT INTO hasil_lab (patient_data, result_data) VALUES ($1::jsonb, $2::jsonb)`

	pData, err := json.MarshalIndent(data.Patient, "", "  ")
	if err != nil {
		return err
	}
	rData, err := json.MarshalIndent(data.Result, "", "  ")
	if err != nil {
		return err
	}

	// DEBUG LOG
	fmt.Printf("DEBUG: Inserting into hasil_lab. Patient=%s Result=%s\n", string(pData), string(rData))

	// lib/pq treats []byte as bytea. For JSONB, we must pass it as string.
	// DEBUG: Print info before exec
	_, err = p.Conn.ExecContext(ctx, query, string(pData), string(rData))
	if err != nil {
		fmt.Printf("DEBUG: Insert Operation Failed: %v\n", err)
		return err
	}
	
	fmt.Println("DEBUG: Insert Operation Reported Success.")

	// DEBUG: Check row count immediately
	var count int
	_ = p.Conn.QueryRowContext(ctx, "SELECT count(*) FROM hasil_lab").Scan(&count)
	fmt.Printf("DEBUG: Current Total Rows in 'hasil_lab': %d\n", count)
	return nil
}

func (p *PostgresRepository) GetLogs(ctx context.Context, date string, limit int) ([]domain.AppLog, error) {
	// Date format Y-m-d, assume filtering by created_at date part
	query := `
		SELECT id, created_at, endpoint, method, status, message
		FROM app_logs
		WHERE DATE(created_at) = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := p.Conn.QueryContext(ctx, query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.AppLog
	for rows.Next() {
		var l domain.AppLog
		if err := rows.Scan(&l.ID, &l.CreatedAt, &l.Endpoint, &l.Method, &l.Status, &l.Message); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (p *PostgresRepository) SaveLog(ctx context.Context, log domain.AppLog) error {
	query := `INSERT INTO app_logs (endpoint, method, status, message, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.Conn.ExecContext(ctx, query, log.Endpoint, log.Method, log.Status, log.Message, time.Now())
	return err
}

func (p *PostgresRepository) GetReport(ctx context.Context) ([]domain.LabResult, error) {
	query := `SELECT id, patient_data, result_data FROM hasil_lab ORDER BY id DESC`
	rows, err := p.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.LabResult
	for rows.Next() {
		var r domain.LabResult
		var pData, rData []byte
		if err := rows.Scan(&r.ID, &pData, &rData); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(pData, &r.PatientData); err != nil {
			// If unmarshal fails, we might want to log it or skip, but for now allow returning error
			return nil, err
		}
		if err := json.Unmarshal(rData, &r.ResultData); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (p *PostgresRepository) RegisterAdmission(ctx context.Context, req domain.AdmissionRequest) (string, error) {
	// Generate UUID for kode_reg
	kodeReg := fmt.Sprintf("REG-%d", time.Now().UnixNano()) // Simple ID generation

	tx, err := p.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 1. Check Pasien Exists
	var count int
	err = tx.QueryRowContext(ctx, "SELECT count(*) FROM pasiens WHERE no_rm = $1", req.NoRM).Scan(&count)
	if err != nil {
		return "", err
	}

	// 2. Insert Pasien if not exists
	if count == 0 {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pasiens (no_rm, nama_pasien, jenis_kelamin, tanggal_lahir)
			VALUES ($1, $2, $3, $4)
		`, req.NoRM, req.NamaPasien, req.JenisKelamin, req.TanggalLahir)
		if err != nil {
			return "", err
		}
	}

	// 3. Insert Kunjungan
	_, err = tx.ExecContext(ctx, `
		INSERT INTO kunjungan (kode_reg, no_rm, no_kunjungan, poli_ruang, kamar, bed)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, kodeReg, req.NoRM, req.NoKunjungan, req.PoliRuang, req.Kamar, req.Bed)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return kodeReg, nil
}

func (p *PostgresRepository) GetPatients(ctx context.Context, limit int) ([]domain.Pasien, error) {
	query := `
		SELECT no_rm, nama_pasien, jenis_kelamin, tanggal_lahir
		FROM pasiens
		ORDER BY no_rm ASC
		LIMIT $1
	`
	rows, err := p.Conn.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []domain.Pasien
	for rows.Next() {
		var pat domain.Pasien
		if err := rows.Scan(&pat.NoRM, &pat.NamaPasien, &pat.JenisKelamin, &pat.TanggalLahir); err != nil {
			return nil, err
		}
		patients = append(patients, pat)
	}
	return patients, nil
}
