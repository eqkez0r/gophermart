package middleware

import (
	"compress/gzip"
	"github.com/eqkez0r/gophermart/internal/server/writers"
	e "github.com/eqkez0r/gophermart/pkg/error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var avaliableTypes = map[string]bool{
	"text/html":        true,
	"application/json": true,
	"html/text":        true,
}

func Gzip(
	logger *zap.SugaredLogger,
) gin.HandlerFunc {
	return func(context *gin.Context) {
		op := "Error in gzip decode handler: "
		//compress response
		ae := context.GetHeader("Accept-Encoding")
		ct := context.GetHeader("Content-Type")
		ac := context.GetHeader("Accept")
		if (strings.Contains(ae, "gzip")) &&
			(avaliableTypes[ct] || avaliableTypes[ac]) {
			w := context.Writer
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			wd := writers.GzipWriter{
				ResponseWriter: context.Writer,
				Writer:         gzipWriter,
			}
			context.Writer = wd
			context.Writer.Header().Set("Content-Encoding", "gzip")
		}

		//uncompressed request
		ce := context.GetHeader("Content-Encoding")
		if strings.Contains(ce, "gzip") {
			gzipReader, err := gzip.NewReader(context.Request.Body)
			if err != nil {
				logger.Error(e.Wrap(op, err))
				context.Status(http.StatusInternalServerError)
				return
			}
			defer gzipReader.Close()

			context.Request.Body = gzipReader
		}
		context.Next()
	}
}
