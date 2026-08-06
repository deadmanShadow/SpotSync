# 🛠️ SpotSync Backend — Step-by-Step Implementation Guide for Puku CLI

> **Target Platform:** Go 1.22+ (Echo v4 + GORM + PostgreSQL + JWT + bcrypt + Validator)  
> **Architecture Pattern:** Strict Clean Architecture (DTO → Handler → Service → Repository → Models)  
> **Key Feature:** High-concurrency EV Charging reservation system with Row-Level DB Locking (`FOR UPDATE`).

---

## 🤖 Instructions for Puku CLI Execution

When instructed to implement a step (e.g., `"Implement Step 1 of Backend.MD"`):
1. **Read the target Step completely.** Each step contains full context, file paths, exact schemas, function signatures, logic rules, and validation commands.
2. **Strictly adhere to the Clean Architecture rules.** Handlers MUST NOT access the database directly; Repositories MUST handle database calls; Services MUST contain business logic.
3. **Execute all file creation and modifications** described in that specific step.
4. **Run the provided Verification Commands** to test your implementation before declaring the step complete.

---

## 📁 Repository Directory Blueprint

```
backend/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point (DI wiring & server startup)
├── config/
│   └── config.go                # Environment variable loader
├── database/
│   └── postgres.go              # Database connection & auto-migrations
├── dto/
│   ├── auth_dto.go              # Register, Login DTOs
│   ├── zone_dto.go              # Zone creation & response DTOs
│   └── reservation_dto.go       # Reservation DTOs
├── handler/
│   ├── auth_handler.go          # Auth HTTP handlers
│   ├── zone_handler.go          # Parking Zone HTTP handlers
│   └── reservation_handler.go   # Reservation HTTP handlers
├── middleware/
│   ├── auth_middleware.go       # JWT Bearer verification
│   └── role_middleware.go       # Role-based access control (admin/driver)
├── models/
│   ├── user.go                  # User GORM model
│   ├── parking_zone.go          # Parking Zone GORM model
│   └── reservation.go           # Reservation GORM model
├── repository/
│   ├── user_repository.go       # DB access for Users
│   ├── zone_repository.go       # DB access for Zones
│   └── reservation_repository.go# DB access for Reservations (FOR UPDATE locks)
├── service/
│   ├── auth_service.go          # Business logic for Auth
│   ├── zone_service.go          # Business logic for Zones
│   └── reservation_service.go   # Business logic for Reservations
├── utils/
│   ├── context.go               # Echo context helpers
│   ├── error_handler.go         # Echo central error handler
│   ├── jwt.go                   # JWT signing & validation
│   ├── password.go              # Bcrypt hashing utilities
│   ├── response.go              # Standardized API response formatters
│   └── validator.go             # Custom struct validator for Echo
├── .env.example
├── go.mod
└── go.sum
```

---

## 📌 Step 1: Project Initialization, Dependencies & Base Infrastructure

### 🎯 Objective
Initialize the Go module, install required third-party libraries, set up the project folder structure, create environment loader, database connector, helper response utilities, custom struct validator, and basic Echo server.

### 📋 Prerequisites
- Go 1.22 or higher installed.

### 📁 Files to Create
- `backend/go.mod`
- `backend/.env.example`
- `backend/.env`
- `backend/config/config.go`
- `backend/database/postgres.go`
- `backend/utils/response.go`
- `backend/utils/validator.go`
- `backend/cmd/api/main.go`

### 📝 Step Context & Implementation Details

1. **Go Module Initialization & Dependencies**
   Run the following dependencies:
   - `github.com/labstack/echo/v4` (Web Framework)
   - `gorm.io/gorm` & `gorm.io/driver/postgres` (ORM & PostgreSQL Driver)
   - `github.com/joho/godotenv` (Environment file reader)
   - `github.com/go-playground/validator/v10` (Struct tag validator)
   - `github.com/golang-jwt/jwt/v5` (JWT token handling)
   - `golang.org/x/crypto/bcrypt` (Password hashing)

2. **Environment Configuration (`config/config.go`)**
   Define a `Config` struct holding:
   - `Port` (default `"8080"`)
   - `DBHost`, `DBUser`, `DBPassword`, `DBName`, `DBPort`, `DBSSLMode` (or `DATABASE_URL`)
   - `JWTSecret` (default `"supersecretjwtkey"`)
   Load from `.env` using `godotenv.Load()`.

