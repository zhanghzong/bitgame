package console

import "os"

// IsConsoleVersion 是否控制台查看版本号
func IsConsoleVersion() bool {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		return true
	}

	return false
}
