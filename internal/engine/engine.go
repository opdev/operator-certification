package engine

import (
	"context"
	"fmt"

	"github.com/opdev/operator-certification/internal/image"
)

type engine struct {
	image        string
	dockerConfig string
}

type OperatorOption func(*engine)

func New(ctx context.Context, options ...OperatorOption) *engine {
	return &engine{}
}

func (engine) Execute(ctx context.Context, imageRef image.Reference) error {
	return fmt.Errorf("not implemented")
}

func WithImage(image string) OperatorOption {
	return func(e *engine) {
		e.image = image
	}
}

func WithDockerConfig(dockerConfig string) OperatorOption {
	return func(e *engine) {
		e.dockerConfig = dockerConfig
	}
}
