package middleware

import (
	"net/http"
	"path/filepath"

	"github.com/unrolled/secure"
)

const (
	allowOriginHeader = "Access-Control-Allow-Origin"
)

func Secure(isProd bool) func(http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect:          true,
		STSSeconds:           315360000,
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		IsDevelopment:        !isProd,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldAllowOrigin(r) {
				w.Header().Add(allowOriginHeader, "*")
			}

			err := secureMiddleware.Process(w, r)
			if err != nil {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func shouldAllowOrigin(req *http.Request) bool {
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 {
		return false
	}

	switch extension {
	case ".eot", ".ttf", ".otf", ".woff", ".woff2":
		return true
	default:
		return false
	}
}
