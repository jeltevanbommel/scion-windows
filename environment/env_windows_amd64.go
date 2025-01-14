package environment

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed windows64
var embeddedWindowsFiles embed.FS

func init() {

	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		fmt.Println("The Program Files directory could not be found.")
		os.Exit(1)
	}

	basePath := filepath.Join(programFiles, "scion-host")
	configPath := filepath.Join(basePath, "windows64")

	EndhostEnv = &EndhostEnvironment{
		Windows:              true,
		ConfigPath:           configPath,
		BasePath:             basePath,
		DispatcherBinaryPath: filepath.Join(configPath, "dispatcher.exe"),
		DispatcherConfigPath: configPath,
		DaemonBinaryPath:     filepath.Join(configPath, "daemon.exe"),
		DaemonConfigPath:     configPath,
		EmbeddedFiles:        embeddedWindowsFiles,
		EmbeddedFolder:       "windows64",
	}
}
