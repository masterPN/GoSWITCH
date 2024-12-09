package wholesaledatabase

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
			wantError: false,
		},
		{
			name:      "empty password",
			dbname:    "test_db",
			password:  "",
			username:  "test_user",
			port:      "1433",
			host:      "localhost",
			wantError: false,
		},
		{
			name:      "empty username",
			dbname:    "test_db",
			password:  "test_password",
			username:  "",
			port:      "1433",
			host:      "localhost",
			wantError: false,
		},
		{
			name:      "empty port",
			dbname:    "test_db",
			password:  "test_password",
			username:  "test_user",
			port:      "",
			host:      "localhost",
			wantError: false,
		},
		{
			name:      "empty host",
			dbname:    "test_db",
			password:  "test_password",
			username:  "test_user",
			port:      "1433",
			host:      "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("WHOLESALE_DB_DATABASE", tt.dbname)
			os.Setenv("WHOLESALE_DB_PASSWORD", tt.password)
			os.Setenv("WHOLESALE_DB_USERNAME", tt.username)
			os.Setenv("WHOLESALE_DB_PORT", tt.port)
			os.Setenv("WHOLESALE_DB_HOST", tt.host)

			got := New()
			assert.NotNil(t, got)
		})
	}
}
