# API Reference

This document provides detailed information about all available API endpoints in the IPKISS API implementation.

## Base URL

```
http://localhost:8080
```

## Endpoints

### Reset State

Resets the application state, removing all accounts and balances.

**Endpoint:** `POST /reset`

**Request:**

- Method: `POST`
- URL: `/reset`
- Body: None

**Response:**

- Status: `200 OK`
- Body: `OK`

**Example:**

```bash
curl -X POST http://localhost:8080/reset
```

**Response:**

```
OK
```

---

### Get Balance

Retrieves the balance for a specific account.

**Endpoint:** `GET /balance`

**Request:**

- Method: `GET`
- URL: `/balance?account_id={account_id}`
- Query Parameters:
    - `account_id` (required): The ID of the account

**Response:**

**Success (200 OK):**

- Status: `200 OK`
- Body: Balance as integer (e.g., `20`)

**Error (404 Not Found):**

- Status: `404 Not Found`
- Body: `0`
- Occurs when: Account does not exist

**Error (400 Bad Request):**

- Status: `400 Bad Request`
- Body: `missing account_id`
- Occurs when: `account_id` parameter is not provided

**Examples:**

```bash
# Get balance for existing account
curl http://localhost:8080/balance?account_id=100

# Response: 20
```

```bash
# Get balance for non-existing account
curl http://localhost:8080/balance?account_id=999

# Response (404): 0
```

---

### Process Event

Processes financial events including deposits, withdrawals, and transfers.

**Endpoint:** `POST /event`

**Request:**

- Method: `POST`
- URL: `/event`
- Headers:
    - `Content-Type: application/json`
- Body: JSON object (see event types below)

**Event Types:**

#### 1. Deposit

**Request Body:**

```json
{
    "type": "deposit",
    "destination": "100",
    "amount": 10
}
```

**Fields:**

- `type`: Must be `"deposit"`
- `destination` (required): Account ID to deposit into (numeric string)
- `amount` (required): Amount to deposit (positive integer)

**Response (201 Created):**

```json
{
    "destination": {
        "id": "100",
        "balance": 10
    }
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"deposit", "destination":"100", "amount":10}'
```

---

#### 2. Withdraw

**Request Body:**

```json
{
    "type": "withdraw",
    "origin": "100",
    "amount": 5
}
```

**Fields:**

- `type`: Must be `"withdraw"`
- `origin` (required): Account ID to withdraw from (numeric string)
- `amount` (required): Amount to withdraw (positive integer)

**Success Response (201 Created):**

```json
{
    "origin": {
        "id": "100",
        "balance": 5
    }
}
```

**Error Response (404 Not Found):**

- Status: `404 Not Found`
- Body: `0`
- Occurs when:
    - Origin account does not exist
    - Insufficient funds in origin account

**Example:**

```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"withdraw", "origin":"100", "amount":5}'
```

---

#### 3. Transfer

**Request Body:**

```json
{
    "type": "transfer",
    "origin": "100",
    "destination": "300",
    "amount": 15
}
```

**Fields:**

- `type`: Must be `"transfer"`
- `origin` (required): Account ID to transfer from (numeric string)
- `destination` (required): Account ID to transfer to (numeric string)
- `amount` (required): Amount to transfer (positive integer)

**Success Response (201 Created):**

```json
{
    "origin": {
        "id": "100",
        "balance": 0
    },
    "destination": {
        "id": "300",
        "balance": 15
    }
}
```

**Error Response (404 Not Found):**

- Status: `404 Not Found`
- Body: `0`
- Occurs when:
    - Origin account does not exist
    - Insufficient funds in origin account

**Example:**

```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"transfer", "origin":"100", "destination":"300", "amount":15}'
```

**Note:** If the destination account doesn't exist, it will be created automatically with the transferred amount as the initial balance.

---

## Validation Rules

The `/event` endpoint validates all requests using the following rules:

- **`type`**: Required, must be one of: `deposit`, `withdraw`, `transfer`
- **`amount`**: Required, must be a positive integer (greater than 0)
- **`origin`**:
    - Required for `withdraw` and `transfer` events
    - Must be a numeric string
    - Account must exist for withdrawal/transfer
- **`destination`**:
    - Required for `deposit` and `transfer` events
    - Must be a numeric string
    - Account will be created if it doesn't exist

**Validation Error Response (400 Bad Request):**

When validation fails, the API returns a `400 Bad Request` status with details about the validation errors.

```json
[
    {
        "EventRequest.Type": "Type must be one of [deposit withdraw transfer]"
    }
]
```

---

## Complete Usage Example

Here's a complete workflow demonstrating all API operations:

```bash
# 1. Reset state
curl -X POST http://localhost:8080/reset
# Response: OK

# 2. Get balance for non-existing account
curl http://localhost:8080/balance?account_id=1234567890
# Response (404): 0

# 3. Create account with deposit
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"deposit", "destination":"100", "amount":10}'
# Response: {"destination":{"id":"100","balance":10}}

# 4. Deposit into existing account
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"deposit", "destination":"100", "amount":10}'
# Response: {"destination":{"id":"100","balance":20}}

# 5. Check balance
curl http://localhost:8080/balance?account_id=100
# Response: 20

# 6. Withdraw from account
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"withdraw", "origin":"100", "amount":5}'
# Response: {"origin":{"id":"100","balance":15}}

# 7. Transfer to new account
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{"type":"transfer", "origin":"100", "destination":"300", "amount":15}'
# Response: {"origin":{"id":"100","balance":0},"destination":{"id":"300","balance":15}}

# 8. Verify balances
curl http://localhost:8080/balance?account_id=100
# Response: 0

curl http://localhost:8080/balance?account_id=300
# Response: 15
```

---

## Error Handling Summary

| Status Code       | Description     | When It Occurs                                                        |
| ----------------- | --------------- | --------------------------------------------------------------------- |
| `200 OK`          | Success         | Balance query successful, Reset successful                            |
| `201 Created`     | Success         | Event processed successfully                                          |
| `400 Bad Request` | Invalid request | Missing required parameters, validation errors                        |
| `404 Not Found`   | Not found       | Account doesn't exist (balance/withdraw/transfer), insufficient funds |

**Note on 404 for Insufficient Funds:** The API returns `404 Not Found` with body `0` for both non-existent accounts and insufficient funds. This is part of the IPKISS API specification.
