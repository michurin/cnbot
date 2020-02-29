package cfgreader

import (
	"path"
)

func pathToScript(cfgPath, scriptPath string) string {
	if scriptPath == "" {
		return scriptPath
	}
	if path.IsAbs(scriptPath) {
		return path.Clean(scriptPath)
	}
	return path.Clean(path.Join(path.Dir(cfgPath), scriptPath))
}
