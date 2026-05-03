# Replatform Progress: safa (Laravel) → safa-went (Go/Chi)

## Tamamlananlar

### Faz 1 — Auth Altyapısı ✅
- `database/models/user.go` — `password`, `email_verified_at` alanları eklendi
- `http/middlewares/auth.go` — JWT doğrulama implemente edildi; `/auth/register`, `/auth/login`, `/swagger/*`, `/ping` public path olarak bırakıldı
- `http/controllers/auth_controller.go` — `Register`, `Login`, `Me`, `Logout` metodları
- `http/requests/auth_request.go` — `RegisterPayload`, `LoginPayload`
- `http/resources/auth_resource.go` — `AuthResource`, `AuthUserResource`, `NewAuthUserResource`
- `routes/auth_router.go` — `/auth/register`, `/auth/login`, `/auth/user`, `/auth/logout`
- `main.go` — Swagger `@securityDefinitions.apikey BearerAuth` eklendi

### Altyapı İyileştirmeleri ✅
- `http/resources/pagination.go` — Generic `PaginateQuery[M, R]` helper eklendi (Go generics); tüm resource'larda pagination tekrarı kaldırıldı
- `perPage` için 100 üst sınırı `PaginateQuery` içinde merkezi olarak uygulanıyor

### Faz 2 — Domain Modelleri (Kısmen)

#### Requester ✅
- `database/models/requester.go`
- `database/migrations/20260503195455_create_requesters_table.up/down.sql`
- `http/requests/requester_request.go` — `RequesterPayload`, `RequesterUpdatePayload`
- `http/resources/requester_resource.go` — `Paginate`, `Find`, `Search` (pagination destekli)
- `http/controllers/requester_controller.go` — CRUD + search + Swagger
- `routes/requester_router.go`

#### Approver ✅
- `database/models/approver.go`
- `database/migrations/20260503201452_create_approvers_table.up/down.sql`
- `http/requests/approver_request.go`
- `http/resources/approver_resource.go`
- `http/controllers/approver_controller.go` — CRUD + search + Swagger
- `routes/approver_router.go`

---

## Yapılacaklar

### Faz 2 — Kalan Modeller

#### PrintRequest ⬜
- Model: `requested_at` date, `color_copies` int, `bw_copies` int, `description` text nullable, `requester_id` FK, `approver_id` FK, soft delete
- Resource: filtreler → `requester_names[]`, `approver_names[]`, `color_copies_min/max`, `bw_copies_min/max`, `requested_at_from/to`, sort, pagination
- Ekstra endpoint'ler: export (Faz 5'te)

#### Publisher ⬜
- Model: `name` unique, soft delete

#### Author ⬜
- Model: `name` unique, soft delete

#### Book ⬜
- Model: `name`, `barcode` unique nullable, `author_id` FK, `publisher_id` FK, `language`, `page_count`, `is_donation` bool, `shelf_code`, `fixture_no` int unique, `level` enum(`ilkokul`/`ortaokul`/`ortak`), soft delete
- Resource filtreler: `name`, `fixture_no`, `search`, `author_id`, `publisher_id`, `level`, `is_donation`, sort (default: `fixture_no`)
- İş mantığı: `IsCurrentlyLoaned() bool`

#### Classroom ⬜
- Model: `name`, soft delete
- İş mantığı: `GetLevel() string` — isimden sınıf numarasını parse et → `ilkokul` (1-4), `ortaokul` (5-8), `ortak`

#### Student ⬜
- Model: `name`, `classroom_id` FK cascade, soft delete
- Resource filtreler: `search`, `classroom_id`
- İş mantığı: `GetLevel() string` (classroom'a delege)

#### Loan ⬜
- Model: `student_id` FK, `book_id` FK, `loan_date`, `due_date`, `return_date` nullable, `status` enum(`active`/`returned`/`overdue`), `notes` nullable, soft delete
- Resource filtreler: `status`, `search`, `student_id`, `book_id`, sort
- Ekstra endpoint'ler: `POST /loans/{id}/return`, `POST /loans/check-availability`
- İş mantığı: `CanStudentBorrowBook(student, book) (bool, string)` — seviye eşleşme + aktif loan kontrolü

---

### Faz 3 — İş Mantığı ⬜
- `Loan.CanStudentBorrowBook` static metodu
- `Book.IsCurrentlyLoaned`
- `Classroom.GetLevel` / `Student.GetLevel`
- `LoanController.ReturnBook`
- `LoanController.CheckAvailability`

### Faz 4 — Gelişmiş Filtreleme ⬜
- PrintRequest, Book, Student, Loan resource'larına filtre parametreleri

### Faz 5 — Excel Export ⬜
- `go get github.com/xuri/excelize/v2`
- `GET /print-requests/export/all` — opsiyonel tarih aralığı
- `GET /print-requests/export/by-requester` — requester başına toplam
- `GET /print-requests/export/comparison` — iki dönem karşılaştırması

### Faz 6 — Güvenlik ⬜
- CORS: `AllowedOrigins: ["*"]` + `AllowCredentials: true` düzeltilecek (production origin listesi)

---

## Mevcut Route Listesi

| Method | Path | Auth |
|--------|------|------|
| POST | /auth/register | Public |
| POST | /auth/login | Public |
| GET | /auth/user | Bearer |
| POST | /auth/logout | Bearer |
| GET | /user | Bearer |
| GET | /user/{id} | Bearer |
| POST | /user | Bearer |
| PUT | /user/{id} | Bearer |
| DELETE | /user/{id} | Bearer |
| GET | /requesters | Bearer |
| GET | /requesters/{id} | Bearer |
| POST | /requesters | Bearer |
| PUT | /requesters/{id} | Bearer |
| DELETE | /requesters/{id} | Bearer |
| GET | /approvers | Bearer |
| GET | /approvers/{id} | Bearer |
| POST | /approvers | Bearer |
| PUT | /approvers/{id} | Bearer |
| DELETE | /approvers/{id} | Bearer |
