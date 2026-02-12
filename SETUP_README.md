# Queue System Backend - Setup Guide

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) ติดตั้งและเปิดใช้งานแล้ว

## Quick Start (คำสั่งเดียวรันได้เลย)

```bash
docker compose up --build
```

เพียงเท่านี้! ระบบจะ:
1. สร้าง PostgreSQL container + สร้าง database `queue_system` ให้อัตโนมัติ
2. Build Go backend จาก Dockerfile
3. รอ PostgreSQL พร้อม (healthcheck) แล้วค่อย start backend
4. AutoMigrate สร้าง tables (`queue_states`, `queue_tickets`) ให้อัตโนมัติ

## Ports

| Service    | URL                        |
|------------|----------------------------|
| Backend API | http://localhost:8080      |
| PostgreSQL  | localhost:5432             |

## API ทดสอบ

### Health Check
```bash
curl http://localhost:8080/api/health
```

### รับบัตรคิว
```bash
curl -X POST http://localhost:8080/api/queue/next
```
Response: `{ "queue_number": "A0", "issued_at": "12/02/2026 14:30" }`

### ดูคิวปัจจุบัน
```bash
curl http://localhost:8080/api/queue/current
```
Response: `{ "queue_number": "A0", "issued_at": "12/02/2026 14:30" }`

### ล้างคิว
```bash
curl -X POST http://localhost:8080/api/queue/clear
```
Response: `{ "queue_number": "00", "message": "Queue has been cleared" }`

## คำสั่ง Docker ที่ใช้บ่อย

```bash
# รัน foreground (เห็น log)
docker compose up --build

# รัน background
docker compose up --build -d

# ดู log
docker compose logs -f

# ดู log เฉพาะ backend
docker compose logs -f backend

# หยุด
docker compose down

# หยุด + ลบ database data
docker compose down -v

# rebuild เฉพาะ backend (ไม่ลบ DB data)
docker compose up --build backend
```

## รันแบบไม่ใช้ Docker (Development)

ต้องมี PostgreSQL รันอยู่แล้วบน localhost:5432

```bash
# 1. แก้ .env ให้ตรงกับ PostgreSQL ของคุณ
#    DB_HOST=localhost
#    DB_PORT=5432
#    DB_USER=postgres
#    DB_PASSWORD=postgres
#    DB_NAME=queue_system

# 2. สร้าง database (ถ้ายังไม่มี)
createdb -U postgres queue_system

# 3. รัน
cd backend
go run main.go
```

## Database Schema

### queue_states (1 row เสมอ)
| Column         | Type       | Description              |
|---------------|------------|--------------------------|
| id            | uint (PK)  | Always 1                 |
| current_letter | varchar(1) | A-Z หรือ "" (cleared)    |
| current_number | int        | 0-9 หรือ -1 (cleared)    |
| updated_at    | timestamp  | อัปเดตอัตโนมัติ            |

### queue_tickets (log ทุกบัตรคิว)
| Column       | Type       | Description              |
|-------------|------------|--------------------------|
| id          | uint (PK)  | Auto increment           |
| queue_number | varchar(2) | เช่น "A5"               |
| issued_at   | timestamp  | เวลาที่ออกบัตร             |
