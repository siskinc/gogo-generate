package common

import (
	"fmt"
	"go/ast"
	"testing"
)

func TestParseStruct(t *testing.T) {
	type args struct {
		filename string
		src      []byte
	}
	tests := []struct {
		name                  string
		args                  args
		wantStructFieldsMap   map[string][]*ast.Field
		wantStructDocumentMap map[string][]string
		wantErr               bool
	}{
		{
			"test1",
			args{
				filename: "test/test1.go",
				src:      nil,
			},
			nil,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStructFieldsMap, gotStructDocumentMap, err := ParseStruct(tt.args.filename, tt.args.src)
			fmt.Println("gotStructFieldsMap", gotStructFieldsMap)
			fmt.Println("gotStructDocumentMap", gotStructDocumentMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
