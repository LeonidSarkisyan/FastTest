package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) GetPassesAndStudents(c *gin.Context) { {
	userID := c.GetInt("userID")
	resultID := MustID(c, "result_id")

	passes, students, err := h.GetPassesAndStudents(resultID, userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"passes": passes,
		"students": students,
	})
}
