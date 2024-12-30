package http

import (
	"fmt"
	"strings"
	"sync"
)

type Header map[string][]string

func (h Header) String() string {
	var builder strings.Builder
	for key, value := range h {
		builder.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(value, ", ")))
	}
	return builder.String()
}

func (h Header) Add(key, value string) {
	key = CanonicalHeaderKey(key)

	h[key] = append(h[key], value)
}

func (h Header) Set(key, value string) {
	h[CanonicalHeaderKey(key)] = []string{value}
}

func (h Header) Get(key string) string {
	values := h.Values(CanonicalHeaderKey(key))

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (h Header) Values(key string) []string {
	return h[CanonicalHeaderKey(key)]
}

var canonicalHeader map[string]string

var canonicalHeaderOnce sync.Once

func initCanonicalHeader() {
	canonicalHeader = make(map[string]string)
	for _, v := range []string{
		"Accept",
		"Host",
		"User-Agent",
	} {
		canonicalHeader[strings.ToLower(v)] = v
	}
}

func CanonicalHeaderKey(key string) string {
	canonicalHeaderOnce.Do(initCanonicalHeader)

	canonicalKey := canonicalHeader[strings.ToLower(key)]

	if canonicalKey == "" {
		return key
	}

	return canonicalKey
}
