package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureCreatesBootstrapConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := Ensure(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Contains(content, "bootstrapAdmin") {
		t.Fatalf("generated config should not include bootstrapAdmin:\n%s", content)
	}
	for _, snippet := range []string{
		"path: ./data/backend.db",
		"path: ./data/media",
	} {
		if !strings.Contains(content, snippet) {
			t.Fatalf("generated config missing %q:\n%s", snippet, content)
		}
	}
	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("config dir mode = %o", info.Mode().Perm())
	}
	if err := Ensure(path); err != nil {
		t.Fatal(err)
	}
}