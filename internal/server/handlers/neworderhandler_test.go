package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewOrderHandler(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger *zap.SugaredLogger
		store  NewOrderProvider
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
			if got := NewOrderHandler(tt.args.ctx, tt.args.logger, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
