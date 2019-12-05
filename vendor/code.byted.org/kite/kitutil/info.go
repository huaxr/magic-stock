package kitutil

import (
	"context"
	"strings"
)

const (
	PrefixPersist         = "RPC_PERSIST_"
	PrefixTransit         = "RPC_TRANSIT_"
	PrefixTransitUpstream = "RPC_TRANSIT_UNSTREAM_"
)

// Using empty string as key or value is not support.

// GetValue retrieves the value set into the context by given key.
func GetValue(ctx context.Context, k string) (string, bool) {
	if v, ok := ctxGetV(ctx, PrefixTransit+k); ok {
		return v, len(v) > 0
	}
	if v, ok := ctxGetV(ctx, PrefixTransitUpstream+k); ok {
		return v, len(v) > 0
	}
	return "", false
}

// GetAllValues retrieves all transient values
func GetAllValues(ctx context.Context) map[string]string {
	res := make(map[string]string)
	if m := ctxGetAll(ctx); m != nil {
		for k, v := range m {
			if len(v) == 0 {
				continue
			}
			if strings.HasPrefix(k, PrefixTransitUpstream) {
				res[k[len(PrefixTransitUpstream):]] = v
				continue
			}
			if strings.HasPrefix(k, PrefixTransit) {
				res[k[len(PrefixTransit):]] = v
			}
		}
	}
	return res
}

// WithValue sets the value into the context by given key.
// This value will be propagated to the next service/endpoint through a RPC call.
//
// Notice that it will not propagate any further beyond the next service/endpoint,
// Use WithPersistValue if you want to pass a key/value pair all the way.
func WithValue(ctx context.Context, k string, v string) context.Context {
	if len(k) == 0 || len(v) == 0 {
		return ctx
	}
	return ctxAddKV(ctx, PrefixTransit+k, v)
}

// DelValue deletes a key/value from current context.
// Since empty string value is not valid, we could just set the value to be empty.
func DelValue(ctx context.Context, k string) context.Context {
	if len(k) == 0 {
		return ctx
	}
	return ctxAddKV(ctx, PrefixTransit+k, "")
}

// GetPersistValue retrieves the persistent value set into the context by given key.
func GetPersistValue(ctx context.Context, k string) (string, bool) {
	if v, ok := ctxGetV(ctx, PrefixPersist+k); ok {
		return v, len(v) > 0
	}
	return "", false
}

// GetAllPersistValues retrieves all persistent value
func GetAllPersistValues(ctx context.Context) map[string]string {
	res := make(map[string]string)
	if m := ctxGetAll(ctx); m != nil {
		for k, v := range m {
			if len(v) == 0 {
				continue
			}
			if strings.HasPrefix(k, PrefixPersist) {
				res[k[len(PrefixPersist):]] = v
			}
		}
	}
	return res
}

// WithPersistValue sets the value info the context by given key.
// This value will be propagated to the services along the RPC call chain.
func WithPersistValue(ctx context.Context, k string, v string) context.Context {
	if len(k) == 0 || len(v) == 0 {
		return ctx
	}
	return ctxAddKV(ctx, PrefixPersist+k, v)
}

// DelPersistValue deletes a persistent key/value from current context.
// Since empty string value is not valid, we could just set the value to be empty.
func DelPersistValue(ctx context.Context, k string) context.Context {
	if len(k) == 0 {
		return ctx
	}
	return ctxAddKV(ctx, PrefixPersist+k, "")
}