3. **Database Connection Setup (`database/postgres.go`)**
   - Implement `ConnectDB(cfg *config.Config) (*gorm.DB, error)`
   - Open GORM connection to PostgreSQL.
   - Configure connection pool:
     - `SetMaxOpenConns(25)`
     - `SetMaxIdleConns(10)`
     - `SetConnMaxLifetime(15 * time.Minute)`

4. **Standard API Response Wrappers (`utils/response.go`)**
   Enforce standard API format across the backend:
   - Success Response:
     ```json
     {
       "success": true,
       "message": "Operation successful",
       "data": { ... }
     }
     ```
   - Error Response:
     ```json
     {
       "success": false,
       "message": "Error message",
       "errors": "Detailed error string"
     }
     ```
   Provide helper functions: `JSONSuccess(c echo.Context, code int, message string, data interface{})` and `JSONError(c echo.Context, code int, message string, err string)`.

5. **Custom Echo Validator (`utils/validator.go`)**
   Create a struct `CustomValidator` wrapping `validator.Validate`. Implement `Validate(i interface{}) error` method so Echo can automatically validate request bodies via `c.Validate(req)`.

6. **Entry Point Starter (`cmd/api/main.go`)**
   - Load config.
   - Initialize Echo instance `e := echo.New()`.
   - Bind `e.Validator = utils.NewValidator()`.
   - Add basic health check route: `GET /health` returning `{"status": "UP"}`.
   - Connect to DB (log success/failure).
   - Start Echo server on configured port.

