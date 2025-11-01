# Passwordless Sign-In with GCP Identity Platform

This project demonstrates how to implement passwordless sign-in using Google Cloud Identity Platform (not Firebase Auth) in Go.

## Endpoints

### 1. Send OOB Code

Sends a sign-in link to the user's email.

**Endpoint:** `/sendOobCode`

**Method:** `POST`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Curl Example:**
```sh
curl -X POST http://localhost:8080/sendOobCode \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

---

### 2. Finish Sign In

Verifies the OOB code and completes the sign-in process.

**Endpoint:** `/finishSignIn`

**Method:** `POST`

**Request Body:**
```json
{
  "email": "user@example.com",
  "oobCode": "CODE_FROM_EMAIL"
}
```

**Curl Example:**
```sh
curl -X POST http://localhost:8080/finishSignIn \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "oobCode": "CODE_FROM_EMAIL"}'
```

Replace `user@example.com` and `CODE_FROM_EMAIL` with the actual user email and the OOB code received in the email.

---

### 3. Signup Passwordless

Sends a passwordless sign-up link to the user's email.

**Endpoint:** `/signupPasswordless`

**Method:** `POST`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Curl Example:**
```sh
curl -X POST http://localhost:8080/signupPasswordless \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

---

## Environment Variables

- `FIREBASE_API_KEY`: Your Firebase API key (required)

---

## Running

1. Set your API key:
   ```sh
   export FIREBASE_API_KEY=your_api_key
   ```
2. Run the server:
   ```sh
   go run main.go
   ```
3. Use the endpoints as shown above.
