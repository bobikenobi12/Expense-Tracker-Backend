package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateWorkspace(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	req := &config.CreateWorkspaceRequest{}

	if err := c.BodyParser(req); err != nil {
		return err
	}

	if err := config.ValidationResponse(req); err != nil {
		return err
	}

	workspace := &models.Workspace{
		Name:    req.Name,
		OwnerId: uint64(userId),
	}

	if err := workspace.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(workspace).Insert(ctx); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(&models.WorkspaceMember{
		UserId:      uint64(userId),
		WorkspaceId: workspace.ID,
	}).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   workspace,
	})
}

func GetWorkspaces(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	// ownedWorkspaces := []models.Workspace{}

	// if err := database.PsqlDb.Model(&ownedWorkspaces).Where("owner_id = ?", userId).Select(ctx); err != nil {
	// 	return err
	// }

	workspaces := []models.Workspace{}

	if err := database.PsqlDb.Model(&workspaces).Column("workspace.*").Join("JOIN workspace_members ON workspace_members.workspace_id = workspace.id").Where("workspace_members.user_id = ?", userId).Select(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   workspaces,
	})
}
