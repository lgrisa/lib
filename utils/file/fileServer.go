package file

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lgrisa/lib/config"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewSimpleFileServer 创建一个简单的文件服务器
func NewSimpleFileServer(StaticPath string, port int) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(static.Serve(StaticPath, static.LocalFile("."+StaticPath, true)))

	router.HEAD("/download", func(c *gin.Context) {
		name := c.Query("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "name is empty",
			})
			return
		}

		fileName := name

		_, errByOpenFile := os.Open(fileName)
		//非空处理
		if errByOpenFile != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": fmt.Sprintf("文件不存在:%v", fileName),
				"error":   "资源不存在",
			})
			return
		}

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+name)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileName)
	})

	router.GET("/download", func(c *gin.Context) {
		name := c.Query("name")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "name is empty",
			})
			return
		}

		fileName := name

		_, errByOpenFile := os.Open(fileName)
		//非空处理
		if errByOpenFile != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": fmt.Sprintf("文件不存在:%v", fileName),
				"error":   "资源不存在",
			})
			return
		}

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+name)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileName)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	go func() {

		certFile := config.StartConfig.HttpConfig.CertFile
		keyFile := config.StartConfig.HttpConfig.KeyFile

		if certFile != "" && keyFile != "" {
			logutil.LogInfoF("https server start at :%v", port)

			if err := srv.ListenAndServeTLS("conf/test46.sgameuser.com.pem", "conf/test46.sgameuser.com.key"); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					logutil.LogErrorF("https server start fail:%v", err)
				}
			}

			logutil.LogInfoF("https server closed")
		} else {
			logutil.LogInfoF("http server start at :%v", port)

			if err := srv.ListenAndServe(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					logutil.LogErrorF("http server start fail:%v", err)
				}
			}

			logutil.LogInfoF("httpServer closed")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logutil.LogInfoF("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logutil.LogErrorF("Server Shutdown:%v", err)
	}
}