### 💻 Execution Commands
```bash
mkdir -p backend/cmd/api backend/config backend/database backend/dto backend/handler backend/middleware backend/models backend/repository backend/service backend/utils
cd backend
go mod init spotsync-backend
go get github.com/labstack/echo/v4
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

### ✅ Step 1 Verification
Run the server:
```bash
go run cmd/api/main.go
```
Test health check endpoint:
```bash
curl http://localhost:8080/health
```
Expected Output: `200 OK` with JSON `{"status": "UP"}`.

---

## 📌 Step 2: GORM Database Models & Auto-Migration System

### 🎯 Objective
Define GORM models corresponding to PostgreSQL tables (`users`, `parking_zones`, `reservations`) with proper tags, relationships, and auto-migration logic.

### 📋 Prerequisites
- Step 1 completed.

### 📁 Files to Create / Modify
- `backend/models/user.go`
- `backend/models/parking_zone.go`
- `backend/models/reservation.go`
- `backend/database/postgres.go` (Update with AutoMigrate call)

### 📝 Step Context & Implementation Details

1. **User Model (`models/user.go`)**
   Struct `User`:
   - `ID` (uint, `gorm:"primaryKey;autoIncrement"`)
   - `Name` (string, `gorm:"type:varchar(100);not null"`)
   - `Email` (string, `gorm:"type:varchar(150);uniqueIndex;not null"`)
   - `Password` (string, `gorm:"type:varchar(255);not null"`)
   - `Role` (string, `gorm:"type:varchar(20);not null;default:'driver'"`) — Must be `'driver'` or `'admin'`
   - `CreatedAt` (time.Time)
   - `UpdatedAt` (time.Time)

2. **ParkingZone Model (`models/parking_zone.go`)**
   Struct `ParkingZone`:
   - `ID` (uint, `gorm:"primaryKey;autoIncrement"`)
   - `Name` (string, `gorm:"type:varchar(100);not null"`)
   - `Type` (string, `gorm:"type:varchar(30);not null"`) — Must be `'general'`, `'ev_charging'`, or `'covered'`
   - `TotalCapacity` (int, `gorm:"not null"`) — Integer > 0
   - `PricePerHour` (float64, `gorm:"type:numeric(10,2);not null"`) — Decimal > 0
   - `CreatedAt` (time.Time)
   - `UpdatedAt` (time.Time)
   - `Reservations` ([]Reservation, `gorm:"foreignKey:ZoneID"`)

3. **Reservation Model (`models/reservation.go`)**
   Struct `Reservation`:
   - `ID` (uint, `gorm:"primaryKey;autoIncrement"`)
   - `UserID` (uint, `gorm:"not null;index"`)
   - `User` (User, `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`)
   - `ZoneID` (uint, `gorm:"not null;index"`)
   - `Zone` (ParkingZone, `gorm:"foreignKey:ZoneID;constraint:OnDelete:CASCADE"`)
   - `LicensePlate` (string, `gorm:"type:varchar(15);not null"`)
   - `Status` (string, `gorm:"type:varchar(20);not null;default:'active'"`) — Must be `'active'`, `'completed'`, or `'cancelled'`
   - `CreatedAt` (time.Time)
   - `UpdatedAt` (time.Time)

4. **Auto-Migration Integration (`database/postgres.go`)**
   Add `AutoMigrate(db *gorm.DB) error` function that calls:
   ```go
   db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{})
   ```
   Call `database.AutoMigrate(db)` in `main.go` after DB connection.

### 💻 Execution Commands
```bash
cd backend
go run cmd/api/main.go
```

### ✅ Step 2 Verification
Check database tables created in PostgreSQL (or PostgreSQL container/NeonDB):
- `users` table created with columns `id`, `name`, `email`, `password`, `role`, `created_at`, `updated_at`.
- `parking_zones` table created with columns `id`, `name`, `type`, `total_capacity`, `price_per_hour`, `created_at`, `updated_at`.
- `reservations` table created with foreign keys `user_id` and `zone_id`.

---

## 📌 Step 3: DTOs, Password Hashing, JWT Utility & Auth Module

### 🎯 Objective
Implement User Registration (`POST /api/v1/auth/register`) and User Login (`POST /api/v1/auth/login`) with bcrypt hashing, JWT token generation, struct validation, and clean layer separation (DTO → Handler → Service → Repository).

### 📋 Prerequisites
- Step 2 completed.

### 📁 Files to Create / Modify
- `backend/dto/auth_dto.go`
- `backend/utils/password.go`
- `backend/utils/jwt.go`
- `backend/repository/user_repository.go`
- `backend/service/auth_service.go`
- `backend/handler/auth_handler.go`
- `backend/cmd/api/main.go` (Wire dependencies & register auth routes)

### 📝 Step Context & Implementation Details

1. **Auth DTOs (`dto/auth_dto.go`)**
   - `RegisterRequest`:
     - `Name` (string, `json:"name" validate:"required,min=2,max=100"`)
     - `Email` (string, `json:"email" validate:"required,email"`)
     - `Password` (string, `json:"password" validate:"required,min=6"`)
     - `Role` (string, `json:"role" validate:"omitempty,oneof=driver admin"`)
   - `LoginRequest`:
     - `Email` (string, `json:"email" validate:"required,email"`)
     - `Password` (string, `json:"password" validate:"required"`)
   - `UserResponse`:
     - `ID` (uint), `Name` (string), `Email` (string), `Role` (string), `CreatedAt` (time.Time), `UpdatedAt` (time.Time)
   - `LoginResponse`:
     - `Token` (string `json:"token"`), `User` (UserResponse `json:"user"`)

2. **Password Utility (`utils/password.go`)**
   - `HashPassword(password string) (string, error)` using `bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)` (cost 10).
   - `CheckPassword(hashedPassword, password string) bool` using `bcrypt.CompareHashAndPassword`.

3. **JWT Utility (`utils/jwt.go`)**
   - Struct `JWTClaims`: `UserID` (uint), `Role` (string), RegisteredClaims.
   - `GenerateToken(userID uint, role string, secret string) (string, error)`: Signs token with `HS256`, 24-hour expiration.
   - `ValidateToken(tokenStr string, secret string) (*JWTClaims, error)`.

4. **User Repository Interface & Struct (`repository/user_repository.go`)**
   - Interface `UserRepository`:
     - `Create(user *models.User) error`
     - `FindByEmail(email string) (*models.User, error)`
     - `FindByID(id uint) (*models.User, error)`
   - Implement with GORM DB pointer.

5. **Auth Service (`service/auth_service.go`)**
   - Interface `AuthService`:
     - `Register(req dto.RegisterRequest) (*dto.UserResponse, error)`
     - `Login(req dto.LoginRequest) (*dto.LoginResponse, error)`
   - Logic Rules:
     - Default role to `"driver"` if empty.
     - Check if email already exists (`ErrEmailAlreadyExists` -> 400 Bad Request).
     - Hash password using bcrypt.
     - Save user model.
     - Login verifies email existence, compares bcrypt password, and generates JWT with `user_id` and `role`.

6. **Auth Handler (`handler/auth_handler.go`)**
   - `Register(c echo.Context) error`:
     - Bind DTO, call `c.Validate(req)`. Return 400 if validation fails.
     - Call `authService.Register(req)`.
     - Return `201 Created` with `utils.JSONSuccess`.
   - `Login(c echo.Context) error`:
     - Bind DTO, call `c.Validate(req)`.
     - Call `authService.Login(req)`.
     - Return `200 OK` with `utils.JSONSuccess`.

7. **Dependency Injection & Route Wiring (`cmd/api/main.go`)**
   Instantiate: `userRepo := repository.NewUserRepository(db)` → `authService := service.NewAuthService(userRepo, cfg)` → `authHandler := handler.NewAuthHandler(authService)`.
   Register Routes:
   - `POST /api/v1/auth/register`
   - `POST /api/v1/auth/login`

### 💻 Execution Commands
```bash
cd backend
go run cmd/api/main.go
```

### ✅ Step 3 Verification
1. **Register a Driver:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"name":"John Driver","email":"john@spotsync.com","password":"password123","role":"driver"}'
   ```
   Expected: `201 Created` with user data JSON.

