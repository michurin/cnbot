package interfaces

import "context"

type Executor interface {
	Run(ctx context.Context, script string, env []string, args []string) ([]byte, error)
}
