package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateStudent(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")

	var student models.Student

	if err := c.BindJSON(&student); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	id, err := h.StudentService.Create(userID, groupID, student)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *Handler) AddStudentsFromExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}
	defer fileContent.Close()

	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")

	students, err := h.StudentService.CreateStudentFromExcel(userID, groupID, fileContent)
	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.JSON(http.StatusOK, responses.NewListResponse(students))
}

func (h *Handler) GetAllStudents(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")

	students, err := h.StudentService.GetAll(userID, groupID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
	}

	c.JSON(http.StatusOK, responses.NewListResponse(students))
}

func (h *Handler) DeleteStudent(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")
	studentID := MustID(c, "student_id")

	err := h.StudentService.Delete(userID, groupID, studentID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
