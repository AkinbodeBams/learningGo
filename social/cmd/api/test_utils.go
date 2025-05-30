package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akinbodeBams/social/internal/store"
	"github.com/akinbodeBams/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application{
	t.Helper()
logger := zap.Must(zap.NewProduction()).Sugar()
mockStore := store.NewMockStorage()
mockCacheStore := cache.NewMockStorage()
	return &application{
logger: logger,
store: mockStore,
cacheStorage: mockCacheStore,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder{
	rr := httptest.NewRecorder()
mux.ServeHTTP(rr,req)

return rr
}