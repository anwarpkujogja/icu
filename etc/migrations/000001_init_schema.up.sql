CREATE TABLE pasiens (
    no_rm VARCHAR(50) PRIMARY KEY,
    nama_pasien VARCHAR(255) NOT NULL,
    jenis_kelamin VARCHAR(10),
    tanggal_lahir DATE
);

CREATE TABLE kunjungan (
    kode_reg VARCHAR(50) PRIMARY KEY,
    no_rm VARCHAR(50) NOT NULL REFERENCES pasiens(no_rm),
    no_kunjungan VARCHAR(50),
    poli_ruang VARCHAR(100),
    kamar VARCHAR(50),
    bed VARCHAR(50)
);

CREATE TABLE hasil_lab (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    patient_data JSONB,
    result_data JSONB
);

CREATE TABLE app_logs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    endpoint VARCHAR(100),
    method VARCHAR(10),
    status INT,
    message TEXT
);

-- Seed Data (Optional, for testing)
INSERT INTO pasiens (no_rm, nama_pasien, jenis_kelamin, tanggal_lahir) VALUES 
('RM001', 'John Doe', 'L', '1980-01-01'),
('RM002', 'Jane Smith', 'P', '1990-05-15');

INSERT INTO kunjungan (kode_reg, no_rm, no_kunjungan, poli_ruang, kamar, bed) VALUES
('REG001', 'RM001', 'K001', 'ICU', '101', '1'),
('REG002', 'RM002', 'K002', 'UGD', '202', '2');
