package console

import (
	"os"
)

// IsConsoleVersion 检测是否查看 `server` 版本号
func IsConsoleVersion() bool {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		return true
	}

	return false
}
