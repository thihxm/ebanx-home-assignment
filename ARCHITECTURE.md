# Architecture

This document describes the architecture and design decisions of the IPKISS API implementation.

## Overview

The application follows a clean, layered architecture with clear separation of concerns. It's designed to be simple, testable, and maintainable while meeting all requirements of the IPKISS API specification.

## Architecture Layers

```
┌─────────────────────────────────────┐
│         HTTP Handler Layer          │
│    (internal/handler/http.go)       │
│  Request/Response, Validation       │
└─────────────┬───────────────────────┘
              │
              ↓
┌─────────────────────────────────────┐
│         Service Layer               │
│  (internal/service/*.go)            │
│  Business Logic, Orchestration      │
└─────────────┬───────────────────────┘
              │
              ↓
┌─────────────────────────────────────┐
│       Repository Layer              │
│  (internal/repository/*.go)         │
│  Data Access, Persistence           │
└─────────────┬───────────────────────┘
              │
              ↓
┌─────────────────────────────────────┐
│         Domain Layer                │
│    (internal/domain/*.go)           │
│  Models, Interfaces                 │
└─────────────────────────────────────┘
```

## Component Details

### 1. Domain Layer (`internal/domain`)

**Purpose:** Defines the core business entities and contracts (interfaces) that other layers must implement.

**Files:**

- `account.go`: Account model and interfaces
- `event.go`: Event request/response models and service interface

**Key Components:**

- **`Account`**: Core entity representing a bank account

    ```go
    type Account struct {
        ID      string `json:"id"`
        Balance int    `json:"balance"`
    }
    ```

- **`AccountRepository` Interface**: Defines data access contract
    - `FindByID(id string) (*Account, error)`
    - `Upsert(account *Account) (*Account, error)`
    - `Reset() error`

- **`AccountService` Interface**: Defines business logic contract
    - `GetBalance(id string) (int, error)`
    - `Deposit(id string, amount int) (*Account, error)`
    - `Withdraw(id string, amount int) (*Account, error)`
    - `Transfer(originID, destinationID string, amount int) (origin, destination *Account, err error)`
    - `Reset() error`

- **`EventService` Interface**: Defines event processing contract
    - `ProcessEvent(event EventRequest) (*EventResponse, error)`

**Design Decision:** Interfaces are defined in the domain layer to enforce dependency inversion. Higher-level modules (services) don't depend on lower-level modules (repositories); both depend on abstractions.

### 2. Repository Layer (`internal/repository`)

**Purpose:** Handles data persistence and retrieval. Implements the `AccountRepository` interface.

**Files:**

- `in_memory.go`: In-memory implementation of AccountRepository
- `in_memory_test.go`: Tests for the repository

**Key Components:**

- **`InMemoryRepository`**: Thread-safe in-memory storage
    ```go
    type InMemoryRepository struct {
        accounts map[string]*domain.Account
        mu       sync.RWMutex
    }
    ```

**Thread Safety:** Uses `sync.RWMutex` for concurrent access:

- Read operations (`FindByID`) use `RLock()` for concurrent reads
- Write operations (`Upsert`, `Reset`) use `Lock()` for exclusive access

**Design Decision:** In-memory implementation keeps the solution simple while maintaining production-quality patterns. The repository pattern makes it easy to swap implementations (e.g., to PostgreSQL) without changing other layers.

### 3. Service Layer (`internal/service`)

**Purpose:** Implements business logic and orchestrates operations. Acts as a bridge between handlers and repositories.

**Files:**

- `account_service.go`: Account operations implementation
- `account_service_test.go`: Tests for account service
- `event_service.go`: Event processing implementation
- `event_service_test.go`: Tests for event service

**Key Components:**

#### AccountService

Implements core banking operations:

- **`GetBalance`**: Retrieves account balance
    - Returns error if account doesn't exist

- **`Deposit`**: Adds money to an account
    - Creates account if it doesn't exist (with 0 initial balance)
    - Increases balance by amount

- **`Withdraw`**: Removes money from an account
    - Validates account exists
    - Validates sufficient funds
    - Decreases balance by amount

- **`Transfer`**: Moves money between accounts
    - Validates origin account exists
    - Validates sufficient funds in origin
    - Creates destination account if it doesn't exist
    - Atomically updates both accounts

#### EventService

Orchestrates event processing by delegating to AccountService based on event type:

