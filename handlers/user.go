package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"ExpenseTracker/tools"

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

func UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	updateFields := &config.UpdateProfileRequest{}

	if err := c.BodyParser(updateFields); err != nil {
		return err
	}

	if err := config.ValidationResponse(updateFields); err != nil {
		return err
	}

	user := &models.User{}

	result, err := database.PsqlDb.Model(user).Set("name = ?", updateFields.Name).Set("country_code = ?", updateFields.CountryCode).Where("id = ?", userId).Update(ctx)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"data": fiber.Map{
			"name":         updateFields.Name,
			"country_code": updateFields.CountryCode,
		},
	})
}

func ChangePassword(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	changePassword := &config.ChangePasswordRequest{}

	if err := c.BodyParser(changePassword); err != nil {
		return err
	}

	if err := config.ValidationResponse(changePassword); err != nil {
		return err
	}

	userSecrets := &models.UserSecrets{
		ID: uint64(userId),
	}

	if err := database.PsqlDb.Model(userSecrets).WherePK().Select(ctx); err != nil {
		return err
	}

	match := tools.CheckPasswordHash(changePassword.OldPassword, userSecrets.Password)

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Incorrect password",
		})
	}

	hashedPassword, err := tools.HashPassword(changePassword.NewPassword)

	if err != nil {
		return err
	}

	result, err := database.PsqlDb.Model(userSecrets).Set("password = ?", hashedPassword).WherePK().Update(ctx)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password changed successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()

	claims := c.Locals("jwtClaims").(jwt.MapClaims)

	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claims["id"].(float64)

	user := &models.User{
		ID: uint64(userId),
	}

	if _, err := database.PsqlDb.Model(user).WherePK().Delete(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"msg":    "User deleted",
	})
}
