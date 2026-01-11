package middlewares

import (
	"context"
	"net/http"
	"runtime/pprof"

	"github.com/gorilla/mux"
	"github.com/grafana/pyroscope-go"
	"golang.org/x/exp/slices"
)

type PyroProfiler struct {
	staticLabels []string
}

func NewPyroProfiler(staticLabels map[string]string) *PyroProfiler {
	return &PyroProfiler{
		staticLabels: flattenLabels(staticLabels),
	}
}

func flattenLabels(m map[string]string) []string {
	out := make([]string, 0, len(m)*2)
	for k, v := range m {
		out = append(out, k, v)
	}

	return out
}

func (pp *PyroProfiler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		labelSet := slices.Clone(pp.staticLabels)
		if route := mux.CurrentRoute(r); route != nil {
			if routeName := route.GetName(); routeName != "" {
				labelSet = append(labelSet, "endpoint", routeName)
			}
		}

		// Run the handler inside the profiling label scope
		pyroscope.TagWrapper(ctx, pprof.Labels(labelSet...), func(ctx context.Context) {
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
}
