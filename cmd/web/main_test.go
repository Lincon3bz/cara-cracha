package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
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

func TestCriarEventosAoMesmoTempo(t *testing.T) {
	antes := len(eventos)
	const total = 1000

	var wg sync.WaitGroup
	for range total {
		wg.Add(1)
		go func() {
			defer wg.Done()
			form := url.Values{"nome": {"Festa_teste"}, "data": {"01/01/2020"}}
			req := httptest.NewRequest("POST", "/eventos", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			rotas().ServeHTTP(rec, req)
		}()
	}
	wg.Wait()
	novos := len(eventos) - antes
	if novos != total {
		t.Errorf("esperava %d eventos novos, vieram %d", total, novos)
	}
}
