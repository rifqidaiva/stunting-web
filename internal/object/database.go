package object

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"root",
		"",
		"127.0.0.1",
		"3306",
		"stuntingdb_new",
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// MARK: Pengguna
type Pengguna struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // "masyarakat", "admin", or "petugas kesehatan"
}

// MARK: Masyarakat
type Masyarakat struct {
	Id         string `json:"id"`
	IdPengguna string `json:"id_pengguna"`
	Nama       string `json:"nama"`
	Alamat     string `json:"alamat"`
}

// MARK: PetugasKesehatan
type PetugasKesehatan struct {
	Id         string `json:"id"`
	IdPengguna string `json:"id_pengguna"`
	IdSkpd     string `json:"id_skpd"`
	Nama       string `json:"nama"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: Skpd
type Skpd struct {
	Id    string `json:"id"`
	Skpd  string `json:"skpd"`
	Jenis string `json:"jenis"` // "puskesmas", "kelurahan", "skpd"

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: Balita
type Balita struct {
	Id           string `json:"id"`
	IdKeluarga   string `json:"id_keluarga"`
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"` // "L" or "P"
	BeratLahir   string `json:"berat_lahir"`
	TinggiLahir  string `json:"tinggi_lahir"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: Keluarga
type Keluarga struct {
	Id          string `json:"id"`
	NomorKk     string `json:"nomor_kk"`
	NamaAyah    string `json:"nama_ayah"`
	NamaIbu     string `json:"nama_ibu"`
	NikAyah     string `json:"nik_ayah"`
	NikIbu      string `json:"nik_ibu"`
	Alamat      string `json:"alamat"`
	Rt          string `json:"rt"`
	Rw          string `json:"rw"`
	IdKelurahan string `json:"id_kelurahan"`
	Koordinat   string `json:"koordinat"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: LaporanMasyarakat
type LaporanMasyarakat struct {
	Id                    string `json:"id"`
	IdMasyarakat          string `json:"id_masyarakat"` // if NULL, this is an admin report
	IdBalita              string `json:"id_balita"`
	IdStatusLaporan       string `json:"id_status_laporan"`
	TanggalLaporan        string `json:"tanggal_laporan"`
	HubunganDenganBalita  string `json:"hubungan_dengan_balita"`
	NomorHpPelapor        string `json:"nomor_hp_pelapor"`
	NomorHpKeluargaBalita string `json:"nomor_hp_keluarga_balita"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: RiwayatPemeriksaan
type RiwayatPemeriksaan struct {
	Id           string `json:"id"`
	IdBalita     string `json:"id_balita"`
	IdIntervensi string `json:"id_intervensi"`
	Tanggal      string `json:"tanggal"`
	BeratBadan   string `json:"berat_badan"`
	TinggiBadan  string `json:"tinggi_badan"`
	StatusGizi   string `json:"status_gizi"` // "normal", "stunting", "gizi buruk"
	Keterangan   string `json:"keterangan"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: Intervensi
type Intervensi struct {
	Id        string `json:"id"`
	Jenis     string `json:"jenis"` // "gizi", "kesehatan", "sosial"
	Tanggal   string `json:"tanggal"`
	Deskripsi string `json:"deskripsi"`
	Hasil     string `json:"hasil"`
	// IdPetugasKesehatan string `json:"id_petugas_kesehatan"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

// MARK: IntervensiPetugas (Junction Table)
type IntervensiPetugas struct {
	Id                 string `json:"id"`
	IdIntervensi       string `json:"id_intervensi"`
	IdPetugasKesehatan string `json:"id_petugas_kesehatan"`
}

// MARK: StatusLaporan
type StatusLaporan struct {
	Id     string `json:"id"`
	Status string `json:"status"` // "Belum diproses", "Diproses dan data tidak sesuai", "Diproses dan data sesuai", "Belum ditindaklanjuti", "Sudah ditindaklanjuti", "Sudah perbaikan gizi"
}

// MARK: Kelurahan
type Kelurahan struct {
	Id          string `json:"id"`
	IdKecamatan string `json:"id_kecamatan"`
	Kelurahan   string `json:"kelurahan"`
	Area        string `json:"area"`
}

// MARK: Kecamatan
type Kecamatan struct {
	Id        string `json:"id"`
	Kecamatan string `json:"kecamatan"`
	Area      string `json:"area"`
}
