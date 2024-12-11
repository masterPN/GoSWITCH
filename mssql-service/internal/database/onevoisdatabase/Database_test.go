package onevoisdatabase

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		dbname    string
		password  string
		username  string
		port      string
		host      string
		wantError bool
	}{
		{
			name:      "valid input",
			dbname:    "test_db",
			password:  "test_password",
			username:  "test_user",
			port:      "1433",
			host:      "localhost",
			wantError: false,
		},
		{
			name:      "empty dbname",
			dbname:    "",
			password:  "test_password",
			username:  "test_user",
			port:      "1433",
			host:      "localhost",
			wantError: true,
		},
		{
			name:      "empty password",
			dbname:    "test_db",
			password:  "",
			username:  "test_user",
			port:      "1433",
			host:      "localhost",
			wantError: true,
		},
		{
			name:      "empty username",
			dbname:    "test_db",
			password:  "test_password",
			username:  "",
			port:      "1433",
			host:      "localhost",
			wantError: true,
		},
		{
			name:      "empty port",
			dbname:    "test_db",
			password:  "test_password",
			username:  "test_user",
			port:      "",
			host:      "localhost",
			wantError: true,
		},
		{
			name:      "empty host",
			dbname:    "test_db",
			password:  "test_password",
			username:  "test_user",
			port:      "1433",
			host:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ONEVOIS_DB_DATABASE", tt.dbname)
			os.Setenv("ONEVOIS_DB_PASSWORD", tt.password)
			os.Setenv("ONEVOIS_DB_USERNAME", tt.username)
			os.Setenv("ONEVOIS_DB_PORT", tt.port)
			os.Setenv("ONEVOIS_DB_HOST", tt.host)

			got := New()
			assert.NotNil(t, got)
		})
	}
}
