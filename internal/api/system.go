package api

import (
	"bc-opp-api/internal/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSystemInfo(c *gin.Context) {
	now := time.Now()
	d := int(now.Sub(model.Creatd).Minutes())
	msg := fmt.Sprintf("啟動時間: %s\n", model.Creatd.Format("2006/01/02 15:04:05"))
	msg += fmt.Sprintf("現在時間: %s\n", now.Format("2006/01/02 15:04:05"))
	msg += fmt.Sprintf("已啟動: %v分鐘", d)
	c.String(200, msg)
}
