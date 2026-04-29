# GXApp Daily Report Service

Bot Slack internal pengganti Geekbot untuk mengisi daily standup langsung lewat DM Slack.

## Fitur

- Bot mengirim reminder DM ke semua user aktif setiap hari kerja (Senin–Jumat)
- Pengisian report step-by-step via DM (5 pertanyaan)
- Hasil report otomatis dikirim ke channel Slack yang dikonfigurasi
- Sesi dilanjutkan dari step terakhir jika terganggu
- REST API (health check + siap untuk endpoint dashboard)

---

## Arsitektur

```
Slack DM (user) → Socket Mode → app/slack/Handler → report/service → DB (PostgreSQL)
                                                                    ↓
                                                          Slack Channel (rangkuman)
Scheduler (cron) → SendReminders → Slack DM (reminder)
```

---

## Prasyarat

- Go 1.22+
- PostgreSQL
- Slack App dengan Socket Mode aktif

### Konfigurasi Slack App

1. Buka [api.slack.com/apps](https://api.slack.com/apps) → buat app baru
2. **Socket Mode**: aktifkan di menu *Socket Mode*, salin **App-Level Token** (`SLACK_APP_TOKEN`)
3. **OAuth & Permissions** → tambahkan Bot Token Scopes:
   - `chat:write` — kirim pesan
   - `im:history` — baca DM
   - `im:read` — list DM
   - `im:write` — buka DM channel
4. Install app ke workspace → salin **Bot User OAuth Token** (`SLACK_BOT_TOKEN`)
5. **Event Subscriptions** → Subscribe to Bot Events: `message.im`
6. Catat ID channel tujuan laporan (`SLACK_CHANNEL_ID`, format `C...`)

---

## Setup

### 1. Konfigurasi `.env`

Salin `.env.example` ke `.env` dan isi nilai-nilainya:

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=secret
DB_DATABASE=daily_report
DB_OWNER=postgres

SLACK_BOT_TOKEN=xoxb-...
SLACK_APP_TOKEN=xapp-...
SLACK_CHANNEL_ID=C0123456789
REMINDER_TIME=09:00

HTTP_PORT=3000
TZ=Asia/Makassar
```

### 2. Jalankan Migration

```shell
go build -o application cmd/main.go

./application xtreme:migration
```

### 3. Jalankan Slack Bot

```shell
./application xtreme:slack --dev
```

Bot akan:
- Mendengarkan DM via Socket Mode
- Menjalankan cron scheduler untuk reminder harian
- Membuka HTTP server di `HTTP_PORT` (default `3000`)

---

## Cara Menambahkan User

User harus didaftarkan terlebih dahulu ke tabel `report_users`. Jalankan perintah SQL berikut di database:

```sql
INSERT INTO report_users ("slackId", "name", "email", "isActive")
VALUES ('U0123456789', 'Budi Santoso', 'budi@example.com', true);
```

Cari Slack ID user:
- Buka profil user di Slack → klik ⋮ (tiga titik) → *Copy member ID*
- Format: `U` diikuti 10 karakter, contoh: `U0123456789`

---

## Perintah Runner

```shell
# Jalankan HTTP API (existing gorilla/mux routes)
./application --dev

# Jalankan Slack Bot + scheduler reminder
./application xtreme:slack --dev

# Jalankan migration database
./application xtreme:migration

# Jalankan schedule (gocron-based)
./application xtreme:schedule

# Jalankan seeder
./application xtreme:seeder
```

---

## Alur Kerja Bot

```
[Scheduler 09:00]
  → Kirim DM reminder ke semua user aktif

[User membalas DM apa saja]
  → Bot: "What did you complete yesterday?"
  → User menjawab
  → Bot: "What will you do today?"
  → ... (5 pertanyaan)
  → Bot: "✅ Daily report kamu sudah tersimpan."
  → Channel Slack menerima rangkuman laporan
```

### Edge Cases

| Kondisi | Respons Bot |
|---------|-------------|
| User belum terdaftar | "Kamu belum terdaftar..." |
| User sudah isi hari ini | "Kamu sudah mengisi daily report hari ini ✅" |
| Sesi terputus di tengah | Bot melanjutkan dari step terakhir |

---

## Stack

- **Language**: Go 1.22+
- **HTTP**: Fiber v2 (future REST API)
- **Slack**: slack-go/slack (Socket Mode)
- **Scheduler**: robfig/cron/v3
- **Database**: PostgreSQL via GORM
- **Framework**: globalxtreme/go-core v2

---

## Struktur Folder (Daily Report)

```
internal/
├── app/
│   ├── api/fiber.go              # Fiber HTTP router (health + future endpoints)
│   └── slack/Handler.go          # Socket Mode event handler
├── report/
│   ├── repository/ReportRepository.go
│   └── service/ReportService.go
├── user/
│   ├── repository/UserRepository.go
│   └── service/UserService.go
└── pkg/
    ├── config/slack.go           # Slack client init
    ├── constant/step.go          # Step numbers & question texts
    ├── error/error.go            # App error variables
    ├── model/{ReportUser,Report,ReportSession}.go
    ├── parser/ReportParser.go    # Format channel message
    └── thirdparty/slack/SlackClient.go
```

---

## Boilerplate (GlobalXtreme)

## ✨Getting Started

### Prerequisites
- **Go Version**: [Go](https://go.dev/) version [1.25](https://go.dev/doc/devel/release#go1.25.0) or above

### Installation

1. Download ["go-create-project.sh"](https://storage.globalxtreme-gateway.net/link/installations/go-create-project.sh) and place it in the project directory.
2. Navigate & allow script to be executed
```shell
cd /path/your-go-project-directory

chmod +x go-create-project.sh
```
3. Create new project.
```shell
./go-create-project.sh <project-name>
```


### Installation with Executable
If you want to install go-create-project globally and call it from any directory, use the following steps:

1. Install shc
```shell
brew install shc
```

2. Convert the shell script "go-create-project.sh" to a binary file
```shell
shc -f go-create-project.sh
```

3. Rename and move binary file into system's PATH
```shell
mv go-create-project.sh.x /usr/local/bin/go-create-project
```

4. Create new project anywhere using global command
```shell
go-create-project <project-and-path-name>
```

### Set Up
Set the environment variables before building project.
Update these variables in the `.env` file
```javascript
VERSION=<version>
SERVICE=<service-name>

DB_HOST=<host>
DB_PORT=<port>
DB_USERNAME=<uname>
DB_PASSWORD=<pass>
DB_DATABASE=<dbname>
```
Set the environment variables for RabbitMQ Connection
```javascript
DB_RABBITMQ_HOST=<host>
DB_RABBITMQ_PORT=<port>
DB_RABBITMQ_USERNAME=<uname>
DB_RABBITMQ_PASSWORD=<pass>
DB_RABBITMQ_DATABASE=<dbname>

RABBITMQ_GLOBAL_HOST=<rabbitmq-host>
RABBITMQ_GLOBAL_PORT=<rabbitmq-port>
RABBITMQ_GLOBAL_USER=<rabbitmq-uname>
RABBITMQ_GLOBAL_PASSWORD=<rabbitmq-pass>
```

### Running the Application
Run the application
```shell
go run cmd/main.go --dev

# or

go build -o application cmd/main.go

./application --dev
```
Runner commands, to run API, gRPC Server, and others.
```shell
# Migration
./application xtreme:migration

# Seeder
./application xtreme:seeder

# gRPC
./application xtreme:grpc

# Queue
./application xtreme:queue

# RabbitMQ
./application xtreme:rabbitmq

# Schedule
./application xtreme:schedule

# Custom Commands (Example)
./application dev-test
```
Generator commands to generate migration file, model, and others.
```shell
# Migration file
./application gen:migration <Name>

# Handler File
./application gen:handler <Name> --type=<web/mobile> --resource

# Model File
./application gen:model <Name> --migration

# Parser File
./application gen:parser <Name> --model
```
⚠️ Add `--dev` flag for development mode.

## �� Documentation
For documentation and Production deployment guide, you can read it on [Go-Lang Backend Service](https://www.notion.so/globalxtreme/Go-Lang-Backend-Service-527f335297b8465f838fc2598538dae7?pvs=4), which we will create later!
