package config

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseStdinSetsSourceFile(t *testing.T) {
	input := `---
type: vars
vars:
  NAME: world
---
type: instance
name: web
image: images:alpine/3.19
`

	result, err := NewParser(0).ParseStdin(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseStdin() error = %v", err)
	}
	if result.SourceFile != "stdin" {
		t.Fatalf("SourceFile = %q, want %q", result.SourceFile, "stdin")
	}
	if len(result.Vars) != 1 || result.Vars[0].SourceFile != "stdin" {
		t.Fatalf("vars source file = %#v, want stdin", result.Vars)
	}
	if len(result.Resources) != 1 || result.Resources[0].SourceFile != "stdin" {
		t.Fatalf("resources source file = %#v, want stdin", result.Resources)
	}
}

func TestParseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("type: instance\nname: web\nimage: images:alpine/3.19\n"))
	}))
	defer server.Close()

	result, err := NewParser(0).ParseURL(server.URL)
	if err != nil {
		t.Fatalf("ParseURL() error = %v", err)
	}
	if result.SourceFile != server.URL {
		t.Fatalf("SourceFile = %q, want %q", result.SourceFile, server.URL)
	}
	if len(result.Resources) != 1 {
		t.Fatalf("resources = %d, want 1", len(result.Resources))
	}
}

func TestParseURLTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte("type: instance\nname: web\nimage: images:alpine/3.19\n"))
	}))
	defer server.Close()

	_, err := NewParser(10 * time.Millisecond).ParseURL(server.URL)
	if err == nil {
		t.Fatal("ParseURL() error = nil, want timeout error")
	}
}
