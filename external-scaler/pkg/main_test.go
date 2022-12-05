package main

import (
	"testing"
)

func TestExternalScaler(t *testing.T) {
	es := &ExternalScaler{}
	if es == nil {
		t.Fatalf(`ExternalScaler, want not nil`)
	}
}
