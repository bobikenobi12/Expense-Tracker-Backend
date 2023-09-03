package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UploadProfilePic(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	file, err := c.FormFile("profile_pic")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No file was uploaded",
		})
	}

	fileType := file.Header.Get("Content-Type")

	if fileType != "image/png" && fileType != "image/jpeg" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file type. Only png and jpeg are allowed",
		})
	}

	fileSize := file.Size

	if fileSize > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "File size too large. Max 5MB allowed",
		})
	}

	fileContent, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	uploadOutput, err := config.UploadToS3Bucket(&fileContent, file.Filename, fileType)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	insertPicResult, err := database.PsqlDb.Model(uploadOutput).Returning("id").Insert(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	result, err := database.PsqlDb.Model(&models.User{}).Set("profile_pic_id = ?", insertPicResult.RowsReturned()).Where("id = ?", userId).Update(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Profile pic uploaded successfully",
		"data": fiber.Map{
			"version_id": uploadOutput.VersionId,
			"location":   uploadOutput.Location,
			"etag":       uploadOutput.ETag,
			"key":        uploadOutput.Key,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	profilePicId := claimData["prof_pic_id"].(float64)

	user := &models.User{}

	err := database.PsqlDb.Model(user).Relation("ProfilePic").Where("profile_pic.id = ?", profilePicId).Select(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":              user.ID,
			"name":            user.Name,
			"email":           user.Email,
			"country_code":    user.CountryCode,
			"created_at":      user.CreatedAt,
			"updated_at":      user.UpdatedAt,
			"profile_pic":     user.ProfilePic.Location,
			"user_secrets":    user.UserSecrets,
			"user_secrets_id": user.UserSecretsId,
		},
	})

}
