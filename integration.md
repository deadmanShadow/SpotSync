# Frontend ↔ Backend Integration with JWT

## Task 1
- JWT Expiration

Configure JWT authentication so that every newly issued JWT token expires after exactly 1 day (24 hours).

Requirements:

Set the JWT expiration (exp) to now + 24 hours.
Apply this consistently to all JWT tokens issued by the backend.
Do not hardcode an already-expired timestamp.
Validate the exp claim when authenticating requests.
Expired tokens must be rejected with the appropriate 401 Unauthorized response.
Keep JWT_SECRET server-side only.
Do not expose JWT_SECRET to the Astro frontend.

Example concept:

expiresAt := time.Now().Add(24 * time.Hour)

The exact implementation should follow the existing JWT library and authentication architecture in the project. Do not unnecessarily refactor the authentication system.



## Task - 2
Configure my existing **Astro frontend** and **Go + Echo backend** so they communicate correctly in production.

## Current Setup

### Go Backend `.env`

```env
APP_PORT=8080
DATABASE_URL=your-database-url
JWT_SECRET=your-jwt-secret
```

Current Go config:

```go
type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}
```

`LoadConfig()` currently loads environment variables using `godotenv`.

### Astro Frontend `.env`

```env
PUBLIC_API_BASE_URL=https://spotsync-s3f4.onrender.com/api/v1
```

## Required Changes

1. Add this to the **Go backend environment/config**:

```env
FRONTEND_URL=http://localhost:4321
```

Add `FrontendURL string` to the existing `Config` struct and load it through the existing `getEnv()` system.

2. Configure **Echo CORS middleware** using `config.FrontendURL`.

Allow the required methods:

```text
GET, POST, PUT, PATCH, DELETE, OPTIONS
```

and headers such as:

```text
Content-Type, Authorization
```

Do not use `*` as the allowed origin.

3. Keep the Astro environment variable:

```env
PUBLIC_API_BASE_URL=https://spotsync-s3f4.onrender.com/api/v1
```

Make sure all frontend API requests use:

```ts
import.meta.env.PUBLIC_API_BASE_URL
```

instead of hardcoded backend URLs.

4. Do not expose any backend secrets to Astro. `DATABASE_URL` and `JWT_SECRET` must remain backend-only.

5. Do not unnecessarily refactor existing code. Inspect the current project structure and modify only what is required.

6. Update `.env.example` / `.gitignore` if necessary.

## Final Check

Verify that:

```text
Astro Frontend
      ↓
PUBLIC_API_BASE_URL
      ↓
Go Backend
      ↓
CORS → FRONTEND_URL
```

works correctly in production without CORS errors.

After implementation, briefly tell me:
- files changed
- environment variables required
- whether backend/frontend redeployment is required
- how to test the connection
