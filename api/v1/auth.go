package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirasildev/note_project/api/models"
	emailPkg "github.com/mirasildev/note_project/pkg/email"
	"github.com/mirasildev/note_project/pkg/utils"
	"github.com/mirasildev/note_project/storage/repo"
)

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrIncorrectCode    = errors.New("incorrect verification code")
	ErrCodeExpired      = errors.New("varification code has been expired")
	ErrEmailExists      = errors.New("email already exists")
)

const (
	RegisterCode      = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

// @Router /auth/register [post]
// @Summary Register a user
// @Description Register a user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Register(c *gin.Context) {
	var req models.RegisterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.storage.User().GetByEmail(req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user := repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.inMemory.Set("user_"+user.Email, string(userData), 10*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(RegisterCode, req.Email)
		if err != nil {
			fmt.Printf("Failed to send verification code: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})

}

func (h *handlerV1) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = h.inMemory.Set(key+email, code, time.Minute)
	if err != nil {
		return err
	}

	err = emailPkg.SendEmail(h.cfg, &emailPkg.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: emailPkg.VerificationEmail,
	})
	if err != nil {
		return err
	}

	return nil
}

// @Router /auth/verify [post]
// @Summary Verify user
// @Description Verify user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Verify(c *gin.Context) {
	var req models.VerifyRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		// log.Print(err)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userData, err := h.inMemory.Get("user_" + req.Email)
	if err != nil {
		// log.Print(err)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user repo.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	code, err := h.inMemory.Get(RegisterCode + user.Email)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		// log.Print(err)
		c.JSON(http.StatusForbidden, errorResponse(ErrIncorrectCode))
		return
	}

	result, err := h.storage.User().Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth/login [post]
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Login(c *gin.Context) {
	var req models.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.User().GetByEmail(req.Email)
	if err != nil {
		log.Println("Hello World 1")
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, result.Password)
	if err != nil {
		log.Println("Hello World 1")
		c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID: result.ID,
		Email: result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}