```go
func (s *EventService) ProcessEvent(event domain.EventRequest) (*domain.EventResponse, error) {
    switch event.Type {
    case "deposit":
        // Delegate to AccountService.Deposit
    case "withdraw":
        // Delegate to AccountService.Withdraw
    case "transfer":
        // Delegate to AccountService.Transfer
    }
}
```

**Design Decision:** Separating AccountService and EventService provides:

- Single Responsibility: Each service has one reason to change
- Reusability: AccountService can be used independently of events
- Testability: Each service can be tested in isolation

### 4. Handler Layer (`internal/handler`)

**Purpose:** HTTP interface to the application. Handles routing, request parsing, validation, and response formatting.

**Files:**

- `http.go`: HTTP handler implementation
- `http_test.go`: Tests for HTTP handlers

**Key Components:**

- **Request Validation**: Uses `go-playground/validator` with struct tags

    ```go
    type EventRequest struct {
        Type        string `json:"type" validate:"required,oneof=deposit withdraw transfer"`
        Origin      string `json:"origin,omitempty" validate:"omitempty,required_if=Type withdraw,required_if=Type transfer,numeric"`
        Destination string `json:"destination,omitempty" validate:"omitempty,required_if=Type deposit,required_if=Type transfer,numeric"`
        Amount      int    `json:"amount" validate:"required,gt=0"`
    }
    ```

- **Internationalization**: Supports error message translation via Accept-Language header

- **Endpoints**:
    - `POST /reset`: Resets application state
    - `GET /balance`: Returns account balance
    - `POST /event`: Processes financial events

**Design Decision:** Validation in the handler layer keeps the domain layer clean and focused on business logic. The handler is the boundary layer responsible for ensuring only valid data reaches the core.

## Data Flow

### Example Flow: Processing a Deposit Event

```
1. HTTP Handler
   ↓ Receives POST /event request
   ↓ Parses JSON body
   ↓ Validates request structure

2. Event Service
   ↓ Receives EventRequest
   ↓ Routes to Deposit operation

3. Account Service
   ↓ Receives account ID and amount
   ↓ Calls repository to find account

4. Repository
   ↓ Locks for reading (RLock)
   ↓ Returns account or nil
   ↑ Unlocks (RUnlock)

5. Account Service
   ↓ Creates new account if nil
   ↓ Updates balance
   ↓ Calls repository to save

6. Repository
   ↓ Locks for writing (Lock)
   ↓ Saves account to map
   ↑ Returns saved account
   ↑ Unlocks (Unlock)

7. Account Service
   ↑ Returns updated account

8. Event Service
   ↑ Wraps in EventResponse

9. HTTP Handler
   ↑ Serializes to JSON
   ↑ Returns 201 Created
```

## Design Patterns

### 1. Repository Pattern

**Purpose:** Abstracts data access logic from business logic.

**Benefits:**

- Easy to swap implementations (in-memory → database)
- Testable: Can mock repository in service tests
- Clear separation of concerns

**Implementation:**

```go
// Domain defines interface
type AccountRepository interface {
    FindByID(id string) (*Account, error)
    Upsert(account *Account) (*Account, error)
    Reset() error
}

// Repository layer implements it
type InMemoryRepository struct {
    accounts map[string]*domain.Account
    mu       sync.RWMutex
}
```

### 2. Dependency Injection

**Purpose:** Loosely couples components and improves testability.

**Implementation:**

```go
func main() {
    // Create dependencies
    repo := repository.NewInMemoryRepository()
    accountService := service.NewAccountService(repo)
    eventService := service.NewEventService(accountService)

    // Inject into handler
    httpHandler := handler.NewAccountHTTPHandler(accountService, eventService)
}
```

**Benefits:**

- Components don't create their own dependencies
- Easy to test with mock implementations
- Clear dependency graph

### 3. Layered Architecture

**Purpose:** Separates concerns and enforces unidirectional dependencies.

**Dependency Flow:**

```
Handler → Service → Repository → Domain
                                    ↑
                All layers depend on Domain interfaces
```

**Benefits:**

- Changes to one layer don't affect others
- Can test each layer independently
- Clear responsibilities

## Concurrency & Thread Safety

### Repository Level

The `InMemoryRepository` uses `sync.RWMutex` to ensure thread-safe concurrent access:

```go
// Read lock: Multiple goroutines can read simultaneously
func (r *InMemoryRepository) FindByID(id string) (*domain.Account, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // Read operation
}

// Write lock: Exclusive access for modifications
func (r *InMemoryRepository) Upsert(account *domain.Account) (*domain.Account, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    // Write operation
}
```

