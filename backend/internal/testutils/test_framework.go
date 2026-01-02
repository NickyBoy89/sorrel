package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	log "log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/NickyBoy89/sorrel/backend/internal/db"
	"github.com/NickyBoy89/sorrel/backend/internal/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbTestFile = "test.db"
	serverPort = 9031
)

type MockFramework struct {
	db     *sql.DB
	mux    *http.ServeMux
	server *http.Server
	client *http.Client
}

func NewMockFramework(registerFunc func(mux *http.ServeMux, middleware middleware.AuthMiddleware), initFunc func(*sql.DB) error) (*MockFramework, error) {
	testdb, err := InitGlobalDB()
	if err != nil {
		return nil, err
	}

	if err := initFunc(testdb); err != nil {
		return nil, err
	}

	m := &MockFramework{
		client: &http.Client{},
		db:     testdb,
	}

	db.DB = testdb

	m.mux = http.NewServeMux()
	registerFunc(m.mux, &middleware.NoAuthHandler{})
	http.DefaultServeMux = m.mux

	m.StartAPIServer()

	return m, err
}

func InitGlobalDB() (*sql.DB, error) {
	testdb, err := sql.Open("sqlite3", dbTestFile)
	if err != nil {
		return nil, err
	}

	return testdb, nil
}

func (m *MockFramework) Close() error {
	if err := m.server.Shutdown(context.Background()); err != nil {
		return err
	}
	if err := db.DB.Close(); err != nil {
		return err
	}

	if err := os.Remove(dbTestFile); err != nil {
		return err
	}

	return nil
}

// Asynchronusly starts a server and returns a handle to it
func (m *MockFramework) StartAPIServer() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: m.mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("error shutting down http server", "error", err)
			return
		}
	}()

	m.server = srv
}

func localURL(relativeUrl string) string {
	return fmt.Sprintf("http://localhost:%d%s", serverPort, relativeUrl)
}

func (m *MockFramework) LocalAPIRequest(req *http.Request) (*http.Response, error) {
	return m.client.Do(req)
}

func (m *MockFramework) CheckedRequest(t *testing.T, method, relativeUrl string, body io.Reader) *http.Response {
	resp := m.Request(t, method, relativeUrl, body)

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		t.Error("request did not return successfully", "code", resp.Status, "message", string(message))
	}

	return resp
}

func (m *MockFramework) Request(t *testing.T, method, relativeUrl string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, localURL(relativeUrl), body)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := m.LocalAPIRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func (m *MockFramework) Get(t *testing.T, relativeUrl string) *http.Response {
	return m.Request(t, http.MethodGet, relativeUrl, nil)
}

func (m *MockFramework) CheckedPost(t *testing.T, relativeUrl string, body io.Reader) *http.Response {
	return m.CheckedRequest(t, http.MethodPost, relativeUrl, body)
}

func (m *MockFramework) CheckedGet(t *testing.T, relativeUrl string) *http.Response {
	return m.CheckedRequest(t, http.MethodGet, relativeUrl, nil)
}
