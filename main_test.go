package main

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestHandlerMain(t *testing.T) {
	type errorCases struct {
		name      string
		serveFunc func(addr string, handler http.Handler) error
		wantError error
	}
	for _, scenario := range []errorCases{
		{
			name: "error while listen and serve",
			serveFunc: func(addr string, handler http.Handler) error {
				return errors.New("error while listen and serve")
			},
			wantError: errors.New("error while listen and serve"),
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			gotErr := handlerMain(scenario.serveFunc)

			if scenario.wantError != nil && gotErr != nil {
				diff := cmp.Diff(gotErr.Error(), scenario.wantError.Error())
				if diff != "" {
					t.Errorf("wanted: %#v, got: %#v", scenario.wantError, gotErr)
				}
			}
		})
	}
}
