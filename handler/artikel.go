package handler

import (
	"fmt"
	"mnc-portal/artikel"
	"mnc-portal/helper"
	"mnc-portal/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type artikelHandler struct {
	service artikel.Service
}

func NewArtikelHandler(service artikel.Service) *artikelHandler {
	return &artikelHandler{service}
}

func (h *artikelHandler) GetArtikels(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	artikels, err := h.service.GetArtikels(userID)
	if err != nil {
		response := helper.APIResponse("Error to get artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of artikel", http.StatusOK, "success", artikel.FormatterArtikels(artikels))
	c.JSON(http.StatusOK, response)
}

func (h *artikelHandler) GetArtikel(c *gin.Context) {
	var input artikel.GetArtikelDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	artikelDetail, err := h.service.GetArtikelByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Artikel detail", http.StatusOK, "success", artikel.FormatArtikelDetail(artikelDetail))
	c.JSON(http.StatusOK, response)
}

func (h *artikelHandler) CreateArtikel(c *gin.Context) {
	var input artikel.CreateArtikelInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create artikel", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newArtikel, err := h.service.CreateArtikel(input)
	if err != nil {
		response := helper.APIResponse("Failed to create artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create artikel", http.StatusOK, "success", artikel.FormatterArtikel(newArtikel))
	c.JSON(http.StatusOK, response)
}

func (h *artikelHandler) UpdateArtikel(c *gin.Context) {
	var inputID artikel.GetArtikelDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData artikel.CreateArtikelInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update artikel", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedArtikel, err := h.service.UpdateArtikel(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update artikel", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update artikel", http.StatusOK, "success", artikel.FormatterArtikel(updatedArtikel))
	c.JSON(http.StatusOK, response)
}

func (h *artikelHandler) UploadImage(c *gin.Context) {
	var input artikel.CreateArtikelImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to upload artikel image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload artikel image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload artikel image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveArtikelImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload artikel image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("artikel image successfuly uploded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
