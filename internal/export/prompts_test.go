package export_test

import (
	"testing"

	"github.com/vigia-io/vigia-io-exporter/internal/export"
)

func TestPasswordFromEnv(t *testing.T) {
	t.Setenv("VIGIA_TEST_PASSWORD", "secret")

	pw, err := export.PasswordFromEnv("VIGIA_TEST_PASSWORD")
	if err != nil || pw != "secret" {
		t.Fatalf("got %q err=%v", pw, err)
	}

	if _, err := export.PasswordFromEnv("VIGIA_MISSING"); err == nil {
		t.Fatal("expected error for missing env")
	}
}
