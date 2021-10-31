package main

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	listener, err := net.Listen("tcp", ":19191")
	if err != nil {
		panic(err)
	}
	http.Serve(listener, &MyHandler{engine})
}

type MyHandler struct {
	currentEngine *gin.Engine
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.currentEngine.ServeHTTP(w, r)
}

func (h *MyHandler) SetNewHandler(newEngine *gin.Engine) {
	h.currentEngine = newEngine
}
