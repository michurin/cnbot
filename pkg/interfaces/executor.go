package interfaces

import "context"

type Executor interface {
	Run(ctx context.Context, env []string, args []string) ([]byte, error)
}
