package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContactSubmissionIsStored(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "index.html")
	submissionsPath := filepath.Join(tmpDir, "submissions.jsonl")

	if err := os.WriteFile(templatePath, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	app, err := NewApp(templatePath, "web/static", submissionsPath)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	form := url.Values{
		"name":    {"Mary"},
		"email":   {"mary@example.com"},
		"message": {"Hello from test"},
	}

	req := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	app.contact(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	body, err := os.ReadFile(submissionsPath)
	if err != nil {
		t.Fatalf("read submissions: %v", err)
	}

	if !strings.Contains(string(body), `"email":"mary@example.com"`) {
		t.Fatalf("submission was not persisted: %s", string(body))
	}
}
