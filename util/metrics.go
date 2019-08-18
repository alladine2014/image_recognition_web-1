package util

import (
	"context"
	"github.com/cgCodeLife/logs"
	"time"
)

const (
	throuput = ".throughput"
	latency  = ".latency"
	err      = ".error"
	speed    = ".speed(mb/s)"
)

// EmitThroughput Method
func EmitThroughput(mkey string, tagkv map[string]string) {
	s := padding(tagkv)
	defer restore(tagkv, s)
}

// EmitLatency Method, this should be called with defer when entering a method
func EmitLatency(ctx context.Context, mkey string, t0 time.Time, tagkv map[string]string) {
	s := padding(tagkv)
	defer restore(tagkv, s)

	cost := time.Since(t0)
	logs.CtxInfo(ctx, "EmitTimer=%s cost=%v", mkey, cost)
}

// EmitError method
func EmitError(mkey string, tagkv map[string]string) {
	s := padding(tagkv)
	defer restore(tagkv, s)
}

// padding empty value with "-"
func padding(tagkv map[string]string) []string {
	s := []string{}
	for k, v := range tagkv {
		if v == "" {
			tagkv[k] = "-"
			s = append(s, k)
		}
	}
	return s
}

// restore value with empty string for given keys
func restore(tagkv map[string]string, keys []string) {
	for _, k := range keys {
		tagkv[k] = ""
	}
}
