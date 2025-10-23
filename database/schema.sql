
DROP TABLE IF EXISTS pekerjaan_alumni;
DROP TABLE IF EXISTS alumni;
DROP TABLE IF EXISTS roles;

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE alumni (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    nim VARCHAR(50),
    nama VARCHAR(255),
    jurusan VARCHAR(255),
    angkatan INT,
    tahun_lulus INT,
    no_telepon VARCHAR(50),
    alamat TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE pekerjaan_alumni (
    id SERIAL PRIMARY KEY,
    alumni_id INT NOT NULL REFERENCES alumni(id) ON DELETE CASCADE,
    nama_perusahaan VARCHAR(255) NOT NULL,
    posisi_jabatan VARCHAR(255) NOT NULL,
    bidang_industri VARCHAR(255) NOT NULL,
    lokasi_kerja VARCHAR(255) NOT NULL,
    gaji_range VARCHAR(100),
    tanggal_mulai_kerja DATE NOT NULL,
    tanggal_selesai_kerja DATE,
    status_pekerjaan VARCHAR(50) NOT NULL,
    deskripsi_pekerjaan TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_delete TIMESTAMP NULL
);

-- Add comment to explain the is_delete column purpose
COMMENT ON COLUMN pekerjaan_alumni.is_delete IS 'Timestamp when the record was soft deleted. NULL means not deleted.';

-- Create an index for better performance on queries filtering by is_delete
CREATE INDEX idx_pekerjaan_alumni_is_delete ON pekerjaan_alumni(is_delete);

INSERT INTO roles (name) VALUES ('admin'), ('user');

INSERT INTO alumni (email, password, role_id, nim, nama, jurusan, angkatan, tahun_lulus, no_telepon, alamat)
VALUES
('sayu@gmail.com', '$2a$12$OsfNwKXSbGNaLm25cHcWW.aHssa9JKYmrnlQG6e4CvNeFqFcEKJIa', 1, '20160001', 'Sayu Amelia', 'Informatika', 2016, 2020, '081200000001', 'Jl. Melati No. 1'),
('rina.pratama@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20170002', 'Rina Pratama', 'Sistem Informasi', 2017, 2021, '081200000002', 'Jl. Mawar No. 2'),
('budi.santoso@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20150003', 'Budi Santoso', 'Informatika', 2015, 2019, '081200000003', 'Jl. Kenanga No. 3'),
('siti.aisyah@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20140004', 'Siti Aisyah', 'Teknik Industri', 2014, 2018, '081200000004', 'Jl. Anggrek No. 4'),
('andi.wijaya@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20130005', 'Andi Wijaya', 'Informatika', 2013, 2017, '081200000005', 'Jl. Dahlia No. 5'),
('dewi.lestari@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20120006', 'Dewi Lestari', 'Sistem Informasi', 2012, 2016, '081200000006', 'Jl. Flamboyan No. 6'),
('fajar.nugraha@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20110007', 'Fajar Nugraha', 'Teknik Elektro', 2011, 2015, '081200000007', 'Jl. Teratai No. 7'),
('intan.safitri@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20100008', 'Intan Safitri', 'Teknik Mesin', 2010, 2014, '081200000008', 'Jl. Sakura No. 8'),
('yoga.prabowo@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20180009', 'Yoga Prabowo', 'Informatika', 2018, 2022, '081200000009', 'Jl. Bougenville No. 9'),
('nabila.putri@gmail.com', '$2a$12$8f7qEV2pqI4rFU9jHGj37.QOOMWx/KMbb0KRR1lzCzjrD/tas4TXe', 2, '20190010', 'Nabila Putri', 'Sistem Informasi', 2019, 2023, '081200000010', 'Jl. Cemara No. 10');

INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan)
VALUES
(2, 'Perusahaan A', 'Software Engineer', 'Teknologi', 'Jakarta', '8-12jt', '2022-01-10', NULL, 'aktif', 'Pengembangan aplikasi web'),
(3, 'Perusahaan B', 'Data Analyst', 'Konsultan', 'Bandung', '7-10jt', '2021-06-01', '2023-06-01', 'selesai', 'Analisis data bisnis'),
(4, 'Perusahaan C', 'Network Engineer', 'Telekomunikasi', 'Surabaya', '6-9jt', '2020-03-15', NULL, 'aktif', 'Administrasi jaringan'),
(5, 'Perusahaan D', 'QA Engineer', 'Teknologi', 'Yogyakarta', '5-8jt', '2023-02-01', NULL, 'aktif', 'Pengujian perangkat lunak'),
(6, 'Perusahaan E', 'Product Manager', 'Teknologi', 'Jakarta', '15-20jt', '2019-08-20', '2021-12-31', 'selesai', 'Manajemen produk');


