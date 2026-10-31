package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	rotas().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status: esperava %d,veio %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("corpo: esperava %q, veio %q", "ok", rec.Body.String())
	}
}
