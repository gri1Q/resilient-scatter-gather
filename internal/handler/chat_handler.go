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

// ChatHandler обработчик маршрута
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
		//если запрос к пользователю не уложился за 10 мс то таймаут
		uCtx, uCancel := context.WithTimeout(gCtx, 10*time.Millisecond)
		defer uCancel()

		us, err := u.userService.GetUser(uCtx, id)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return apperrors.Wrap(apperrors.ErrTimeout, "user service timeout 10ms")
			}
			return apperrors.Wrap(err, "user service failed")
		}
		userResp = us
		return nil
	})

	g.Go(func() error {
		//Делаем индивидуальный контекст, что если запрос к правам не уложился за 50 мс то таймаут
		pCtx, pCancel := context.WithTimeout(gCtx, 50*time.Millisecond)
		defer pCancel()
		perm, err := u.permissionsService.CheckAccess(pCtx, id)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return apperrors.Wrap(apperrors.ErrTimeout, "permissions service timeout 50ms")
			}
			return apperrors.Wrap(err, "permissions service failed")
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
		var appErr *apperrors.AppError

		if errors.As(err, &appErr) {
			if errors.Is(appErr, apperrors.ErrTimeout) {
				// конкретный таймаут сервиса
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "timeout",
					"message": appErr.Error(), // покажет что именно таймаут User или Permissions
				})
				return
			}
			if errors.Is(appErr, apperrors.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not_found", "message": appErr.Error()})
				return
			}
			if errors.Is(appErr, apperrors.ErrPermission) {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": appErr.Error()})
				return
			}

			// Другие ошибки — 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": appErr.Error(),
			})
			return
		}

		// fallback — внутренняя ошибка
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	// Пытаемся считать vector, если он уже пришёл.
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
