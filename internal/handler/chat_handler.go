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

	g, ctx := errgroup.WithContext(ctx)

	var userResp *dto.UserResponse
	var permissionsResp *dto.PermissionsResponse

	g.Go(func() error {
		us, err := u.userService.GetUser(ctx, id)
		if err != nil {
			return apperrors.Wrap(err, "user service failed")
		}
		userResp = us
		return nil
	})

	g.Go(func() error {
		perm, err := u.permissionsService.CheckAccess(ctx, id)
		if err != nil {
			return apperrors.Wrap(err, "permissions service failed")
		}
		permissionsResp = perm
		return nil
	})

	vectorChannel := make(chan string)
	go func() {
		vector, err := u.vectorMemoryService.GetContext(ctx, id)
		if err == nil {
			select {
			case vectorChannel <- vector:
			default:
			}
		}
	}()

	err := g.Wait()

	if err != nil {
		var appErr *apperrors.AppError

		// Если это контекстный таймаут (context deadline exceeded)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "timeout",
				"message": "выполнение больше 200 мс",
			})
			return
		}

		if errors.As(err, &appErr) {
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
	default:
	}

	resp := ChatSummaryResponse{
		User:        userResp,
		Permissions: permissionsResp,
		Vector:      vectorResp,
	}

	c.JSON(http.StatusOK, resp)
}