2. **Register an Admin:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"name":"Admin User","email":"admin@spotsync.com","password":"adminpassword","role":"admin"}'
   ```
   Expected: `201 Created`.

3. **Login:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"john@spotsync.com","password":"password123"}'
   ```
   Expected: `200 OK` returning `token` and `user` object.

---

## 📌 Step 4: Echo JWT & Role Middleware (Authentication & Authorization)

### 🎯 Objective
Create secure middleware for Echo framework to protect endpoints. Parse standard `Authorization: Bearer <token>` headers, store `user_id` and `role` in Echo context, and enforce Role-Based Access Control (`driver` / `admin`).

### 📋 Prerequisites
- Step 3 completed.

### 📁 Files to Create / Modify
- `backend/middleware/auth_middleware.go`
- `backend/middleware/role_middleware.go`
- `backend/utils/context.go`

### 📝 Step Context & Implementation Details

1. **Context Helpers (`utils/context.go`)**
   - `GetUserID(c echo.Context) (uint, error)`: Extracts `user_id` set in Echo context. Return error if missing/invalid type.
   - `GetUserRole(c echo.Context) (string, error)`: Extracts `role` set in Echo context.

2. **JWT Auth Middleware (`middleware/auth_middleware.go`)**
   - `JWTMiddleware(jwtSecret string) echo.MiddlewareFunc`:
     - Reads header `Authorization`. If missing or doesn't start with `"Bearer "`, return `401 Unauthorized` using `utils.JSONError`.
     - Extracts raw token string, calls `utils.ValidateToken(tokenStr, jwtSecret)`.
     - If invalid or expired, return `401 Unauthorized`.
     - Set claims into Echo context: `c.Set("user_id", claims.UserID)` and `c.Set("role", claims.Role)`.
     - Call `next(c)`.

3. **Role Authorization Middleware (`middleware/role_middleware.go`)**
   - `RequireRole(allowedRoles ...string) echo.MiddlewareFunc`:
     - Retrieves `role` from Echo context using `utils.GetUserRole(c)`.
     - Checks if user role is included in `allowedRoles`.
     - If not authorized, return `403 Forbidden` (`{"success": false, "message": "Forbidden: Insufficient permissions", "errors": "Role 'driver' is not authorized to access this resource"}`).
     - If authorized, call `next(c)`.

### 💻 Execution Commands
```bash
cd backend
go build ./...
```

### ✅ Step 4 Verification
Unit / Integration test verification will occur when applying this middleware to protected routes in Step 5 & Step 6.

---

## 📌 Step 5: Parking Zones Module (CRUD & Dynamic Capacity Calculation)

