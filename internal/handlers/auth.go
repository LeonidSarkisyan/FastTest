package handlers

import (
	"App/internal/models"
	"App/internal/service"
	"App/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

const Domain = "https://фаст-тест.рф/"
const DomainWithPort = "https://фаст-тест.рф/"

func (h *Handler) CreateAccount(c *gin.Context) {
	var in models.UserIn

	if err := c.BindJSON(&in); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	if h.Exist(in.Email) {
		SendErrorResponse(c, 422, service.UserAlreadyExists.Error())
		return
	}

	exists := true
	var code []int64

	var mu sync.Mutex

	for exists {
		code = utils.GenerateSixDigitNumber(1)
		mu.Lock()
		_, exists = h.EmailCodeMap[code[0]]
		mu.Unlock()
	}

	mu.Lock()
	h.EmailCodeMap[code[0]] = in
	mu.Unlock()

	err := h.EmailClient.SendCodeToEmail(in.Email, code[0])

	if err != nil {
		SendErrorResponse(c, 400, "некорректная почта")
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, "")
}

func (h *Handler) ConfirmAccount(c *gin.Context) {
	code := int64(MustID(c, "code"))

	var mu sync.Mutex
	mu.Lock()
	in, ok := h.EmailCodeMap[code]
	mu.Unlock()

	if !ok {
		SendErrorResponse(c, 404, "неизвестный код")
		c.Abort()
		return
	}

	err := h.UserService.Register(in)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	token, err := h.UserService.Login(in)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.SetCookie("Authorization", token, 2592000, "/", Domain, false, true)
	c.HTML(200, "success_register.html", gin.H{})
}

func (h *Handler) Login(c *gin.Context) {
	var in models.UserIn

	if err := c.BindJSON(&in); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	token, err := h.UserService.Login(in)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.SetCookie("Authorization", token, 2592000, "/", Domain, false, true)
	c.AbortWithStatusJSON(http.StatusNoContent, "")
}
