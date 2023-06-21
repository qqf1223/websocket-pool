package nodeManager

import (
	"context"
	"testing"
)

func TestNodeManager_GetServers(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		n    *NodeManager
		args args
	}{
		// TODO: Add test cases.
		{
			name: "f",
			n:    &NodeManager{},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodeManager{}
			n.GetServers(tt.args.ctx)
		})
	}
}
