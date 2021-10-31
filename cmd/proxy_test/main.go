package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	enginx := gin.New()
	enginx.POST("/addPath", func(context *gin.Context) {
		data := &struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}{}
		if err := context.ShouldBind(data); err != nil {
			context.JSON(http.StatusBadRequest, nil)
			return
		}
		switch data.Method {
		case http.MethodPost:
			enginx.POST(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodGet:
			enginx.GET(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodPut:
			enginx.PUT(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodDelete:
			enginx.DELETE(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		}
	})
	enginx.POST("deletePath", func(context *gin.Context) {
		data := &struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}{}
		if err := context.ShouldBind(data); err != nil {
			context.JSON(http.StatusBadRequest, nil)
			return
		}
		switch data.Method {
		case http.MethodPost:
			enginx.POST(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodGet:
			enginx.GET(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodPut:
			enginx.PUT(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		case http.MethodDelete:
			enginx.DELETE(data.Path, func(context *gin.Context) {
				fmt.Println("Request %s.Method %s", data.Path, data.Method)
			})
		}
	})

	listener, err := net.Listen("tcp", ":19292")
	if err != nil {
		panic(err)
	}
	err = enginx.RunListener(listener)
	if err != nil {
		panic(err)
	}
}

func AddReservedProxy(ctx context.Context, remoteUrl, localPath, method string, ginRouter *gin.RouterGroup) error {
	urlObj, err := url.Parse(remoteUrl)
	if err != nil {
		return err
	}
	reservedProxy := httputil.NewSingleHostReverseProxy(urlObj)
	reservedProxy.Director = func(request *http.Request) {
		//实现http复写请求header的部分逻辑
		request.Host = urlObj.Host
	}
	switch method {
	case http.MethodPost:
		ginRouter.POST(localPath, func(context *gin.Context) {
		})
	case http.MethodGet:
		ginRouter.GET(localPath, func(context *gin.Context) {
		})
	case http.MethodPut:
		ginRouter.PUT(localPath, func(context *gin.Context) {
		})
	case http.MethodDelete:
		ginRouter.DELETE(localPath, func(context *gin.Context) {
		})
	}
	return nil
}
