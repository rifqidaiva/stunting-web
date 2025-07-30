package object

import "fmt"

type Report struct {
	Pengguna           Pengguna             `json:"pengguna"`
	Balita             Balita               `json:"balita"`
	Keluarga           Keluarga             `json:"keluarga"`
	LaporanMasyarakat  LaporanMasyarakat    `json:"laporan_masyarakat"`
	RiwayatPemeriksaan []RiwayatPemeriksaan `json:"riwayat_pemeriksaan"`
	Intervensi         []Intervensi         `json:"intervensi"`
}

// MARK: Pengguna
type Pengguna struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Nama         string `json:"nama"`
	Password     string `json:"password"` // Unhashed password for login, not stored in DB
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // "masyarakat" or "admin"
	Alamat       string `json:"alamat"`
}

func (p *Pengguna) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if p.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "Email":
			if p.Email == "" {
				return fmt.Errorf("email is required")
			}
		case "Nama":
			if p.Nama == "" {
				return fmt.Errorf("name is required")
			}
		case "Password":
			if p.Password == "" {
				return fmt.Errorf("password is required")
			}
		case "PasswordHash":
			if p.PasswordHash == "" {
				return fmt.Errorf("password hash is required")
			}
		case "Role":
			if p.Role != "Masyarakat" && p.Role != "Admin" {
				return fmt.Errorf("role must be either 'Masyarakat' or 'Admin'")
			}
		case "Alamat":
			if p.Alamat == "" {
				return fmt.Errorf("address is required")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}

// MARK: LaporanMasyarakat
type LaporanMasyarakat struct {
	Id                    string `json:"id"`
	IdPengguna            string `json:"id_pengguna"`
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

func (l *LaporanMasyarakat) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if l.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "IdPengguna":
			if l.IdPengguna == "" {
				return fmt.Errorf("ID Pengguna is required")
			}
		case "IdBalita":
			if l.IdBalita == "" {
				return fmt.Errorf("ID Balita is required")
			}
		case "IdStatusLaporan":
			if l.IdStatusLaporan == "" {
				return fmt.Errorf("ID Status Laporan is required")
			}
		case "TanggalLaporan":
			if l.TanggalLaporan == "" {
				return fmt.Errorf("tanggal Laporan is required")
			}
		case "HubunganDenganBalita":
			if l.HubunganDenganBalita == "" {
				return fmt.Errorf("hubungan Dengan Balita is required")
			}
		case "NomorHpPelapor":
			if l.NomorHpPelapor == "" {
				return fmt.Errorf("nomor HP Pelapor is required")
			}
		case "NomorHpKeluargaBalita":
			if l.NomorHpKeluargaBalita == "" {
				return fmt.Errorf("nomor HP Keluarga Balita is required")
			}
		}
	}
	return nil
}

