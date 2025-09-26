package build

import (
	"context"

	"github.com/rs/zerolog"
)

func (b *Builder) Shutdown(ctx context.Context) {
	zerolog.Ctx(ctx).Info().Msgf("got os signal. application will be stopped")
	b.shutdown.do(ctx)
}

type shutdownFn func(context.Context) error

type shutdown struct {
	fn []shutdownFn
}

func (s *shutdown) add(fn shutdownFn) {
	s.fn = append(s.fn, fn)
}

func (s *shutdown) do(ctx context.Context) {
	for i := len(s.fn) - 1; i >= 0; i-- {
		if err := s.fn[i](ctx); err != nil {
			zerolog.Ctx(ctx).Err(err).Send()
		}
	}
}
