package cli

import (
	"testing"
	"time"

	"github.com/abiosoft/incus-apply/internal/apply"
	"github.com/abiosoft/incus-apply/internal/renderer"
)

func TestNewRenderer(t *testing.T) {
	// Default options should return TextRenderer.
	opts := &apply.Options{}
	r := newRenderer(opts)
	if _, ok := r.(*renderer.TextRenderer); !ok {
		t.Error("expected TextRenderer for default options")
	}

	// --diff=text should return TextRenderer.
	opts = &apply.Options{Diff: "text"}
	r = newRenderer(opts)
	if _, ok := r.(*renderer.TextRenderer); !ok {
		t.Error("expected TextRenderer for --diff=text")
	}

	// --diff=json should return JSONRenderer.
	opts = &apply.Options{Diff: "json"}
	r = newRenderer(opts)
	if _, ok := r.(*renderer.JSONRenderer); !ok {
		t.Error("expected JSONRenderer for --diff=json")
	}
}

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    apply.Options
		wantErr bool
	}{
		{name: "valid default", opts: apply.Options{}, wantErr: false},
		{name: "valid json diff", opts: apply.Options{Diff: "json"}, wantErr: false},
		{name: "valid replace", opts: apply.Options{Replace: true}, wantErr: false},
		{name: "valid show env", opts: apply.Options{ShowEnv: true}, wantErr: false},
		{name: "invalid diff", opts: apply.Options{Diff: "yaml"}, wantErr: true},
		{name: "negative fetch timeout", opts: apply.Options{FetchTimeout: -time.Second}, wantErr: true},
		{name: "negative command timeout", opts: apply.Options{CommandTimeout: -time.Second}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(&tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
