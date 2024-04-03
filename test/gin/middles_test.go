package gin

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddle(t *testing.T) {
	// 新建一个没有任何默认中间件的路由
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		t.Log("The first middle...")
	})
	r.Use(func(c *gin.Context) {
		t.Log("The second middle...")
		c.AbortWithStatus(666)
		return
	})
	r.GET("/", func(c *gin.Context) {
		c.String(200, "test middle...")

		return
	})
	r.Run(":8000")
}