**Why RWMutex?**

- Allows multiple concurrent readers (balance queries)
- Ensures exclusive access for writers (deposits, withdrawals, transfers)
- Better performance than regular Mutex for read-heavy workloads

### HTTP Handler Level

Go's `http.Server` handles each request in a separate goroutine automatically. The repository's thread safety ensures correct behavior under concurrent requests.

## Testing Strategy

### Unit Tests

Each layer has corresponding test files:

- `handler/http_test.go`: Tests HTTP endpoints
- `service/account_service_test.go`: Tests account operations
- `service/event_service_test.go`: Tests event processing
- `repository/in_memory_test.go`: Tests data access

### Test Structure

Tests follow the Arrange-Act-Assert pattern:

```go
func TestDeposit(t *testing.T) {
    // Arrange: Set up dependencies
    repo := NewInMemoryRepository()
    service := NewAccountService(repo)

    // Act: Perform operation
    account, err := service.Deposit("100", 10)

    // Assert: Verify results
    assert.NoError(t, err)
    assert.Equal(t, 10, account.Balance)
}
```

## Technology Choices

### Go Standard Library

**Why:** Minimal dependencies, rock-solid stability, excellent performance.

- `net/http`: HTTP server and routing
- `encoding/json`: JSON parsing and serialization
- `sync`: Thread synchronization primitives

### go-playground/validator

**Why:** Industry-standard validation library with declarative struct tags.

**Features Used:**

- Field validation (`required`, `oneof`, `gt`)
- Conditional validation (`required_if`)
- Custom error messages with i18n

### No External Database

**Why:** Requirement specifies simplest implementation.

**Trade-offs:**

- ✅ Simple: No database setup required
- ✅ Fast: In-memory operations
- ✅ Easy to test: Tests don't need database cleanup
- ❌ Not persistent: Data lost on restart
- ❌ No transactions: But not needed for this use case

## Project Structure

```
ebanx-home-assignment/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   ├── account.go           # Account model & interfaces
│   │   └── event.go             # Event models & interfaces
│   ├── handler/
│   │   ├── http.go              # HTTP handlers
│   │   └── http_test.go         # Handler tests
│   ├── repository/
│   │   ├── in_memory.go         # In-memory repository
│   │   └── in_memory_test.go    # Repository tests
│   └── service/
│       ├── account_service.go   # Account business logic
│       ├── account_service_test.go
│       ├── event_service.go     # Event processing
│       └── event_service_test.go
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── README.md                    # Project overview
├── API_REFERENCE.md            # API documentation
└── ARCHITECTURE.md             # This file
```

### Why `internal/` Package?

The `internal/` package is a Go convention that prevents other projects from importing these packages. It signals that this code is for internal use only and should not be considered a public API.

## Scalability Considerations

### Current Implementation

The in-memory implementation is suitable for:

- Development and testing
- Small-scale deployments
- Single-server scenarios

### Future Enhancements

To scale this application, consider:

1. **Database Persistence**
    - Implement `AccountRepository` with PostgreSQL/MySQL
    - Add transaction support for transfers
    - No changes needed to service or handler layers (thanks to repository pattern!)

2. **Distributed Deployment**
    - Add Redis for distributed locking
    - Implement distributed transactions
    - Consider event sourcing for audit trails

3. **Performance**
    - Add caching layer (Redis)
    - Implement connection pooling
    - Add database read replicas

4. **Observability**
    - Add structured logging
    - Implement metrics (Prometheus)
    - Add distributed tracing (OpenTelemetry)

## Security Considerations

### Current Implementation

- Input validation prevents invalid data
- Thread-safe operations prevent race conditions
- No authentication (not required for this assignment)

### Production Considerations

For a production system, add:

- Authentication & authorization
- Rate limiting
- Audit logging
- HTTPS/TLS
- Input sanitization
- SQL injection prevention (when using database)
- CORS configuration

## Conclusion

This architecture demonstrates:

- ✅ Clean separation of concerns
- ✅ Dependency inversion principle
- ✅ Easy testability
- ✅ Thread safety
- ✅ Maintainability
- ✅ Extensibility

The layered architecture with clear interfaces makes it easy to:

- Add new features
- Swap implementations
- Test in isolation
- Understand the codebase

While simple, the architecture follows production-quality patterns and could scale to more complex requirements with minimal refactoring.
