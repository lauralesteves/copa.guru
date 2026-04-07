package middlewares

import (
	"errors"
	"testing"
)

func TestServiceError_Error_WithWrapped(t *testing.T) {
	err := NewInternalError("db failed", errors.New("connection refused"))
	want := "db failed: connection refused"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestServiceError_Error_WithoutWrapped(t *testing.T) {
	err := NewValidationError("invalid email")
	if got := err.Error(); got != "invalid email" {
		t.Errorf("Error() = %q, want %q", got, "invalid email")
	}
}

func TestServiceError_Unwrap(t *testing.T) {
	original := errors.New("original")
	err := NewInternalError("wrapped", original)
	if !errors.Is(err, original) {
		t.Error("Unwrap() should return the original error")
	}
}

func TestNewServiceError(t *testing.T) {
	err := NewServiceError(ErrConflict, "duplicate", nil)
	if err.Code != ErrConflict {
		t.Errorf("Code = %d, want %d", err.Code, ErrConflict)
	}
	if err.Message != "duplicate" {
		t.Errorf("Message = %q, want %q", err.Message, "duplicate")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("bad input")
	if err.Code != ErrValidation {
		t.Errorf("Code = %d, want %d", err.Code, ErrValidation)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("invalid token")
	if err.Code != ErrUnauthorized {
		t.Errorf("Code = %d, want %d", err.Code, ErrUnauthorized)
	}
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("access denied")
	if err.Code != ErrForbidden {
		t.Errorf("Code = %d, want %d", err.Code, ErrForbidden)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("user not found")
	if err.Code != ErrNotFound {
		t.Errorf("Code = %d, want %d", err.Code, ErrNotFound)
	}
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("already exists")
	if err.Code != ErrConflict {
		t.Errorf("Code = %d, want %d", err.Code, ErrConflict)
	}
}

func TestNewInternalError(t *testing.T) {
	original := errors.New("boom")
	err := NewInternalError("something broke", original)
	if err.Code != ErrInternal {
		t.Errorf("Code = %d, want %d", err.Code, ErrInternal)
	}
	if err.Err != original {
		t.Error("Err should be the original error")
	}
}