// MARK: Balita
type Balita struct {
	Id           string `json:"id"`
	IdKeluarga   string `json:"id_keluarga"`
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	BeratLahir   string `json:"berat_lahir"`
	TinggiLahir  string `json:"tinggi_lahir"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

func (b *Balita) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if b.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "IdKeluarga":
			if b.IdKeluarga == "" {
				return fmt.Errorf("ID Keluarga is required")
			}
		case "Nama":
			if b.Nama == "" {
				return fmt.Errorf("name is required")
			}
		case "TanggalLahir":
			if b.TanggalLahir == "" {
				return fmt.Errorf("tanggal Lahir is required")
			}
		case "JenisKelamin":
			if b.JenisKelamin != "Laki-laki" && b.JenisKelamin != "Perempuan" {
				return fmt.Errorf("jenis Kelamin must be either 'Laki-laki' or 'Perempuan'")
			}
		case "BeratLahir":
			if b.BeratLahir == "" {
				return fmt.Errorf("berat Lahir is required")
			}
		case "TinggiLahir":
			if b.TinggiLahir == "" {
				return fmt.Errorf("tinggi Lahir is required")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}

// MARK: Keluarga
type Keluarga struct {
	Id          string     `json:"id"`
	NomorKk     string     `json:"nomor_kk"`
	NamaAyah    string     `json:"nama_ayah"`
	NamaIbu     string     `json:"nama_ibu"`
	NikAyah     string     `json:"nik_ayah"`
	NikIbu      string     `json:"nik_ibu"`
	Alamat      string     `json:"alamat"`
	Rt          string     `json:"rt"`
	Rw          string     `json:"rw"`
	IdKelurahan string     `json:"id_kelurahan"`
	Koordinat   [2]float64 `json:"koordinat"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

func (k *Keluarga) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if k.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "NomorKk":
			if k.NomorKk == "" {
				return fmt.Errorf("nomor KK is required")
			}
		case "NamaAyah":
			if k.NamaAyah == "" {
				return fmt.Errorf("nama Ayah is required")
			}
		case "NamaIbu":
			if k.NamaIbu == "" {
				return fmt.Errorf("nama Ibu is required")
			}
		case "NikAyah":
			if k.NikAyah == "" {
				return fmt.Errorf("NIK Ayah is required")
			}
		case "NikIbu":
			if k.NikIbu == "" {
				return fmt.Errorf("NIK Ibu is required")
			}
		case "Alamat":
			if k.Alamat == "" {
				return fmt.Errorf("alamat is required")
			}
		case "Rt":
			if k.Rt == "" {
				return fmt.Errorf("RT is required")
			}
		case "Rw":
			if k.Rw == "" {
				return fmt.Errorf("RW is required")
			}
		case "IdKelurahan":
			if k.IdKelurahan == "" {
				return fmt.Errorf("ID Kelurahan is required")
			}
		case "Koordinat":
			if len(k.Koordinat) != 2 || (k.Koordinat[0] == 0 && k.Koordinat[1] == 0) {
				return fmt.Errorf("koordinat must be a valid [longitude, latitude] pair and cannot be [0, 0]")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}

// MARK: RiwayatPemeriksaan
type RiwayatPemeriksaan struct {
	Id          string `json:"id"`
	IdBalita    string `json:"id_balita"`
	Tanggal     string `json:"tanggal"`
	BeratBadan  string `json:"berat_badan"`
	TinggiBadan string `json:"tinggi_badan"`
	StatusGizi  string `json:"status_gizi"`
	Keterangan  string `json:"keterangan"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

func (r *RiwayatPemeriksaan) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if r.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "IdBalita":
			if r.IdBalita == "" {
				return fmt.Errorf("ID Balita is required")
			}
		case "Tanggal":
			if r.Tanggal == "" {
				return fmt.Errorf("tanggal is required")
			}
		case "BeratBadan":
			if r.BeratBadan == "" {
				return fmt.Errorf("berat Badan is required")
			}
		case "TinggiBadan":
			if r.TinggiBadan == "" {
				return fmt.Errorf("tinggi Badan is required")
			}
		case "StatusGizi":
			if r.StatusGizi == "" {
				return fmt.Errorf("status Gizi is required")
			}
		case "Keterangan":
			if r.Keterangan == "" {
				return fmt.Errorf("keterangan is required")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}

// MARK: Intervensi
type Intervensi struct {
	Id        string `json:"id"`
	IdBalita  string `json:"id_balita"`
	Jenis     string `json:"jenis"`
	Tanggal   string `json:"tanggal"`
	Deskripsi string `json:"deskripsi"`
	Hasil     string `json:"hasil"`

	CreatedId   string `json:"created_id"`
	CreatedDate string `json:"created_date"`
	UpdatedId   string `json:"updated_id"`
	UpdatedDate string `json:"updated_date"`
	DeletedId   string `json:"deleted_id"`
	DeletedDate string `json:"deleted_date"`
}

func (i *Intervensi) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if i.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "IdBalita":
			if i.IdBalita == "" {
				return fmt.Errorf("ID Balita is required")
			}
		case "Jenis":
			if i.Jenis == "" {
				return fmt.Errorf("jenis is required")
			}
		case "Tanggal":
			if i.Tanggal == "" {
				return fmt.Errorf("tanggal is required")
			}
		case "Deskripsi":
			if i.Deskripsi == "" {
				return fmt.Errorf("deskripsi is required")
			}
		case "Hasil":
			if i.Hasil == "" {
				return fmt.Errorf("hasil is required")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}
