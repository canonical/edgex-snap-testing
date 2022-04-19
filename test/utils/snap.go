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
	Exec(t, fmt.Sprintf(
		"sudo snap install %s --channel=%s",
		name,
		channel,
	))
}

func SnapInstallFromFile(t *testing.T, path string) {
	Exec(t, fmt.Sprintf(
		"sudo snap install --dangerous %s",
		path,
	))
}

func SnapRemove(t *testing.T, names ...string) {
	for _, name := range names {
		Exec(t, fmt.Sprintf(
			"sudo snap remove --purge %s",
			name,
		))
	}
}

func SnapBuild(t *testing.T, workDir string) {
	Exec(t, fmt.Sprintf(
		"cd %s && snapcraft",
		workDir,
	))
}

func SnapConnect(t *testing.T, plug, slot string) {
	Exec(t, fmt.Sprintf(
		"sudo snap connect %s %s",
		plug, slot,
	))
}

func SnapDisconnect(t *testing.T, plug, slot string) {
	Exec(t, fmt.Sprintf(
		"sudo snap disconnect %s %s",
		plug, slot,
	))
}

func SnapVersion(t *testing.T, name string) string {
	out, _ := Exec(t, fmt.Sprintf(
		"snap info %s | grep installed | awk '{print $2}'",
		name,
	))
	return strings.TrimSpace(out)
}

// TODO: should the logs be fetched in each test?
// for that, need to use journalctl instead with --since
func SnapDumpLogs(t *testing.T, start time.Time, name string) {
	const filename = "snap.log" // used in action.yml
	Exec(t, fmt.Sprintf(
		"(sudo journalctl --since \"%s\" --no-pager | grep \"%s\"|| true) > %s",
		start.Format("2006-01-02 15:04:05"),
		name,
		filename))

	wd, _ := os.Getwd()
	fmt.Printf("Wrote snap logs to %s/%s\n", wd, filename)
}

func SnapLogsJournal(t *testing.T, start time.Time, name string) string {
	logs, _ := Exec(t, fmt.Sprintf(
		"sudo journalctl --since \"%s\" --no-pager | grep \"%s\"|| true\n\n",
		start.Format("2006-01-02 15:04:05"),
		name))

	return logs
}

func SnapSet(t *testing.T, name, key, value string) {
	Exec(t, fmt.Sprintf(
		"sudo snap set %s %s=%s",
		name,
		key,
		value,
	))
}

func SnapUnset(t *testing.T, name, key string) {
	Exec(t, fmt.Sprintf(
		"sudo snap unset %s %s",
		name,
		key,
	))
}

func SnapStart(t *testing.T, names ...string) {
	for _, name := range names {
		Exec(t, fmt.Sprintf(
			"sudo snap start %s",
			name,
		))
	}
}

func SnapStop(t *testing.T, names ...string) {
	for _, name := range names {
		Exec(t, fmt.Sprintf(
			"sudo snap stop %s",
			name,
		))
	}
}

func SnapRestart(t *testing.T, names ...string) {
	for _, name := range names {
		Exec(t, fmt.Sprintf(
			"sudo snap restart %s",
			name,
		))
	}
}
