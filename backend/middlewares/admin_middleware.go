package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, exists := c.Get("rol")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token requerido",
			})
			c.Abort()
			return
		}

		if rol != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: se requiere rol administrador",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
