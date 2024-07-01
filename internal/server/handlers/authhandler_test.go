package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthHandler(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
