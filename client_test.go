package main

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig()
	t.Log(AccessToken)
}
