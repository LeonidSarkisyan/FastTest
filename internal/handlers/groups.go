package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateGroup(c *gin.Context) {
	userID := c.GetInt("userID")

	var in models.GroupIn

	if err := c.BindJSON(&in); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	id, err := h.GroupService.Create(in, userID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *Handler) GetGroups(c *gin.Context) {
	userID := c.GetInt("userID")

	groups, err := h.GroupService.GetAll(userID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, responses.NewListResponse(groups))
}

func (h *Handler) GetGroup(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")

	group, err := h.GroupService.Get(groupID, userID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, group)
}

func (h *Handler) UpdateGroupTitle(c *gin.Context) {
	userID := c.GetInt("userID")
	groupID := MustID(c, "group_id")

	var groupUpdate models.GroupUpdate

	if err := c.BindJSON(&groupUpdate); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err := h.GroupService.UpdateTitle(userID, groupID, groupUpdate.Title)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, gin.H{})
}