### 🎯 Objective
Implement Parking Zone management: Admin creation (`POST /api/v1/zones`), Public catalog with dynamically calculated `available_spots` (`GET /api/v1/zones`), and Single Zone detail (`GET /api/v1/zones/:id`).

### 📋 Prerequisites
- Step 4 completed.

### 📁 Files to Create / Modify
- `backend/dto/zone_dto.go`
- `backend/repository/zone_repository.go`
- `backend/service/zone_service.go`
- `backend/handler/zone_handler.go`
- `backend/cmd/api/main.go` (Wire zone module routes)

### 📝 Step Context & Implementation Details

1. **Zone DTOs (`dto/zone_dto.go`)**
   - `CreateZoneRequest`:
     - `Name` (string `json:"name" validate:"required,min=3,max=100"`)
     - `Type` (string `json:"type" validate:"required,oneof=general ev_charging covered"`)
     - `TotalCapacity` (int `json:"total_capacity" validate:"required,gt=0"`)
     - `PricePerHour` (float64 `json:"price_per_hour" validate:"required,gt=0"`)
   - `ZoneResponse`:
     - `ID` (uint `json:"id"`)
     - `Name` (string `json:"name"`)
     - `Type` (string `json:"type"`)
     - `TotalCapacity` (int `json:"total_capacity"`)
     - `AvailableSpots` (int `json:"available_spots"`)
     - `PricePerHour` (float64 `json:"price_per_hour"`)
     - `CreatedAt` (time.Time `json:"created_at"`)
     - `UpdatedAt` (time.Time `json:"updated_at,omitempty"`)

2. **Zone Repository (`repository/zone_repository.go`)**
   - Interface `ZoneRepository`:
     - `Create(zone *models.ParkingZone) error`
     - `FindAllWithAvailability() ([]dto.ZoneResponse, error)`
     - `FindByIDWithAvailability(id uint) (*dto.ZoneResponse, error)`
     - `FindByID(id uint) (*models.ParkingZone, error)`
   - Calculation Logic for `available_spots`:
     Use GORM subquery or join query counting active reservations for each zone:
     ```sql
     SELECT 
       pz.id, pz.name, pz.type, pz.total_capacity, pz.price_per_hour, pz.created_at, pz.updated_at,
       (pz.total_capacity - COALESCE((SELECT COUNT(*) FROM reservations r WHERE r.zone_id = pz.id AND r.status = 'active'), 0)) as available_spots
     FROM parking_zones pz
     ```

3. **Zone Service (`service/zone_service.go`)**
   - Interface `ZoneService`:
     - `CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)`
     - `GetAllZones() ([]dto.ZoneResponse, error)`
     - `GetZoneByID(id uint) (*dto.ZoneResponse, error)`

4. **Zone Handler (`handler/zone_handler.go`)**
   - `CreateZone`: Validate DTO, call service, return `201 Created`.
   - `GetAllZones`: Call service, return `200 OK`.
   - `GetZoneByID`: Extract `id` path param, call service, return `200 OK` (or `404 Not Found`).

5. **Route Registration (`cmd/api/main.go`)**
   - `POST /api/v1/zones` → Protected by `JWTMiddleware` + `RequireRole("admin")`.
   - `GET /api/v1/zones` → Public.
   - `GET /api/v1/zones/:id` → Public.

### 💻 Execution Commands
```bash
cd backend
go run cmd/api/main.go
```

### ✅ Step 5 Verification
1. **Create Zone as Admin (Using Admin JWT token from Step 3):**
   ```bash
   curl -X POST http://localhost:8080/api/v1/zones \
     -H "Authorization: Bearer <ADMIN_JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Terminal 1 EV Charging","type":"ev_charging","total_capacity":20,"price_per_hour":5.50}'
   ```
   Expected: `201 Created` with zone data.

