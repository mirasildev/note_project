package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirasildev/note_project/api/models"
	"github.com/mirasildev/note_project/storage/repo"
)

// @Security ApiKeyAuth
// @Router /notes [post]
// @Summary Create a note
// @Description Create a note
// @Tags notes
// @Accept json
// @Produce json
// @Param note body models.CreateNoteRequest true "Note"
// @Success 201 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateNote(c *gin.Context) {
	var req models.CreateNoteRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Note().Create(&repo.Note{
		UserID:      payload.UserID,
		Title:       req.Title,
		Description: *req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.Note{
		ID:          resp.ID,
		UserID:      resp.UserID,
		Title:       resp.Title,
		Description: &resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   &resp.UpdatedAt,
		DeletedAt:   resp.DeletedAt,
	})
}

// @Security ApiKeyAuth
// @Router /notes/{id} [get]
// @Summary Get note by id
// @Description Get note by id
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.Note().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Router /notes [get]
// @Summary Get all notes
// @Description Get all notes
// @Tags notes
// @Accept json
// @Produce json
// @Param filter query models.GetAllNotesParams false "Filter"
// @Success 200 {object} models.GetAllNotesResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllNotes(c *gin.Context) {
	req, err := validateGetAllNotesParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Note().GetAllNotes(&repo.GetAllNotesParams{
		Page:   req.Page,
		Limit:  req.Limit,
		UserID: req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getNotesResponse(result))
}

func validateGetAllNotesParams(c *gin.Context) (*models.GetAllNotesParams, error) {
	var (
		limit  int = 10
		page   int = 1
		err    error
		userID int
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("user_id") != "" {
		userID, err = strconv.Atoi(c.Query("user_id"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllNotesParams{
		Limit:  int32(limit),
		Page:   int32(page),
		UserID: int64(userID),
	}, nil
}

func getNotesResponse(data *repo.GetAllNotesResult) *models.GetAllNotesResponse {
	response := models.GetAllNotesResponse{
		Notes: make([]*models.Note, 0),
		Count: data.Count,
	}

	for _, note := range data.Notes {
		p := models.Note{
			ID:          note.ID,
			UserID:      note.UserID,
			Title:       note.Title,
			Description: &note.Description,
			CreatedAt:   note.CreatedAt,
			UpdatedAt:   &note.UpdatedAt,
			DeletedAt:   note.DeletedAt,
		}
		response.Notes = append(response.Notes, &p)
	}

	return &response
}

// @Router /notes/{id} [put]
// @Summary Update a note
// @Description Update a note
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param note body models.UpdateNoteRequest true "Note"
// @Success 200 {object} models.Note
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateNote(c *gin.Context) {
	var req models.UpdateNoteRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	updated, err := h.storage.Note().Update(&repo.Note{
			ID:          id,
			UserID:      req.UserID,
			Title:       req.Title,
			Description: *req.Description,
			UpdatedAt:   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Router /notes/{id} [delete]
// @Summary Delete a note
// @Description Delete a note
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteNote(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.storage.Note().Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted!",
	})
}
