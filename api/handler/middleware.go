package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type user struct {
	Id   int
	Role string
}

func (h *HandlerApi) AuthMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"reason": "empty auth header",
		})
		return
	}

	request, err := h.Nats.Request("user.middleware", []byte(header), 5*time.Second)
	if err != nil {
		c.AbortWithStatusJSON(500, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send a request. Error is:", err.Error())
		return
	}

	var u user
	if err := json.Unmarshal(request.Data, &u); err != nil {
		c.AbortWithStatusJSON(400, map[string]any{
			"message": string(request.Data),
		})
		h.Logger.Println("couldn't unmarshal the data. Error is:", string(request.Data))
		return
	}
	//fmt.Println(u.Id)
	//fmt.Println(u.Role)

	c.Set("userId", u.Id)
	c.Set("userRole", u.Role)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("userId not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("invalid type of userId")
	}

	return idInt, nil
}

func GetUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get("userRole")
	if !ok {
		return "", errors.New("userRole not found")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("invalid type of userRole")
	}

	return roleStr, nil
}
