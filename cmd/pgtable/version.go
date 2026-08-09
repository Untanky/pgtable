package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

var AppVersion = ""

func displayVersions() {
	version, goVersion := getVersion()

	fmt.Fprintf(os.Stdin, "Version: %s\nGo Version: %s\n", version, goVersion)
}

func getVersion() (string, string) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown", "unknown"
	}

	if AppVersion != "" {
		return AppVersion, info.GoVersion
	}

	if info.Main.Version != "" {
		return info.Main.Version, info.GoVersion
	}

	return "development", info.GoVersion
}
