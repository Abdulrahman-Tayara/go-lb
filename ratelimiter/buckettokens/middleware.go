package buckettokens

import (
	"net"
	"net/http"
	"strings"
)

type Options struct {
	KeyFunc      func(req *http.Request) string
	ErrorHandler func(output *LimitOutput, w http.ResponseWriter, r *http.Request)
}

func ipFuncKey(req *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr))
	if err != nil {
		return ""
	}
	if netIP := net.ParseIP(ip); netIP != nil {
		return netIP.String()
	}
	return ""
}

func defaultErrorHandler(output *LimitOutput, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write([]byte("Too many requests"))
}

var DefaultOptions = &Options{
	KeyFunc:      ipFuncKey,
	ErrorHandler: defaultErrorHandler,
}

func LimiterMiddleware(next http.Handler, cfg *Configs, opts ...*Options) http.Handler {
	limiter := NewLimiter(cfg)
	opt := DefaultOptions

	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		doNext := func() {
			next.ServeHTTP(writer, request)
		}

		if key := opt.KeyFunc(request); key == "" {
			doNext()
			return
		} else if output := limiter.Limit(key); output.Allowed {
			doNext()
		} else {
			opt.ErrorHandler(output, writer, request)
		}
	})
}
