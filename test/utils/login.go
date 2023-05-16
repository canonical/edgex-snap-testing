package utils

import (
	"strings"
	"testing"
)

func LoginTestUser(t *testing.T) (idToken string) {
	// The script path relative to the testing suites
	const loginScriptPath = "../../scripts/login-test-user.sh"

	idToken, _, _ = Exec(t, loginScriptPath)
	t.Log("ID Token for 'example' user:", idToken)
	return strings.TrimSpace(idToken)
}
