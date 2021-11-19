package models

import "testing"

func TestDeriveCacheKeyFromRequest(t *testing.T) {
	type args struct {
		request *SearchRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveCacheKeyFromRequest(tt.args.request); got != tt.want {
				t.Errorf("DeriveCacheKeyFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
