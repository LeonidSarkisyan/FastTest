package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DOMAIN = "http://localhost:8080/"

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

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

func (h *Handler) GroupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "groups.html", gin.H{})
}

func (h *Handler) OneGroupPage(c *gin.Context) {
	groupIDStr := c.Param("group_id")

	userID := c.GetInt("userID")
	groupID, err := strconv.Atoi(groupIDStr)

	if err != nil {
		c.HTML(422, "error.html", gin.H{
			"error": "Ошибка, такой ID не может быть у группы",
		})
		c.Abort()
		return
	}

	group, err := h.GroupService.Get(groupID, userID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такая группа не найдена",
		})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "one_group.html", gin.H{
		"name": group.Name,
	})
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

func (h *Handler) OneTestAccessPage(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")

	test, err := h.TestService.Get(testID, userID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такой тест не найден",
		})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "access.html", gin.H{
		"title": test.Title,
		"count": test.Count,
		"url":   fmt.Sprintf("/p/tests/%d", test.ID),
	})
}

func (h *Handler) ResultPage(c *gin.Context) {
	c.HTML(http.StatusOK, "results.html", gin.H{})
}

func (h *Handler) OneResultPage(c *gin.Context) {
	userID := c.GetInt("userID")
	resultID := MustID(c, "result_id")

	access, err := h.TestService.GetAccess(userID, resultID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такой результат не найден",
		})
		c.Abort()
		return
	}

	test, err := h.TestService.Get(access.TestID, userID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такой тест не найден",
		})
		c.Abort()
		return
	}

	group, err := h.GroupService.Get(access.GroupID, userID)

	if err != nil {
		c.HTML(404, "error.html", gin.H{
			"error": "Ошибка 404, такой тест не найден",
		})
		c.Abort()
		return
	}

	access.DateStart = strings.ReplaceAll(
		strings.ReplaceAll(access.DateStart, "T00:00:00Z", ""), "-", ".")

	access.DateEnd = strings.ReplaceAll(
		strings.ReplaceAll(access.DateEnd, "T00:00:00Z", ""), "-", ".")

	c.HTML(http.StatusOK, "one_result.html", gin.H{
		"title":   test.Title,
		"test":    test,
		"group":   group,
		"access":  access,
		"url":     fmt.Sprintf("/p/tests/%d", test.ID),
		"urlPass": DOMAIN + fmt.Sprintf("passing/%d", access.ID),
	})
}

func (h *Handler) PassingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "passing.html", gin.H{})
}
