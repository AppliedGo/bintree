package main

import (
	"reflect"
	"testing"
)

func TestTree_Delete(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		tree, want Tree
		args       args
		wantErr    bool
	}{
		{
			name: "Delete root in tree with three nodes",
			tree: Tree{
				Root: &Node{
					Value: "b",
					Data:  "b",
					Left: &Node{
						Value: "a",
						Data:  "a",
					},
					Right: &Node{
						Value: "c",
						Data:  "c",
					},
				},
			},
			want: Tree{
				Root: &Node{
					Value: "a",
					Data:  "a",
					Right: &Node{
						Value: "c",
						Data:  "c",
					},
				},
			},
			args: args{
				s: "b",
			},
			wantErr: false,
		},
		{
			name: "Delete root in root-only tree",
			tree: Tree{
				Root: &Node{
					Value: "a",
					Data:  "a",
				},
			},
			want: Tree{
				Root: nil,
			},
			args: args{
				s: "a",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tree.Delete(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tree.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(tt.tree, tt.want) {
				t.Errorf("Tree.Delete() = %v, want %v", tt.tree, tt.want)
			}
		})
	}
}
