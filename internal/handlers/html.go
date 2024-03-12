package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Users Mother Fucker!",
	})
}

func (h *Handler) AuthPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{})
}

func (h *Handler) Header(c *gin.Context) {
	c.HTML(http.StatusOK, "header.html", gin.H{})
}

func (h *Handler) TestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tests.html", gin.H{})
}

func (h *Handler) OneTestPage(c *gin.Context) {
	testIDStr := c.Param("test_id")

	userID := c.GetInt("userID")
	testID, err := strconv.Atoi(testIDStr)

	if err != nil {
		c.HTML(422, "error.html", gin.H{
			"error": "Ошибка, такой ID не может быть у теста",
		})
		c.Abort()
		return
	}

	test, err := h.TestService.Get(testID, userID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такой тест не найден",
		})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "one_test.html", gin.H{
		"title": test.Title,
	})
}
