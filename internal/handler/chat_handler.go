package handler

import (
	"context"
	"errors"
	"net/http"
	"testGolang/internal/apperrors"
	"testGolang/internal/dto"
	"testGolang/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type ChatSummaryResponse struct {
	User        *dto.UserResponse        `json:"user"`
	Permissions *dto.PermissionsResponse `json:"permissions"`
	Vector      *string                  `json:"vector,omitempty"`
}

type ChatHandler struct {
	userService         *service.UserService
	permissionsService  *service.PermissionsService
	vectorMemoryService *service.VectorMemoryService
}

func NewChatHandler(userService *service.UserService, permissionsService *service.PermissionsService, vectorMemoryService *service.VectorMemoryService) *ChatHandler {
	return &ChatHandler{userService: userService, permissionsService: permissionsService, vectorMemoryService: vectorMemoryService}
}

// GetChatSummary получаем сводку по чату
func (u *ChatHandler) GetChatSummary(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 200*time.Millisecond)
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)

	var userResp *dto.UserResponse
	var permissionsResp *dto.PermissionsResponse

	g.Go(func() error {
		uCtx, uCancel := context.WithTimeout(gCtx, 10*time.Millisecond)
		defer uCancel()

		us, err := u.userService.GetUser(uCtx, id)
		if err != nil {
			return err
		}

		userResp = us
		return nil
	})

	g.Go(func() error {
		pCtx, pCancel := context.WithTimeout(gCtx, 50*time.Millisecond)
		defer pCancel()
		perm, err := u.permissionsService.CheckAccess(pCtx, id)

		if err != nil {
			return err
		}

		permissionsResp = perm
		return nil
	})

	vectorChannel := make(chan string)
	go func() {
		vector, err := u.vectorMemoryService.GetContext(ctx, id)
		if err != nil {
			return
		}
		select {
		case vectorChannel <- vector:
		case <-ctx.Done():
		}
	}()

	err := g.Wait()

	if err != nil {
		if errors.Is(err, apperrors.ErrTimeout) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "timeout",
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found", "message": err.Error()})
			return
		}
		if errors.Is(err, apperrors.ErrPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": err.Error()})
			return
		}

		// другие ошибки — Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	var vectorResp *string
	select {
	case v := <-vectorChannel:
		vectorResp = &v
	case <-ctx.Done():
		vectorResp = nil
	}

	resp := ChatSummaryResponse{
		User:        userResp,
		Permissions: permissionsResp,
		Vector:      vectorResp,
	}

	c.JSON(http.StatusOK, resp)
}
