package config

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	type want struct {
		cfg Config
		err error
	}
	type test struct {
		name  string
		flags []string
		env   func()
		want  want
	}

	tests := []test{
		{
			name:  "successfull read config with flags",
			flags: []string{"test", "-host", "123.123.111.111", "-port", "1234", "-debug"},
			want: want{
				cfg: Config{
					Host:        "123.123.111.111",
					Port:        1234,
					DbDSN:       "postgres://user:password@localhost:5432/gt4?sslmode=disable",
					MigratePath: "migrations",
					Debug:       true,
				},
				err: nil,
			},
		},
		{
			name:  "successfull read config with envs",
			flags: []string{"test", "-debug"},
			env: func() {
				t.Setenv("DB_DSN", "postgres://user:password@171.191.1.11:5432/test-db")
				t.Setenv("MIGRATE_PATH", "db/migrate")
				t.Setenv("SRV_HOST", "171.191.15.11")
				t.Setenv("SRV_PORT", "7777")
			},
			want: want{
				cfg: Config{
					Host:        "171.191.15.11",
					Port:        7777,
					DbDSN:       "postgres://user:password@171.191.1.11:5432/test-db",
					MigratePath: "db/migrate",
					Debug:       true,
				},
				err: nil,
			},
		},
		{
			name:  "invalid host",
			flags: []string{"test", "-host", "999.191.1.11", "-port", "7777"},
			env: func() {
				t.Setenv("DB_DSN", "postgres://user:password@171.191.1.11:5432/test-db")
				t.Setenv("MIGRATE_PATH", "db/migrate")
			},
			want: want{
				cfg: Config{},
				err: ErrInvalidHost,
			},
		},
		{
			name:  "invalid port",
			flags: []string{"test", "-host", "999.191.1.11"},
			env: func() {
				t.Setenv("DB_DSN", "postgres://user:password@171.191.1.11:5432/test-db")
				t.Setenv("MIGRATE_PATH", "db/migrate")
				t.Setenv("SRV_PORT", "abcd")
			},
			want: want{
				cfg: Config{},
				err: strconv.ErrSyntax,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tc.flags
			if tc.env != nil {
				tc.env()
				defer os.Unsetenv("DB_DSN")
				defer os.Unsetenv("MIGRATE_PATH")
				defer os.Unsetenv("SRV_HOST")
				defer os.Unsetenv("SRV_PORT")
			}
			cfg, err := ReadConfig()
			assert.ErrorIs(t, err, tc.want.err)
			assert.Equal(t, tc.want.cfg, cfg)
		})
	}
}
