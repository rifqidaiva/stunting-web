# Website Stunting Kota Cirebon

## Deskripsi

Website sistem informasi stunting untuk Kota Cirebon yang memungkinkan masyarakat melaporkan kasus stunting, petugas kesehatan menangani intervensi, dan admin mengelola data secara terpusat.

## Fitur Utama

- **Admin**: Manajemen lengkap semua data (balita, keluarga, laporan, intervensi, petugas kesehatan, SKPD)
- **Masyarakat**: Pendaftaran keluarga, balita, dan pelaporan kasus stunting
- **Petugas Kesehatan**: Penanganan intervensi yang ditugaskan
- **Dashboard & Mapping**: Visualisasi data dengan peta interaktif menggunakan GeoJSON

## Teknologi

- **Backend**: Go (Golang) dengan MySQL database
- **Frontend**: Vue.js 3 dengan TypeScript dan Vite
- **Database**: MySQL dengan data geospasial (koordinat latitude/longitude)
- **Authentication**: JWT (JSON Web Token)

## Prasyarat

- Go 1.21 atau lebih baru
- Node.js 18 atau lebih baru
- MySQL 8.0 atau lebih baru
- Git

## Instalasi dan Setup

### 1. Clone Repository

```bash
git clone https://github.com/rifqidaiva/stunting-web.git
cd stunting-web
```

### 2. Setup Database

1. Buat database MySQL baru dengan nama stuntingdb_new.
2. Import stuntingdb_new.

### 3. Setup Backend

1. Install dependencies Go:

```bash
go mod download
```

2. Jalankan aplikasi backend:

```bash
go run main.go
```

- Backend akan berjalan di `http://localhost:8080`
- Dokumentasi swagger berada di `http://localhost:8080/swagger`

### 4. Setup Frontend

> [!NOTE]
> Jalankan Frontend dalam terminal yang berbeda dari backend.

1. Masuk ke direktori web:

```bash
cd web
```

2. Install dependencies:

```bash
npm install
```

3. Jalankan development server:

```bash
npm run dev
```

Frontend akan berjalan di `http://localhost:3000`

## Struktur Project

```
stunting-web/
├── main.go                # Entry point aplikasi
├── go.mod                 # Go dependencies
├── stuntingdb_new.sql     # Database schema dan data
├── docs/                  # Swagger documentation
├── internal/
│   ├── api/               # API handlers
│   │   ├── admin/         # Admin endpoints
│   │   ├── auth/          # Authentication endpoints
│   │   ├── community/     # Community endpoints
│   │   └── health_worker/ # Health worker endpoints
│   └── object/            # Data structures dan utilities
└── web/                   # Vue.js frontend
    ├── src/
    │   ├── components/    # Vue components
    │   │   ├── admin/     # Admin interface
    │   │   ├── auth/      # Authentication forms
    │   │   ├── community/ # Community interface
    │   │   └── ui/        # Reusable UI components
    │   └── assets/        # Frontend assets
    └── public/            # Public files
```

## Penggunaan

### Login Sebagai Admin

1. Akses `http://localhost:3000/auth`
2. Gunakan credentials admin yang sudah terdaftar:
   - email: dwiki@gmail.com
   - password: 123456
3. Kelola data balita, keluarga, laporan, dan intervensi

### Login Sebagai Masyarakat

1. Register akun baru sebagai masyarakat
2. Login dengan akun yang telah dibuat
3. Daftarkan keluarga dan balita
4. Buat laporan jika menemukan kasus stunting

### Login Sebagai Petugas Kesehatan

1. Gunakan akun petugas kesehatan yang telah didaftarkan admin
2. Lihat daftar intervensi yang ditugaskan
3. Update status penanganan intervensi

## Development

### Database Schema

Database menggunakan soft delete pattern dengan kolom:

- `created_at`: Timestamp pembuatan
- `updated_at`: Timestamp update terakhir
- `deleted_at`: Timestamp penghapusan (NULL jika aktif)

### Authentication Flow

1. User login → JWT token digenerate
2. Token disimpan di localStorage frontend
3. Setiap request API menyertakan token di header Authorization
4. Backend memvalidasi token dan role user
