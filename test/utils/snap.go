package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// func SnapInstall(t *testing.T, name string) {
// 	if strings.HasSuffix(name, ".snap") {
// 		SnapInstallFromFile(nil, name)
// 	} else {
// 		SnapInstallFromStore(nil, name, ServiceChannel)
// 	}
// }

func SnapInstallFromStore(t *testing.T, name, channel string) error {
	_, stderr, err := exec(t, fmt.Sprintf(
		"sudo snap install %s --channel=%s",
		name,
		channel,
	), true)
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapInstallFromFile(t *testing.T, path string) error {
	_, stderr, err := exec(t, fmt.Sprintf(
		"sudo snap install --dangerous %s",
		path,
	), true)
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapRemove(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap remove --purge %s",
			name,
		), true)
	}
}

func SnapBuild(t *testing.T, workDir string) error {
	_, stderr, err := exec(t, fmt.Sprintf(
		"cd %s && snapcraft",
		workDir,
	), true)
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapConnect(t *testing.T, plug, slot string) error {
	_, stderr, err := exec(t, fmt.Sprintf(
		"sudo snap connect %s %s",
		plug, slot,
	), true)
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr)
	}
	return nil
}

func SnapConnectSecretstoreToken(t *testing.T, snap string) error {
	return SnapConnect(t,
		"edgexfoundry:edgex-secretstore-token",
		snap+":edgex-secretstore-token")
}

func SnapDisconnect(t *testing.T, plug, slot string) {
	exec(t, fmt.Sprintf(
		"sudo snap disconnect %s %s",
		plug, slot,
	), true)
}

func SnapVersion(t *testing.T, name string) string {
	out, _, _ := exec(t, fmt.Sprintf(
		"snap info %s | grep installed | awk '{print $2}'",
		name,
	), true)
	return strings.TrimSpace(out)
}

func SnapRevision(t *testing.T, name string) string {
	out, _, _ := exec(t, fmt.Sprintf(
		"snap list %s | awk 'NR==2 {print $3}'",
		name,
	), true)
	return strings.TrimSpace(out)
}

func snapJournalCommand(start time.Time, name string) string {
	// The command should not return error even if nothing is grepped, hence the "|| true"
	return fmt.Sprintf("sudo journalctl --since \"%s\" --no-pager | grep \"%s\"|| true",
		start.Format("2006-01-02 15:04:05"),
		name)
}

func SnapDumpLogs(t *testing.T, start time.Time, name string) {
	filename := name + ".log" // used in action.yml
	exec(t, fmt.Sprintf("(%s) > %s",
		snapJournalCommand(start, name),
		filename), true)

	wd, _ := os.Getwd()
	fmt.Printf("Wrote snap logs to %s/%s\n", wd, filename)
}

func SnapLogs(t *testing.T, start time.Time, name string) string {
	logs, _, _ := exec(t, snapJournalCommand(start, name), false)
	return logs
}

func SnapSet(t *testing.T, name, key, value string) {
	exec(t, fmt.Sprintf(
		"sudo snap set %s %s='%s'",
		name,
		key,
		value,
	), true)
}

func SnapUnset(t *testing.T, name, key string) {
	exec(t, fmt.Sprintf(
		"sudo snap unset %s %s",
		name,
		key,
	), true)
}

func SnapStart(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap start --enable %s",
			name,
		), true)
	}
}

func SnapStop(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap stop --disable %s",
			name,
		), true)
	}
}

func SnapRestart(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap restart %s",
			name,
		), true)
	}
}

func SnapRefresh(t *testing.T, name, channel string) {
	exec(t, fmt.Sprintf(
		"sudo snap refresh %s --channel=%s --amend",
		name,
		channel,
	), true)
}

func SnapServicesEnabled(t *testing.T, name string) bool {
	out, _, _ := exec(t, fmt.Sprintf(
		"snap services %s | awk 'FNR == 2 {print $2}'",
		name,
	), true)
	return strings.TrimSpace(out) == "enabled"
}

func SnapServicesActive(t *testing.T, name string) bool {
	out, _, _ := exec(t, fmt.Sprintf(
		"snap services %s | awk 'FNR == 2 {print $3}'",
		name,
	), true)
	return strings.TrimSpace(out) == "active"
}

func LocalSnap() bool {
	return LocalSnapPath != ""
}
