package app

import (
	"fmt"
	"runtime/debug"
)

func ShowVersionInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("No build info")
		return
	}
	fmt.Println(info.Main.Version)
	fmt.Println(info.String())
}

func MainVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "no version"
	}
	return info.Main.Version
}
