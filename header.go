package lowerkeys

import (
	"fmt"
	"net/http"
	"strings"
)

// Header is a http.Header with HTTP 2 lower-case key-compliant overrides
// for methods Add(), Del(), Get(), Set() and Values().
type Header http.Header

// New returns a new lowerkeys.Header.
func New() Header {
	return make(Header)
}

// From returns a new Header for the given http.Header
// with all key names in lower-case.
//
// If there are two or more keys which have the same lower-case name,
// this function will panic.
func From(header http.Header) Header {
	newKeyValues := make(Header)
	for key, values := range header {
		if values == nil {
			newKeyValues[key] = nil // I don't know what you are trying to accomplish, but I don't want to stop you, either.
		} else {
			newKeyValues[key] = append(make([]string, 0, len(values)), values...)
		}
	}
	return makeLowercase(newKeyValues)
}

// Using casts the given http.Header to a lowerkeys.Header,
// changes all header names to lower-case, and returns it.
//
// If there are two or more keys which have the same lower-case name,
// this function will panic.
func Using(header http.Header) Header {
	return makeLowercase(Header(header))
}

func makeLowercase(header map[string][]string) Header {
	for oldKey, values := range header {
		newKey := strings.ToLower(oldKey)
		if newKey == oldKey {
			continue
		}
		delete(header, oldKey)
		if _, ok := header[newKey]; ok {
			panic(fmt.Sprintf(`encountered two keys with identical lower-case representation: "%s" and "%s"`, newKey, oldKey))
		}
		header[newKey] = values
	}
	return header
}

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is made lower-case beforehand.
func (h Header) Add(key, value string) {
	key = strings.ToLower(key)
	h[key] = append(h[key], value)
}

// Del deletes the values associated with strings.ToLower(key).
func (h Header) Del(key string) {
	delete(h, strings.ToLower(key))
}

// Get gets the first value associated with strings.ToLower(key).
// If there are no values associated with the key, Get returns "".
func (h Header) Get(key string) string {
	key = strings.ToLower(key)
	if values, ok := h[key]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// Set sets the header entries associated with strings.ToLower(key)
// to the single element value. It replaces any existing values
// associated with strings.ToLower(key).
func (h Header) Set(key, value string) {
	h[strings.ToLower(key)] = []string{value}
}

// Values returns all values associated with strings.ToLower(key).
// The returned slice is not a copy.
func (h Header) Values(key string) []string {
	key = strings.ToLower(key)
	if values, ok := h[key]; ok {
		return values
	}
	return nil
}
