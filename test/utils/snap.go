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

func SnapInstallFromStore(t *testing.T, name, channel string) {
	exec(t, fmt.Sprintf(
		"sudo snap install %s --channel=%s",
		name,
		channel,
	), true)
}

func SnapInstallFromFile(t *testing.T, path string) {
	exec(t, fmt.Sprintf(
		"sudo snap install --dangerous %s",
		path,
	), true)
}

func SnapRemove(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap remove --purge %s",
			name,
		), true)
	}
}

func SnapBuild(t *testing.T, workDir string) {
	exec(t, fmt.Sprintf(
		"cd %s && snapcraft",
		workDir,
	), true)
}

func SnapConnect(t *testing.T, plug, slot string) {
	exec(t, fmt.Sprintf(
		"sudo snap connect %s %s",
		plug, slot,
	), true)
}

func SnapDisconnect(t *testing.T, plug, slot string) {
	exec(t, fmt.Sprintf(
		"sudo snap disconnect %s %s",
		plug, slot,
	), true)
}

func SnapVersion(t *testing.T, name string) string {
	out, _ := exec(t, fmt.Sprintf(
		"snap info %s | grep installed | awk '{print $2}'",
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
	const filename = "snap.log" // used in action.yml
	exec(t, fmt.Sprintf("(%s) > %s",
		snapJournalCommand(start, name),
		filename), true)

	wd, _ := os.Getwd()
	fmt.Printf("Wrote snap logs to %s/%s\n", wd, filename)
}

func SnapLogs(t *testing.T, start time.Time, name string) string {
	logs, _ := exec(t, snapJournalCommand(start, name), false)
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
			"sudo snap start %s",
			name,
		), true)
	}
}

func SnapStop(t *testing.T, names ...string) {
	for _, name := range names {
		exec(t, fmt.Sprintf(
			"sudo snap stop %s",
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
