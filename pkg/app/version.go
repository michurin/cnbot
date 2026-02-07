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
	fmt.Println(modInfo(info.Main))
	fmt.Println(info.String())
}

func MainVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "no build info"
	}
	return modInfo(info.Main)
}

func modInfo(m debug.Module) string {
	s := m.Path + " " + m.Version + " " + m.Sum
	if m.Replace != nil {
		return s + " => " + modInfo(*m.Replace)
	}
	return s
}
