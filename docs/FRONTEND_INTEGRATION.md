# Frontend Integration Guide

## Configuration

```
API Base URL: http://localhost:8080
Frontend URL: http://localhost:3000
```

**Important:** Always set `credentials: 'include'` in fetch requests for session cookies.

---

## API Endpoints

### Authentication

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| POST | `/auth/register` | `{email, username, password}` | Register new user |
| POST | `/auth/login` | `{email_or_username, password}` | Login user |
| POST | `/auth/logout` | - | Logout user |
| GET | `/auth/google/login` | - | Initiate Google OAuth |
| POST | `/auth/verify-email` | `{token}` | Verify email |
| POST | `/auth/resend-verification` | `{email}` | Resend verification email |

### Protected Endpoints (require session)

| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| GET | `/api/me` | - | Get current user |
| GET | `/api/validate` | - | Validate session |
| POST | `/api/query-agent` | `{query}` | Natural language to SQL |
| GET | `/api/query-history` | - | Get query history |

---

## Auth Flow

```
Registration:
1. POST /auth/register → redirect to /verify-email
2. User clicks email link → /verify-email?token=xxx
3. POST /auth/verify-email → redirect to /login

Login:
1. POST /auth/login → sets session cookie
2. Response: {message, redirect, user}

Google OAuth:
1. GET /auth/google/login → redirect to Google
2. Callback → redirect to /dashboard

Logout:
1. POST /auth/logout → clears session
```

---

## Data Types

```typescript
interface User {
  ID: number;
  email: string;
  username: string;
  role: 'user' | 'admin';
  provider: 'local' | 'google';
  avatar_url: string | null;
  email_verified: boolean;
}

interface RegisterResponse { message, redirect, user }
interface LoginResponse { message, redirect, user }
interface QueryResponse { response: string }  // SQL query
interface QueryHistoryResponse { queries: [{ID, Query, Response, CreatedAt}] }
```

---

## Key Implementation Notes

1. **Session Cookie:** HttpOnly, name=`session`, automatically managed by browser
2. **Protected Routes:** Check `/api/validate` or `/api/me` on app init
3. **Email Verification:** Local users must verify before accessing protected routes
4. **Error Handling:** API returns `{error: string}` on failure
5. **CORS:** Accepts `http://localhost:*` by default

---

## Quick cURL Examples

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"test","password":"pass123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email_or_username":"test@example.com","password":"pass123"}'

# Query Agent (with session)
curl -X POST http://localhost:8080/api/query-agent \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"query":"Show all users"}'
```
