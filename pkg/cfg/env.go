package cfg

import (
	"fmt"
	"os"
)

func prepareEnv(toPass []string, toForce []string) []string {
	res := []string(nil)
	for _, k := range toPass {
		v, ok := os.LookupEnv(k)
		if ok {
			res = append(res, fmt.Sprintf("%s=%s", k, v))
		}
	}
	// TODO check toForce format
	return append(res, toForce...)
}