2. **Attempt Create Zone as Driver:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/zones \
     -H "Authorization: Bearer <DRIVER_JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Unauth Zone","type":"general","total_capacity":10,"price_per_hour":2.00}'
   ```
   Expected: `403 Forbidden`.

3. **Get All Zones (Public):**
   ```bash
   curl http://localhost:8080/api/v1/zones
   ```
   Expected: `200 OK` returning list of zones with `available_spots: 20`.

---

## 📌 Step 6: Reservations Module (Core Logic & Concurrency Control `FOR UPDATE`)

### 🎯 Objective
Implement the core reservation logic with **Row-Level DB Locking (`FOR UPDATE`)** inside a GORM Database Transaction to resolve the EV Spot race condition. Implement driver reservation, viewing driver reservations, cancelling reservations, and admin viewing all reservations.

### 📋 Prerequisites
- Step 5 completed.

### 📁 Files to Create / Modify
- `backend/dto/reservation_dto.go`
- `backend/repository/reservation_repository.go`
- `backend/service/reservation_service.go`
- `backend/handler/reservation_handler.go`
- `backend/cmd/api/main.go` (Wire reservation endpoints)

### 📝 Step Context & Implementation Details

1. **Reservation DTOs (`dto/reservation_dto.go`)**
   - `CreateReservationRequest`:
     - `ZoneID` (uint `json:"zone_id" validate:"required"`)
     - `LicensePlate` (string `json:"license_plate" validate:"required,min=2,max=15"`)
   - `ReservationResponse`:
     - `ID` (uint `json:"id"`)
     - `UserID` (uint `json:"user_id"`)
     - `ZoneID` (uint `json:"zone_id"`)
     - `LicensePlate` (string `json:"license_plate"`)
     - `Status` (string `json:"status"`)
     - `Zone` (ZoneBriefResponse `json:"zone,omitempty"`)
     - `User` (UserBriefResponse `json:"user,omitempty"`)
     - `CreatedAt` (time.Time `json:"created_at"`)
     - `UpdatedAt` (time.Time `json:"updated_at"`)

2. **CRITICAL Concurrency Implementation in Reservation Repository (`repository/reservation_repository.go`)**
   Implement `CreateReservation(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error)`:
   ```go
   var reservation models.Reservation
   
   err := r.db.Transaction(func(tx *gorm.DB) error {
       var zone models.ParkingZone
       // 1. ROW-LEVEL LOCKING: Lock parking zone record using FOR UPDATE
       if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
           if errors.Is(err, gorm.ErrRecordNotFound) {
               return errors.New("parking zone not found")
           }
           return err
       }
   
       // 2. Count active reservations for this zone
       var activeCount int64
       if err := tx.Model(&models.Reservation{}).
           Where("zone_id = ? AND status = ?", zoneID, "active").
           Count(&activeCount).Error; err != nil {
           return err
       }
   
       // 3. CAPACITY CHECK: If activeCount >= zone.TotalCapacity, return custom conflict error
       if int(activeCount) >= zone.TotalCapacity {
           return errors.New("parking zone is full") // Maps to 409 Conflict
       }
   
       // 4. Create new reservation record
       reservation = models.Reservation{
           UserID:       userID,
           ZoneID:       zoneID,
           LicensePlate: licensePlate,
           Status:       "active",
       }
       if err := tx.Create(&reservation).Error; err != nil {
           return err
       }
   
       return nil
   })
   ```

3. **Additional Repository Methods**
   - `FindByUserID(userID uint) ([]models.Reservation, error)` with `Preload("Zone")`.
   - `FindByID(id uint) (*models.Reservation, error)`.
   - `CancelReservation(id uint, userID uint, isAdmin bool) error`:
     - Checks if reservation exists.
     - If not admin and `reservation.UserID != userID`, return forbidden error.
     - Update `status` to `"cancelled"`.
   - `FindAll() ([]models.Reservation, error)` with `Preload("Zone")` and `Preload("User")`.

4. **Reservation Service (`service/reservation_service.go`)**
   - `CreateReservation(userID uint, req dto.CreateReservationRequest)`: Call repo transaction. Map full error string ("parking zone is full" -> 409 status code).
   - `GetMyReservations(userID uint)`: Fetch driver reservations.
   - `CancelReservation(id uint, userID uint, userRole string)`: Check permission & cancel.
   - `GetAllReservations()`: Admin view of all reservations.

5. **Reservation Handler (`handler/reservation_handler.go`)**
   - `CreateReservation`: Bind & validate DTO, get `userID` from Echo context, call service.
   - `GetMyReservations`: Get `userID` from Echo context, call service.
   - `CancelReservation`: Get `id` from URL param, `userID` & `role` from context, call service.
   - `GetAllReservations`: Admin check, call service.

6. **Route Wiring (`cmd/api/main.go`)**
   - `POST /api/v1/reservations` → Auth (`JWTMiddleware`)
   - `GET /api/v1/reservations/my-reservations` → Auth (`JWTMiddleware`)
   - `DELETE /api/v1/reservations/:id` → Auth (`JWTMiddleware`)
   - `GET /api/v1/reservations` → Auth (`JWTMiddleware`) + Admin (`RequireRole("admin")`)

### 💻 Execution Commands
```bash
cd backend
go run cmd/api/main.go
```

### ✅ Step 6 Verification
1. **Driver Booking:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/reservations \
     -H "Authorization: Bearer <DRIVER_JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"zone_id":1,"license_plate":"EV-9999"}'
   ```
   Expected: `201 Created` with reservation details.

