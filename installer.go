package torgo

import (
	"os/exec"
	"runtime"
)

func CheckInstalledTor() bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("where", "tor")
	case "darwin", "linux":
		cmd = exec.Command("which", "tor")
	default:
		return false
	}

	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
