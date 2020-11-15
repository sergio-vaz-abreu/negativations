package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/negativations/modules/negativation/application"
	"github.com/samora/gin-jsend"
	"net/http"
)

func createHttpServer(port int, negativationController *application.NegativationController, legacyNegativationController *application.LegacyNegativationController) *http.Server {
	handler := createHandler(negativationController, legacyNegativationController)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
}

func createHandler(negativationController *application.NegativationController, legacyNegativationController *application.LegacyNegativationController) http.Handler {
	handler := gin.Default()
	handler.GET("negativation", func(context *gin.Context) {
		negativations, err := negativationController.GetByCPF(context.Query("cpf"))
		sendJsendResponse(context, negativations, err)
	})
	handler.POST("negativation/synchronize", func(context *gin.Context) {
		err := legacyNegativationController.Synchronize()
		sendJsendResponse(context, nil, err)
	})
	return handler
}

func sendJsendResponse(context *gin.Context, data interface{}, err error) {
	if err != nil {
		jsend.Error(context, http.StatusInternalServerError, err.Error(), 0, nil)
		return
	}
	jsend.Success(context, http.StatusOK, data)
}
