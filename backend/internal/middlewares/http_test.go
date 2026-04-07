package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestSuccessResponse_WithData(t *testing.T) {
	data := map[string]string{"status": "ok"}
	resp := SuccessResponse(data, http.StatusOK)

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if resp.Headers["Content-Type"] != "application/json" {
		t.Error("missing Content-Type header")
	}
	if resp.Headers["Access-Control-Allow-Origin"] != "*" {
		t.Error("missing CORS header")
	}

	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["status"] != "ok" {
		t.Errorf("body status = %q, want %q", body["status"], "ok")
	}
}

func TestSuccessResponse_NilData(t *testing.T) {
	resp := SuccessResponse(nil, http.StatusNoContent)
	if resp.StatusCode != 204 {
		t.Errorf("StatusCode = %d, want 204", resp.StatusCode)
	}
	if resp.Body != "" {
		t.Errorf("Body = %q, want empty", resp.Body)
	}
}

func TestBadRequestResponse(t *testing.T) {
	resp := BadRequestResponse([]string{"invalid email"})
	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
	var body map[string]interface{}
	json.Unmarshal([]byte(resp.Body), &body)
	errs, ok := body["errors"].([]interface{})
	if !ok || len(errs) != 1 || errs[0] != "invalid email" {
		t.Errorf("errors = %v, want [\"invalid email\"]", body["errors"])
	}
}

func TestNotFoundResponse(t *testing.T) {
	resp := NotFoundResponse("user")
	if resp.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", resp.StatusCode)
	}
	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["error"] != "user not found" {
		t.Errorf("error = %q, want %q", body["error"], "user not found")
	}
}

func TestUnauthorizedResponse(t *testing.T) {
	resp := UnauthorizedResponse()
	if resp.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", resp.StatusCode)
	}
	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["error"] != "unauthorized" {
		t.Errorf("error = %q, want %q", body["error"], "unauthorized")
	}
}

func TestForbiddenResponse(t *testing.T) {
	resp := ForbiddenResponse("access denied")
	if resp.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", resp.StatusCode)
	}
	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["error"] != "access denied" {
		t.Errorf("error = %q, want %q", body["error"], "access denied")
	}
}

func TestServiceErrorResponse(t *testing.T) {
	svcErr := NewConflictError("duplicate entry")
	resp := ServiceErrorResponse(svcErr)
	if resp.StatusCode != 409 {
		t.Errorf("StatusCode = %d, want 409", resp.StatusCode)
	}
	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["error"] != "duplicate entry" {
		t.Errorf("error = %q, want %q", body["error"], "duplicate entry")
	}
}

func TestHandleServiceError_WithServiceError(t *testing.T) {
	err := NewNotFoundError("match not found")
	resp := HandleServiceError(err)
	if resp.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", resp.StatusCode)
	}
}

func TestHandleServiceError_WithGenericError(t *testing.T) {
	resp := HandleServiceError(errors.New("unexpected"))
	if resp.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", resp.StatusCode)
	}
	var body map[string]string
	json.Unmarshal([]byte(resp.Body), &body)
	if body["error"] != "internal server error" {
		t.Errorf("error = %q, want %q", body["error"], "internal server error")
	}
}

func TestAllResponses_HaveCorsHeaders(t *testing.T) {
	responses := []struct {
		name string
		resp func() map[string]string
	}{
		{"Success", func() map[string]string { return SuccessResponse(nil, 200).Headers }},
		{"BadRequest", func() map[string]string { return BadRequestResponse([]string{"err"}).Headers }},
		{"NotFound", func() map[string]string { return NotFoundResponse("x").Headers }},
		{"Unauthorized", func() map[string]string { return UnauthorizedResponse().Headers }},
		{"Forbidden", func() map[string]string { return ForbiddenResponse("x").Headers }},
		{"HandleServiceError", func() map[string]string { return HandleServiceError(errors.New("x")).Headers }},
	}

	for _, tc := range responses {
		t.Run(tc.name, func(t *testing.T) {
			headers := tc.resp()
			if headers["Access-Control-Allow-Origin"] != "*" {
				t.Error("missing Access-Control-Allow-Origin")
			}
			if headers["Access-Control-Allow-Methods"] == "" {
				t.Error("missing Access-Control-Allow-Methods")
			}
		})
	}
}
