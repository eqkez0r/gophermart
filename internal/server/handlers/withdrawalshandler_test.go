package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestWithdrawalsHandler(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger *zap.SugaredLogger
		store  WithdrawalsProvider
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
			if got := WithdrawalsHandler(tt.args.ctx, tt.args.logger, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithdrawalsHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
