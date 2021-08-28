package api

import (
	"Go_news/pkg/dbnews"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		db *dbnews.DB
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
