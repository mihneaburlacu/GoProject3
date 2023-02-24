package main

import (
	"errors"
	"net/http"
	"testing"
)

func TestHandlerMain(t *testing.T) {
	type errorCases struct {
		name      string
		serveFunc func(addr string, handler http.Handler) error
		wantError bool
	}
	for _, scenario := range []errorCases{
		{
			name: "error while listen and serve",
			serveFunc: func(addr string, handler http.Handler) error {
				return errors.New("error while listen and serve")
			},
			wantError: true,
		},
		{
			name: "no error",
			serveFunc: func(addr string, handler http.Handler) error {
				return nil
			},
			wantError: false,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			var got bool

			gotErr := handlerMain(scenario.serveFunc)
			if gotErr != nil {
				got = true
			} else {
				got = false
			}

			if got != scenario.wantError {
				t.Errorf("wanted: %#v, got: %#v", scenario.wantError, gotErr)
			}
		})
	}
}
