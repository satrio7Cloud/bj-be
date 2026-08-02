# Biro Jasa API

Proyek ini adalah backend REST API untuk sistem Biro Jasa (Perpanjangan STNK, Pengurusan Surat Kendaraan, dll), dibangun menggunakan **Golang** dengan arsitektur **Domain-Driven Design (DDD)**.

## 🚀 Teknologi Utama
- **Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: SQL Server
- **ORM**: [GORM](https://gorm.io/)
- **Caching**: Redis
- **Storage**: MinIO (S3-compatible Object Storage)
- **Autentikasi**: JWT (JSON Web Token)
- **Hashing**: Argon2id
- **Dokumentasi**: Swagger / OpenAPI

---

## 📁 Struktur Proyek (DDD)
Proyek ini mengadopsi struktur *Domain-Driven Design* untuk menjaga kode tetap rapi, terukur (*scalable*), dan mudah di- *maintenance*.

```text
📦 BE
 ┣ 📂 cmd
 ┃ ┗ 📂 server
 ┃   ┗ 📜 main.go           # Entry point aplikasi (Wiring Dependency)
 ┣ 📂 docs                  # Hasil generate Swagger API Docs
 ┣ 📂 internal
 ┃ ┣ 📂 config              # Memuat variabel environment (.env)
 ┃ ┣ 📂 constants           # Berisi konstanta statis
 ┃ ┣ 📂 database            # Koneksi dan inisialisasi SQL Server
 ┃ ┣ 📂 di                  # Dependency Injection (Container)
 ┃ ┣ 📂 domain              # Folder utama untuk setiap domain bisnis (Admin, Order, Service, dll)
 ┃ ┃ ┗ 📂 user              # Contoh Domain User/Admin
 ┃ ┃   ┣ 📂 dto             # Data Transfer Object (Request & Response body)
 ┃ ┃   ┣ 📂 entity          # Struktur data tabel di Database
 ┃ ┃   ┣ 📂 handler         # Fiber HTTP Handlers (Menerima request)
 ┃ ┃   ┣ 📂 repository      # Logika query database GORM
 ┃ ┃   ┣ 📂 routes          # Pendaftaran endpoint API
 ┃ ┃   ┗ 📂 usecase         # Logika bisnis utama
 ┃ ┣ 📂 pkg                 # Driver eksternal (Redis, MinIO)
 ┃ ┣ 📂 routes              # Registrasi utama semua route
 ┃ ┣ 📂 shared              # Kode yang dipakai bersama (misal Storage Service)
 ┃ ┗ 📂 utils               # Fungsi bantuan (Hashing Password, Error Handling, dll)
 ┣ 📂 middleware            # Middleware (Auth, CORS, Logging)
 ┣ 📜 .env                  # Konfigurasi rahasia lokal
 ┣ 📜 docker-compose.yml    # Setup Database, MinIO, Redis via Docker
 ┣ 📜 go.mod                # Dependency Go
 ┗ 📜 README.md             # Dokumentasi ini
```

---

## 🛠️ Persiapan Lingkungan (Setup)

### 1. File Konfigurasi `.env`
Pastikan file `.env` dibuat dengan mengacu pada informasi *credentials* yang diperlukan. Anda perlu memasukkan konfigurasi untuk SQL Server, Redis, dan MinIO:

```env
APP_PORT=
APP_ENV=

# SQL Server
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

# JWT Secret
JWT_SECRET=supersecret

# Redis
REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=

# MinIO
MINIO_ENDPOINT=
MINIO_ACCESS_KEY=
MINIO_SECRET_KEY=
MINIO_USE_SSL=
```

### 2. Menjalankan Infrastruktur Dasar (Docker)
Sistem sangat bergantung pada **SQL Server**, **Redis**, dan **MinIO**. Kami merekomendasikan menggunakan `docker-compose` agar proses *setup* menjadi mudah. Pastikan Docker sudah ter-install, lalu jalankan:

```bash
docker-compose up -d bj_sql_server bj_redis bj_minio
```

- **SQL Server** berjalan di port `1434`
- **Redis** berjalan di port `6379`
- **MinIO** berjalan di port `9000` (Console Web di port `9001`)

---
p
## 🏃 Menjalankan Aplikasi

Jalankan perintah berikut untuk mengunduh semua library dan memulai server di lokal Anda:
```bash
go mod tidy
go run ./cmd/server
```

*(Opsional)* Anda juga bisa menjalankan backend beserta seluruh *services* di dalam Docker:
```bash
docker-compose up --build -d
```

Jika berhasil, terminal akan menampilkan:
> `Berhasil terhubung ke database bj_db`  
> `Server di jalankan di port 3000...`

---

## 📖 Dokumentasi API (Swagger)

Aplikasi ini menggunakan Swagger untuk dokumentasi interaktif yang sudah sepenuhnya di-*generate* secara otomatis.
Setelah aplikasi berjalan, buka browser dan akses URL berikut:

👉 **http://localhost:3000/docs**

*Setiap ada penambahan endpoint baru atau perubahan, jangan lupa jalankan perintah ini untuk memperbarui folder `docs/`:*
```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go
```
