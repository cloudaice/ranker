package main

import (
	"testing"
)

func TestLoadBucket(t *testing.T) {
	data, err := LoadBucket("china")
	if err != nil {
		t.Errorf("LoadBucket error: %s\n", err)
		return
	}

	users := LoadUsers(data)
	for _, user := range users {
		t.Log(user.Score)
	}
}
