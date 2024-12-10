package helpers

import (
	"os"
	"testing"
)

const (
	exampleDomainPhrase = "example.com"
)

func TestLoadConfiguration(t *testing.T) {
	tests := []struct {
		name                   string
		envVars                map[string]string
		expectedSIPPort        int
		expectedExternalDomain string
		expectedError          bool
	}{
		{
			name: "valid SIP_PORT and EXTERNAL_DOMAIN",
			envVars: map[string]string{
				"SIP_PORT":        "5060",
				"EXTERNAL_DOMAIN": exampleDomainPhrase,
			},
			expectedSIPPort:        5060,
			expectedExternalDomain: exampleDomainPhrase,
			expectedError:          false,
		},
		{
			name: "invalid SIP_PORT",
			envVars: map[string]string{
				"SIP_PORT":        "not an integer",
				"EXTERNAL_DOMAIN": exampleDomainPhrase,
			},
			expectedSIPPort:        0,
			expectedExternalDomain: exampleDomainPhrase,
			expectedError:          true,
		},
		{
			name: "missing SIP_PORT",
			envVars: map[string]string{
				"EXTERNAL_DOMAIN": exampleDomainPhrase,
			},
			expectedSIPPort:        0,
			expectedExternalDomain: exampleDomainPhrase,
			expectedError:          true,
		},
		{
			name: "missing EXTERNAL_DOMAIN",
			envVars: map[string]string{
				"SIP_PORT": "5060",
			},
			expectedSIPPort:        5060,
			expectedExternalDomain: "",
			expectedError:          false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set environment variables
			os.Clearenv()
			for key, value := range test.envVars {
				os.Setenv(key, value)
			}

			// Call LoadConfiguration
			sipPort, externalDomain := LoadConfiguration()

			// Check results
			if sipPort != test.expectedSIPPort {
				t.Errorf("expected SIP port %d, got %d", test.expectedSIPPort, sipPort)
			}
			if externalDomain != test.expectedExternalDomain {
				t.Errorf("expected external domain %q, got %q", test.expectedExternalDomain, externalDomain)
			}
			if test.expectedError && sipPort != 0 {
				t.Errorf("expected error, but got SIP port %d", sipPort)
			}
		})
	}
}
