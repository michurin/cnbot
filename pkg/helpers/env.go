package helpers

func Env(pairs ...string) []string {
	if len(pairs)%2 != 0 {
		panic(len(pairs))
	}
	env := []string(nil)
	for i := 0; i < len(pairs); i += 2 {
		env = append(env, pairs[i]+"="+pairs[i+1])
	}
	return env
}
