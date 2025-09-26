package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"slices"

	"github.com/rs/zerolog"
)

func HTTPRequestBodyLoggerWithContext(ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(ctx)

			if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, r.Method) {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					logger.Info().
						Str("http.host", r.URL.Host).
						Str("http.target", r.URL.Path).
						Str("X-USER header", r.Header.Get("X-User")).
						Str("body", string(bodyBytes)).
						Str("HTTPClientIPKey", r.RemoteAddr).
						Int64("HTTPRequestContentLengthKey", r.ContentLength).
						Str("HTTPTargetKey", r.URL.RequestURI()).
						Str("ServerAddressKey", r.Host).
						Str("URLPathKey", r.URL.Path).
						Str("URLQueryKey", r.URL.RawQuery).
						Str("HTTPRefererKey", r.Referer()).
						Str("HTTPRequestIDKey", r.Header.Get("X-Request-ID")).
						Str("UserAgentOriginalKey", r.UserAgent()).
						Str("MethodKey", r.Method).
						Msg("request body")
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				} else {
					logger.Error().Err(err).Msg("failed to read request body")
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
