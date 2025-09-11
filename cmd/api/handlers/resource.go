package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProtectedHandler struct {
}

func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

// @Summary Protected resource
// @Description Example protected endpoint
// @Tags protected
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /protected [get]
func (*ProtectedHandler) Handler(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Protected content",
	})
}
