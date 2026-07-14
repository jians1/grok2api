package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeriveJWTSecretIsDeterministic(t *testing.T) {
	key := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	first, err := deriveJWTSecret(key)
	if err != nil {
		t.Fatal(err)
	}
	second, err := deriveJWTSecret(key)
	if err != nil {
		t.Fatal(err)
	}
	if first != second || len(first) != 64 {
		t.Fatalf("derived jwtSecret = %q", first)
	}
}

func TestLoadDerivesJWTSecretFromCredentialKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte(`secrets:
  credentialEncryptionKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
bootstrapAdmin:
  username: "admin"
  password: "password123"
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := deriveJWTSecret("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Secrets.JWTSecret != expected {
		t.Fatalf("jwtSecret = %q, want %q", cfg.Secrets.JWTSecret, expected)
	}
}

func TestLoadCredentialKeyFromEnvOverridesYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte(`secrets:
  credentialEncryptionKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
bootstrapAdmin:
  password: "password123"
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvCredentialEncryptionKey, "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := deriveJWTSecret("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Secrets.CredentialEncryptionKey != "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=" || cfg.Secrets.JWTSecret != expected {
		t.Fatalf("secrets = %#v", cfg.Secrets)
	}
}

func TestLoadRequiresCredentialEncryptionKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("bootstrapAdmin:\n  password: password123\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("missing credential encryption key was accepted")
	}
}