package main_test

import (
	"testing"

	api "github.com/sakthipriyan/diet-display/diet-display-api"
)

func TestIsDatabaseCreated(t *testing.T) {
	actual := api.IsDatabaseCreated("database.db")
	if actual != true {
		t.Errorf("expected 'false', got '%t'", actual)
	}
}
