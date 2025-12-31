package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		spanName := fmt.Sprintf(
			"%s %s",
			c.Request.Method,
			c.FullPath(),
		)

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
