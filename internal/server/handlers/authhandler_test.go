package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

		})
	}
}
