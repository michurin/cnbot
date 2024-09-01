package xjson

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func jsonToEnv(pfx string, x any) ([]string, error) {
	switch e := x.(type) {
	case nil:
		return nil, nil
	case bool:
		return []string{fmt.Sprintf("%s=%t", pfx, e)}, nil
	case float64:
		return []string{fmt.Sprintf("%s=%.20g", pfx, e)}, nil
	case string:
		return []string{fmt.Sprintf("%s=%s", pfx, e)}, nil
	case []any:
		l := len(e)
		if l == 0 {
			return nil, nil
		}
		res := make([]string, 0, l+1) // we cannot predict final length, however it's not less than l+1
		lst := make([]string, l)
		for i, v := range e {
			p := pfx + "_" + strconv.Itoa(i)
			lst[i] = p
			t, err := jsonToEnv(p, v)
			if err != nil {
				return nil, err
			}
			res = append(res, t...)
		}
		res = append(res, fmt.Sprintf("%s=%s", pfx, strings.Join(lst, " ")))
		return res, nil
	case map[string]any:
		res := make([]string, 0, len(e)) // the same case as above
		for k, v := range e {
			t, err := jsonToEnv(pfx+"_"+k, v)
			if err != nil {
				return nil, err
			}
			res = append(res, t...)
		}
		return res, nil
	}
	return nil, fmt.Errorf("invalid type [pfx=%s]: %T", pfx, x)
}

func JSONToEnv(x any) ([]string, error) {
	a, err := jsonToEnv("tg", x)
	if err != nil {
		return nil, err
	}
	sort.Strings(a) // to be reproducible
	return a, nil
}

func Slice(x any, k ...string) ([]any, error) {
	v, ok, err := Any(x, k...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("key not found: %#v", k)
	}
	b, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid type of value: %T", v)
	}
	return b, nil
}

func String(x any, k ...string) (string, error) {
	v, ok, err := Any(x, k...)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("key not found: %#v", k)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid type of value: %T", v)
	}
	return s, nil
}

func Int(x any, k ...string) (int64, error) {
	v, err := float(x, k...)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func float(x any, k ...string) (float64, error) {
	v, ok, err := Any(x, k...)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("key not found: %#v", k)
	}
	b, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid type of value: %T", v)
	}
	return b, nil
}

func Bool(x any, k ...string) (bool, error) {
	v, ok, err := Any(x, k...)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("key not found: %#v", k)
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("invalid type of value: %T", v)
	}
	return b, nil
}

func Any(x any, k ...string) (any, bool, error) {
	if len(k) == 0 {
		panic("no key") // invalid function usage; panic is reasonable
	}
	kv, ok := x.(map[string]any)
	if !ok {
		return nil, false, fmt.Errorf("invalid type of x: %T", x)
	}
	v, ok := kv[k[0]]
	if !ok {
		return nil, false, nil
	}
	if len(k) == 1 {
		return v, true, nil
	}
	return Any(v, k[1:]...)
}
