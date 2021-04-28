package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_findMostRecent(t *testing.T) {
	dir := t.TempDir()
	os.Create(filepath.Join(dir, "a"))
	os.Create(filepath.Join(dir, "b"))
	os.Create(filepath.Join(dir, "c"))
	os.Create(filepath.Join(dir, "last"))
	os.Create(filepath.Join(dir, "last.ext"))
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"simple",
			args{dir},
			"last",
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findMostRecent(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMostRecent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMostRecent() = %v, want %v", got, tt.want)
			}
		})
	}
}
