# Website Stunting Kota Cirebon

## Instalasi dan Setup

### 1. Clone branch `geotagging` pada Repository

```bash
git clone -b geotagging https://github.com/rifqidaiva/stunting-web.git
cd stunting-web
```

### 2. Setup Database

1. Buat database MySQL baru dengan nama stuntingdb.
2. Import stuntingdb.sql.

### 3. Setup Aplikasi

1. Install dependencies Go:

```bash
go mod download
```

2. Jalankan aplikasi backend:

```bash
go run cmd/main.go
```

- Aplikasi berjalan di `http://localhost:8080`
