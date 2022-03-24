package utils

import (
	"fmt"
	"strings"
	"testing"
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

func SnapVersion(t *testing.T, name string) string {
	out, _ := Exec(t, fmt.Sprintf(
		"snap info %s | grep installed | awk '{print $2}'",
		name,
	))
	return strings.TrimSpace(out)
}

// TODO: should the logs be fetched in each test?
// for that, need to use journalctl instead with --since
func SnapLogs(t *testing.T, name string) string {
	stdout, _ := Exec(t, fmt.Sprintf(
		"sudo snap logs -n=all %s",
		name))
	return stdout
}
