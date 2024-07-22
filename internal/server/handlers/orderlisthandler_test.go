package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestOrderListHandler(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger *zap.SugaredLogger
		store  OrderListProvider
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
			if got := OrderListHandler(tt.args.ctx, tt.args.logger, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderListHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
