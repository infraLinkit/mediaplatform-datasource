package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource/config"
	"github.com/infraLinkit/mediaplatform-datasource/entity"
)

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(bytes)
}

var (
	uploadAllowedExt = map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
		".pdf": true, ".csv": true, ".xls": true, ".xlsx": true,
	}
	uploadFolderPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

// UploadFile handles generic file/image uploads to ./storage/{folder}
func (h *IncomingHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		// Try alternative form field names
		file, err = c.FormFile("image")
		if err != nil {
			file, err = c.FormFile("sample_image_banner")
		}
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "No file uploaded. Expected 'file', 'image', or 'sample_image_banner' in multipart form data.",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !uploadAllowedExt[ext] {
		// /storage is served publicly from the app origin, so .html/.svg/.js would be stored XSS.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "File type not allowed",
		})
	}

	folder := c.FormValue("folder")
	if folder == "" {
		folder = c.Query("folder", "uploads")
	}
	// Sanitize folder name: a single path segment only
	if !uploadFolderPattern.MatchString(folder) {
		folder = "uploads"
	}

	uploadDir := filepath.Join(".", "storage", folder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.Logs.WithError(err).Error("Failed to create upload directory")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to create storage directory",
		})
	}

	fileName := fmt.Sprintf("%s_%d_%s%s", folder, time.Now().Unix(), randomHex(4), ext)
	dstPath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, dstPath); err != nil {
		h.Logs.WithError(err).Error("Failed to save uploaded file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to save uploaded file",
		})
	}

	relativeURL := fmt.Sprintf("/storage/%s/%s", folder, fileName)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "File uploaded successfully",
		"url":     relativeURL,
		"data": fiber.Map{
			"url":       relativeURL,
			"filename":  fileName,
			"size":      file.Size,
			"mime_type": file.Header.Get("Content-Type"),
		},
	})
}

// StoreLPDesign stores a new LP Design Type and its banner image
// StoreLPDesign follows CampaignMainstreamSettingController::storeLPDesign. The file type comes from
// the uploaded bytes, and the stored name from the server, because /storage is served publicly.
func (h *IncomingHandler) StoreLPDesign(c *fiber.Ctx) error {
	mainstreamLPType := c.FormValue("mainstream_lp_type")
	description := c.FormValue("description")
	errs := lpDesignErrors(mainstreamLPType, description)

	extension := ""
	file, err := c.FormFile("sample_image_banner")
	if err != nil {
		errs = append(errs, "The sample image banner field is required.")
	} else {
		opened, openErr := file.Open()
		if openErr != nil {
			errs = append(errs, "The sample image banner failed to upload.")
		} else {
			head := make([]byte, 512)
			n, _ := opened.Read(head)
			opened.Close()
			extension = lpBannerExtension(http.DetectContentType(head[:n]))
			if extension == "" {
				errs = append(errs, "The sample image banner field must be a file of type: jpeg, png, jpg, gif.")
			} else if file.Size > 2048*1024 {
				errs = append(errs, "The sample image banner field must not be greater than 2048 kilobytes.")
			}
		}
	}
	if len(errs) > 0 {
		return validationFailed(c, errs)
	}

	uploadDir := filepath.Join(".", "storage", "lp_design_banners")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.Logs.WithError(err).Error("Failed to create LP banner directory")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to save LP Design Type"})
	}
	fileName := fmt.Sprintf("banner_%d_%s%s", time.Now().Unix(), randomHex(4), extension)
	dstPath := filepath.Join(uploadDir, fileName)
	if err := c.SaveFile(file, dstPath); err != nil {
		h.Logs.WithError(err).Error("Failed to save LP banner image file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to save LP Design Type"})
	}

	lpDesign := entity.LpDesignType{
		MainstreamLPType:  mainstreamLPType,
		Description:       description,
		SampleImageBanner: fmt.Sprintf("/storage/lp_design_banners/%s", fileName),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := h.DB.Create(&lpDesign).Error; err != nil {
		os.Remove(dstPath)
		h.Logs.WithError(err).Error("Failed to create LPDesignType in DB")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save LP Design Type",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "LP Design Type saved successfully!",
		"data":    lpDesign,
	})
}

// GetLPDesignTypes retrieves all registered LP Design Types
func (h *IncomingHandler) GetLPDesignTypes(c *fiber.Ctx) error {
	var list []entity.LpDesignType
	if err := h.DB.Order("id DESC").Find(&list).Error; err != nil {
		h.Logs.WithError(err).Error("Failed to fetch LPDesignTypes")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch LP Design Types",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{
		Code:    fiber.StatusOK,
		Message: config.OK_DESC,
		Data:    list,
	})
}

// DeleteLPDesignType removes an LP Design Type by ID
func (h *IncomingHandler) DeleteLPDesignType(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	if err := h.DB.Delete(&entity.LpDesignType{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete LP Design Type",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "LP Design Type deleted successfully",
	})
}
