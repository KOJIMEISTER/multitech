package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProtectedHandler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Protected content",
	})
}
