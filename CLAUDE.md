# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Information

This is a Go project: `github.com/ia-edev-sindireceita/todo`
- Go version: 1.24.5
- Frontend: HTMX + Tailwind CSS
- Architecture: Hexagonal Architecture
- Development approach: Test-Driven Development (TDD)

## Workflow

### Branching Strategy

**MANDATORY**: All development must follow this workflow:

1. **Create an issue first** - NEVER make code changes without a corresponding GitHub issue
   - Document what needs to be done
   - Explain the problem or feature request
   - Reference the issue number in all related work

2. **Create a feature branch** for each task/issue:
   ```bash
   git checkout -b feature/descriptive-name
   # or
   git checkout -b fix/bug-description
   ```

3. **Work on the branch** following TDD principles

4. **Commit changes** with clear, descriptive messages

5. **Create a Pull Request** to `main` when done:
   ```bash
   gh pr create --base main --head feature/descriptive-name
   ```

6. **Never commit directly to `main`**

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code refactoring
- `docs/` - Documentation updates
- `test/` - Test additions or fixes

## Commands

### Build
```bash
export PATH=$PATH:/usr/local/go/bin/ && go build
```

### Run
```bash
export PATH=$PATH:/usr/local/go/bin/ && go run cmd/server/main.go
```

### Test
```bash
# Run all tests
export PATH=$PATH:/usr/local/go/bin/ && go test ./...

# Run tests with verbose output
export PATH=$PATH:/usr/local/go/bin/ && go test -v ./...

# Run a specific test
export PATH=$PATH:/usr/local/go/bin/ && go test -run TestName ./...
```

### Dependencies
```bash
# Add a dependency
export PATH=$PATH:/usr/local/go/bin/ && go get <package>

# Tidy dependencies
export PATH=$PATH:/usr/local/go/bin/ && go mod tidy
```

### Formatting and Linting
```bash
# Format code
export PATH=$PATH:/usr/local/go/bin/ && go fmt ./...

# Vet code
export PATH=$PATH:/usr/local/go/bin/ && go vet ./...
```

## Architecture

### Hexagonal Architecture (Ports and Adapters)

The project follows hexagonal architecture with clear separation between domain and infrastructure:

**Domain Layer** (core business logic, no external dependencies):
- `application/` - Domain entities and value objects. Entities validate data before persisting to database
- `service/` - General business rules and domain services
- `repository/` - Repository interfaces (ports)

**Application Layer**:
- `usecases/` - Specific business use cases that orchestrate services and repositories

**Infrastructure Layer** (adapters):
- Repository implementations
- HTTP handlers
- Database connections
- External service integrations

### Dependency Flow
```
usecases → services → repositories
         ↓
    application (entities validate data)
```

Use cases orchestrate services and repositories. Services contain general business logic. Repositories handle data persistence. Entities in `application/` validate data before it reaches the database.

### Frontend

- Use **HTMX** for dynamic interactions without writing JavaScript
- Use **Tailwind CSS** for styling
- Server-side rendering with Go templates
- HTMX endpoints should return HTML fragments

## Security Practices

### Input Validation and Sanitization
- **Entities are the first line of defense**: all validation happens in `application/` layer
- Validate type, length, format, and business rules in entity constructors/factories
- Return explicit validation errors, never silent failures
- Sanitize all user input before rendering in templates (use `html/template` auto-escaping)
- Never trust client-side validation alone

### Authentication and Authorization
- Implement authentication in middleware layer
- Use secure session management (httponly, secure, samesite cookies)
- Store passwords using bcrypt or argon2
- Implement proper authorization checks in use cases
- Never expose sensitive data in error messages or logs

### Database Security
- **ALWAYS use prepared statements/parameterized queries** in repository implementations
- Never construct SQL with string concatenation
- Apply principle of least privilege to database users
- Encrypt sensitive data at rest
- Use connection pooling securely
- Version control all migrations

### Web Security (HTMX)
- **CSRF Protection**: include CSRF tokens in all mutating requests (POST, PUT, DELETE)
- Set security headers:
  - `Content-Security-Policy`
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `Strict-Transport-Security`
- Validate all HTMX requests server-side
- Sanitize HTML fragments before returning
- Use HTTPS only in production

### Secrets Management
- Store secrets in environment variables, never hardcode
- Use `.env` files for local development (add to `.gitignore`)
- Rotate secrets regularly
- Minimum necessary access for API keys and tokens

### Error Handling and Logging
- Log security events (failed auth, suspicious activity)
- Never log sensitive information (passwords, tokens, PII)
- Return generic error messages to clients
- Fail securely: on error, deny access by default

### Rate Limiting and DoS Protection
- Implement rate limiting on public endpoints
- Set request timeout limits
- Limit request body size
- Implement circuit breakers for external services

## Design Principles (Minimalism)

### HTMX Best Practices
- Use minimal attributes: `hx-get`, `hx-post`, `hx-swap`, `hx-target`
- Return HTML fragments (partials), not full pages
- Progressive enhancement: basic functionality works without JS
- Avoid complex client-side logic
- Use semantic HTTP methods (GET for reads, POST/PUT/DELETE for writes)
- Keep responses small and focused

### Tailwind CSS Guidelines
- Utility-first approach: use Tailwind classes directly
- Avoid custom CSS unless absolutely necessary
- Create reusable components for repeated patterns
- Use Tailwind's design system (spacing, colors, typography)
- Purge unused CSS in production
- Keep markup semantic and accessible

### UI/UX Minimalism
- Simple, clean interfaces
- Clear visual hierarchy
- Minimal color palette
- Consistent spacing and typography
- Accessible by default (ARIA labels, keyboard navigation)
- Mobile-first responsive design

## Go Development Best Practices

### Test-Driven Development (TDD)

**MANDATORY**: Follow TDD for all development:
1. Write a test that fails (Red)
2. Write the minimum code to make the test pass (Green)
3. Refactor if needed (Refactor)

All actions must be preceded by a failing test. The test must fail on first execution and pass after implementation.

### Code Organization
- Entities in `application/` must validate all business rules before data persistence
- Entity constructors/factories should return `(Entity, error)` to enforce validation
- Use value objects for domain concepts (e.g., Email, CPF, Money)
- Repository interfaces belong in the domain; implementations in infrastructure
- Use dependency injection to wire dependencies
- Use cases depend on service and repository interfaces, not concrete implementations
- Keep packages small and focused (Single Responsibility Principle)

### Error Handling
- Return errors explicitly; don't panic except for truly exceptional cases
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Validate data in entities and return validation errors

### Testing
- Write table-driven tests
- Mock repository and service interfaces for use case testing
- Test entity validation logic thoroughly
- Test domain logic in isolation from infrastructure

### Naming Conventions
- Use Go standard naming (PascalCase for exported, camelCase for unexported)
- Interface names: single-method interfaces end with "-er" (e.g., `Repository`, `Service`)
- Avoid stuttering (e.g., `user.UserRepository` → `user.Repository`)

### Dependency Direction

Dependencies flow: **usecases → services → repositories**

All layers use entities from `application/` for data validation. The domain has no dependencies on infrastructure.

### Security Architecture Principles

- **Defense in Depth**: multiple layers of security controls
- **Fail Securely**: default to deny access on errors
- **Least Privilege**: grant minimum necessary permissions
- **Zero Trust**: validate and authenticate everything
- **Security by Default**: secure configurations out of the box
- **Separation of Concerns**: security logic isolated and testable
