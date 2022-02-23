package utils

import (
	"fmt"
	"testing"
)

func SnapInstall(t *testing.T, name, channel string) {
	Exec(t, fmt.Sprintf(
		"sudo snap install %s --channel=%s",
		name,
		channel,
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

func SnapInstallLocal(t *testing.T, path string) {
	Exec(t, fmt.Sprintf(
		"sudo snap install --dangerous %s",
		path,
	))
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