2. **Get My Reservations:**
   ```bash
   curl http://localhost:8080/api/v1/reservations/my-reservations \
     -H "Authorization: Bearer <DRIVER_JWT_TOKEN>"
   ```
   Expected: `200 OK` with list containing zone details.

3. **Check Dynamic Capacity decrease:**
   ```bash
   curl http://localhost:8080/api/v1/zones/1
   ```
   Expected: `available_spots` reduced by 1 (e.g. 19/20).

4. **Cancel Reservation:**
   ```bash
   curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
     -H "Authorization: Bearer <DRIVER_JWT_TOKEN>"
   ```
   Expected: `200 OK` with `"Reservation cancelled successfully"`.

---

## 📌 Step 7: Centralized Error Handler, CORS & Final Verification

### 🎯 Objective
Configure standard Echo central HTTP error handler to catch invalid JSON, route 404s, and internal errors into standard JSON format. Add CORS middleware for frontend communication, and verify complete backend operation.

### 📋 Prerequisites
- Step 6 completed.

### 📁 Files to Create / Modify
- `backend/utils/error_handler.go`
- `backend/cmd/api/main.go` (Attach error handler & CORS middleware)

### 📝 Step Context & Implementation Details

1. **Central HTTP Error Handler (`utils/error_handler.go`)**
   ```go
   func CustomHTTPErrorHandler(err error, c echo.Context) {
       code := http.StatusInternalServerError
       msg := "Internal Server Error"
       errDetail := err.Error()

       if he, ok := err.(*echo.HTTPError); ok {
           code = he.Code
           if m, ok := he.Message.(string); ok {
               msg = m
           } else {
               msg = fmt.Sprintf("%v", he.Message)
           }
       }

       if !c.Response().Committed {
           c.JSON(code, map[string]interface{}{
               "success": false,
               "message": msg,
               "errors":  errDetail,
           })
       }
   }
   ```

2. **CORS & Middleware Assembly (`cmd/api/main.go`)**
   Add middleware:
   ```go
   e.Use(middleware.Logger())
   e.Use(middleware.Recover())
   e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
       AllowOrigins: []string{"*"},
       AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
       AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
   }))
   e.HTTPErrorHandler = utils.CustomHTTPErrorHandler
   ```

### 💻 Execution Commands
```bash
cd backend
go build -o spotsync-api cmd/api/main.go
./spotsync-api
```

### ✅ Step 7 Verification Checklist
- `POST /api/v1/auth/register` (Driver & Admin)
- `POST /api/v1/auth/login` (Driver & Admin tokens)
- `POST /api/v1/zones` (Admin creates zone)
- `GET /api/v1/zones` (Public view with dynamic `available_spots`)
- `POST /api/v1/reservations` (Driver reserves spot)
- `GET /api/v1/reservations/my-reservations` (Driver views bookings)
- `DELETE /api/v1/reservations/:id` (Driver cancels booking)
- `GET /api/v1/reservations` (Admin views all system bookings)
- Verify `FOR UPDATE` transaction locks when zone is filled to capacity (21st reservation returns `409 Conflict`).


// chatgpt
Do not ask for confirmation unless the step is ambiguous or missing information.

Do not skip any subtask.

If verification fails, fix the implementation and rerun verification until it passes.

Only declare the step complete after all verification commands succeed.