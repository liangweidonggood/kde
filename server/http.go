package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kde/config"
	"log"
	"net/http"
	"time"
)

/**
处理http请求
比如主动向设备发指令
获取设备在线列表
*/
func HttpServer() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		//time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	v1 := router.Group("/api/v1")
	{
		v1.GET("deviceList", listHandler)
		v1.POST("deviceCmd", controlHandler)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GVA_CONFIG.System.WebPort),
		Handler: router,
	}
	// 服务连接
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

//设备列表l
func listHandler(c *gin.Context) {
	datalist := make([]Device, 0)
	for _, val := range connManager {
		var item Device
		item.Sn = val.sn
		item.Data = val.data
		item.FlowId = val.flowId
		item.MacType = val.macType
		item.LastTimeStamp = time.Unix(val.lastTimeStamp/1000, 0).Format("2006-01-02 15:04:05")
		datalist = append(datalist, item)
	}
	c.JSON(http.StatusOK, datalist)
}

//发送指令
func controlHandler(c *gin.Context) {
	//请求参数
	type DataReq struct {
		Sn    string `json:"sn" binding:"required"`
		Param string `json:"param" binding:"required"`
	}
	var json DataReq
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, val := range connManager {
		if val.sn == json.Sn {
			val.Conn.Write([]byte(json.Param + "\r"))
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": 0})
}
