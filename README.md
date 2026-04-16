# 🚀 Account Service - Clean Architecture

[![Go Version](https://img.shields.io/github/go-mod/go-version/codebayu/account-service?style=flat-square)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/codebayu/account-service?style=flat-square)](https://goreportcard.com/report/github.com/codebayu/account-service)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

Account Service adalah sebuah microservice tangguh yang menangani autentikasi dan manajemen pengguna, dibangun menggunakan **Golang** dengan standar **Clean Architecture**.

🎨 Proyek ini dibangun dengan semangat **Vibe Coding** menggunakan **Antigravity** (Powerful AI Coding Assistant).

---

## ✨ Fitur Utama

- **Authentication**: Registrasi dan Login menggunakan JWT (Access & Refresh Token).
- **Profile Management**: Mendapatkan informasi user yang sedang login.
- **Security Middleware**: 
  - **Signature Validation**: Verifikasi integritas setiap request menggunakan HMAC-SHA256.
  - **Header Validation**: Pengecekan wajib untuk `x-signature`, `x-datetime`, dan `x-channel`.
- **Database**: PostgreSQL dengan abstraksi GORM.
- **Clean Architecture**: Pemisahan layer yang jelas antara Handler, Service, Repository, dan DTO.
- **Robust Testing**: Cakupan unit test mencapai **>87%** pada package inti.

---

## 🛠️ Tech Stack

- **Core**: [Go (Golang)](https://golang.org/)
- **Framework**: [Echo v5](https://echo.labstack.com/v5)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Testing**: [Testify](https://github.com/stretchr/testify), [Go-Sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- **Security**: JWT, HMAC-SHA256

---

## 🚀 Cara Menjalankan

### 1. Prasyarat
Pastikan Anda sudah menginstal:
- Go 1.22+
- PostgreSQL
- Make

### 2. Konfigurasi
Salin `.env` dan sesuaikan nilainya:
```bash
cp .env.example .env # Jika ada, atau buat baru sesuai spesifikasi
```
Pastikan `API_KEY`, `API_SECRET`, dan `CHANNEL_ID` sudah diset untuk keamanan signature.

### 3. Jalankan Aplikasi
Gunakan `Makefile` untuk mempermudah eksekusi:

```bash
# Menjalankan database migration
make migrate

# Menjalankan aplikasi
make run
```

---

## 🧪 Testing & Coverage

Kami menjaga kualitas kode dengan pengujian yang ketat.

```bash
# Menjalankan semua unit test
make test

# Mendapatkan laporan coverage (Layar: Handler, Repo, Service, Utils)
make test-coverage

# Melihat laporan coverage visual (HTML)
make coverage-html
```

---

## 🔒 Keamanan: Digital Signature

Setiap request wajib menyertakan:
- `x-signature`: `HMAC-SHA256(apiKey + unixTimestamp, apiSecret)`
- `x-datetime`: Unix Timestamp (detik)
- `x-channel`: Identitas Channel (misal: "WEB")

Contoh logic generasinya dapat dilihat di [index.md](./index.md).

---

## 🤝 Kontribusi

Dibuat dengan ❤️ melalui **Vibe Coding** bersama **Antigravity**.

---

*Copyright © 2026 codebayu. All rights reserved.*
