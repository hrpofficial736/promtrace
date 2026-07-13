package certmanager

import (
	"testing"

	"github.com/hrpofficial736/promtrace/internal/config"
)

func TestGenerateAndLoadCA(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Dir:    tmpDir,
		CACert: tmpDir + "/ca.crt",
		CAKey:  tmpDir + "/ca.key",
	}

	cm, err := NewCertManager(cfg)
	if err != nil {
		t.Fatalf("NewCertManager failed: %v", err)
	}

	if err := cm.GenerateRootCACertificate(); err != nil {
		t.Fatalf("GenerateRootCACertificate failed: %v", err)
	}

	cm2, err := NewCertManager(cfg)
	if err != nil {
		t.Fatalf("NewCertManager (2nd) failed: %v", err)
	}

	if err := cm2.LoadCA(); err != nil {
		t.Fatalf("LoadCA failed: %v", err)
	}
}

func TestGetOrCreateHostCertificate(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Dir:    tmpDir,
		CACert: tmpDir + "/ca.crt",
		CAKey:  tmpDir + "/ca.key",
	}

	cm, _ := NewCertManager(cfg)
	cm.GenerateRootCACertificate()

	cert, err := cm.GetOrCreateHostCertificate("example.com")
	if err != nil {
		t.Fatalf("GetOrCreateHostCertificate failed: %v", err)
	}

	if cert == nil {
		t.Fatal("expected non-nil certificate")
	}

	cert2, _ := cm.GetOrCreateHostCertificate("example.com")
	if cert != cert2 {
		t.Error("expected cached certificate, got a new one")
	}

	cert3, _ := cm.GetOrCreateHostCertificate("other.com")
	if cert3 == cert {
		t.Error("expected a different certificate for the different host")
	}
}
