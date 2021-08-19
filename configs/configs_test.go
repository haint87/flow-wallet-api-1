package configs

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Run("Read environment variables from file", func(t *testing.T) {
		opts := &Options{EnvFilePath: ".testfile"}
		cfg, err := ParseConfig(opts)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.AdminAccountAddress != "admin-address" {
			t.Errorf(`expected "admin-address", got "%s"`, cfg.AdminAccountAddress)
		}
	})

	t.Run("", func(t *testing.T) {
		cfg, err := ParseConfig(nil)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.AdminAccountAddress != "admin-address" {
			t.Errorf(`expected "admin-address", got "%s"`, cfg.AdminAccountAddress)
		}
	})
}
