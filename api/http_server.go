package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/negativations/modules/negativation/application"
	"github.com/negativations/modules/negativation/domain"
	"github.com/samora/gin-jsend"
	"net/http"
)

func createHttpServer(port int, negativationController *application.NegativationController) *http.Server {
	handler := createHandler(negativationController)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
}

func createHandler(negativationController *application.NegativationController) http.Handler {
	handler := gin.Default()
	handler.GET("negativations", func(context *gin.Context) {
		negativations, err := negativationController.GetByCPF(context.Query("cpf"))
		sendJsendResponse(context, negativations, err)
	})
	return handler
}

func sendJsendResponse(context *gin.Context, negativations []*domain.Negativation, err error) {
	if err != nil {
		jsend.Error(context, http.StatusInternalServerError, err.Error(), 0, nil)
		return
	}
	jsend.Success(context, http.StatusOK, negativations)
}
