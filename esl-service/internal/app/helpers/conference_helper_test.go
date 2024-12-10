package helpers

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/0x19/goesl"
	"github.com/go-playground/assert/v2"
)

const (
	exampleDomain  = "example.com"
	operatorPrefix = "test prefix"
	jobCommand     = "Job-Command"
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
				"EXTERNAL_DOMAIN": exampleDomain,
			},
			expectedSIPPort:        5060,
			expectedExternalDomain: exampleDomain,
			expectedError:          false,
		},
		{
			name: "invalid SIP_PORT",
			envVars: map[string]string{
				"SIP_PORT":        "not an integer",
				"EXTERNAL_DOMAIN": exampleDomain,
			},
			expectedSIPPort:        0,
			expectedExternalDomain: exampleDomain,
			expectedError:          true,
		},
		{
			name: "missing SIP_PORT",
			envVars: map[string]string{
				"EXTERNAL_DOMAIN": exampleDomain,
			},
			expectedSIPPort:        0,
			expectedExternalDomain: exampleDomain,
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

func TestHandleReadError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantLog string
	}{
		{
			name: "EOF error",
			err:  errors.New("EOF"),
		},
		{
			name: "unexpected end of JSON input error",
			err:  errors.New("unexpected end of JSON input"),
		},
		{
			name:    "other error",
			err:     errors.New("other error"),
			wantLog: "Error reading Freeswitch message: other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			logBuf := &bytes.Buffer{}
			log.SetOutput(logBuf)
			defer log.SetOutput(os.Stderr)

			handleReadError(tt.err)

			// Check log output
			logOutput, err := io.ReadAll(logBuf)
			if err != nil {
				t.Fatal(err)
			}

			if (tt.wantLog == "" && len(logOutput) > 0) || (tt.wantLog != "" && !strings.Contains(string(logOutput), tt.wantLog)) {
				t.Errorf("unexpected log output: got '%s', expected to contain '%s'", string(logOutput), tt.wantLog)
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	tests := []struct {
		name           string
		msg            *goesl.Message
		operatorPrefix string
		destination    string
		expected       bool
	}{
		{
			name: "add-member action and matching caller destination",
			msg: &goesl.Message{
				Headers: map[string]string{
					"Action":                "add-member",
					answerStateHeader:       "early",
					callerDestinationHeader: "operatorPrefixdestination",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       true,
		},
		{
			name: "add-member action and non-matching caller destination",
			msg: &goesl.Message{
				Headers: map[string]string{
					"Action":                "add-member",
					answerStateHeader:       "early",
					callerDestinationHeader: "wrongPrefixdestination",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "show Job-Command and matching conference body",
			msg: &goesl.Message{
				Headers: map[string]string{
					jobCommand:               "show",
					"Job-Command-Arg":        "channels",
					"Event-Calling-Function": "bgapi_exec",
					"body":                   "operatorPrefixdestination,conference",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       true,
		},
		{
			name: "show Job-Command and non-matching conference body",
			msg: &goesl.Message{
				Headers: map[string]string{
					jobCommand:               "show",
					"Job-Command-Arg":        "channels",
					"Event-Calling-Function": "bgapi_exec",
					"body":                   "wrongPrefixdestination,conference",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "non-matching action and Job-Command",
			msg: &goesl.Message{
				Headers: map[string]string{
					"Action":   "wrong-action",
					jobCommand: "wrong-command",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "empty message headers",
			msg: &goesl.Message{
				Headers: map[string]string{},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isConnected(test.msg, test.operatorPrefix, test.destination)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestIsCalleeUnavailable(t *testing.T) {
	tests := []struct {
		name           string
		msg            *goesl.Message
		operatorPrefix string
		destination    string
		expected       bool
	}{
		{
			name: "valid callee unavailable message",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "17",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       true,
		},
		{
			name: "invalid answer state",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "ringing",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "17",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "mismatched caller destination",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "wrongPrefixdestination",
					"variable_hangup_cause_q850": "17",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "missing hangup cause code",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:       "hangup",
					callerDestinationHeader: "operatorPrefixdestination",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "hangup cause code not in sipCalleeUnavailableCode",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "123",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isCalleeUnavailable(tt.msg, tt.operatorPrefix, tt.destination)
			if actual != tt.expected {
				t.Errorf("isCalleeUnavailable() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestIsOperatorUnavailable(t *testing.T) {
	tests := []struct {
		name           string
		msg            *goesl.Message
		operatorPrefix string
		destination    string
		expected       bool
	}{
		{
			name: "valid hangup cause code",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "1",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       true,
		},
		{
			name: "invalid hangup cause code",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "2",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "different answer state",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "ringing",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "1",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name: "different caller destination",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "differentPrefix+destination",
					"variable_hangup_cause_q850": "1",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
		{
			name:     "nil message",
			msg:      &goesl.Message{},
			expected: false,
		},
		{
			name: "empty hangup cause code",
			msg: &goesl.Message{
				Headers: map[string]string{
					answerStateHeader:            "hangup",
					callerDestinationHeader:      "operatorPrefixdestination",
					"variable_hangup_cause_q850": "",
				},
			},
			operatorPrefix: "operatorPrefix",
			destination:    "destination",
			expected:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isOperatorUnavailable(test.msg, test.operatorPrefix, test.destination)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestLogOperatorIssue(t *testing.T) {
	tests := []struct {
		name           string
		msg            *goesl.Message
		operatorPrefix string
		wantOutput     string
	}{
		{
			name: "valid message and operator prefix",
			msg: &goesl.Message{
				Headers: map[string]string{
					"variable_hangup_cause_q850":         "123",
					"variable_sip_invite_failure_phrase": "test reason",
				},
			},
			operatorPrefix: operatorPrefix,
			wantOutput:     "test prefix has a problem, please contact operator test prefix.\ncode - 123, reason - test reason\n",
		},
		{
			name:           "nil message",
			msg:            &goesl.Message{},
			operatorPrefix: operatorPrefix,
			wantOutput:     "test prefix has a problem, please contact operator test prefix.\ncode - , reason - \n",
		},
		{
			name: "empty operator prefix",
			msg: &goesl.Message{
				Headers: map[string]string{
					"variable_hangup_cause_q850":         "123",
					"variable_sip_invite_failure_phrase": "test reason",
				},
			},
			operatorPrefix: "",
			wantOutput:     " has a problem, please contact operator .\ncode - 123, reason - test reason\n",
		},
		{
			name: "missing headers in message",
			msg: &goesl.Message{
				Headers: map[string]string{},
			},
			operatorPrefix: operatorPrefix,
			wantOutput:     "test prefix has a problem, please contact operator test prefix.\ncode - , reason - \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			logOperatorIssue(tt.msg, tt.operatorPrefix)
			log.SetOutput(os.Stdout)
			gotOutput := buf.String()
			// Use a regular expression to remove the timestamp
			re := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `)
			gotOutput = re.ReplaceAllString(gotOutput, "")
			assert.Equal(t, tt.wantOutput, gotOutput)
		})
	}
}
