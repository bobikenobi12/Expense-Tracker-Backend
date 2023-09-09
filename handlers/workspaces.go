package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"
	"log"
	"sync"

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

	wm := &models.WorkspaceMember{
		UserId:      uint64(userId),
		WorkspaceId: workspace.ID,
	}

	if err := wm.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(wm).Insert(ctx); err != nil {
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

func InviteUsersToWorkspace(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	req := &config.InviteUsersToWorkspaceRequest{}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if err := config.ValidationResponse(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	errChan := make(chan string, len(req.Emails))
	emailChan := make(chan string, len(req.Emails))

	var wg sync.WaitGroup

	for _, email := range req.Emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			if err := database.CheckIfEmailExists(email); err != nil {
				errChan <- email
				return
			}
			// if err := database.PsqlDb.Model(&models.WorkspaceInvitation{}).Where("email = ?", email).Delete(ctx); err != nil {

			emailChan <- email
		}(email)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(emailChan)
	}()

	log.Println(<-errChan, "and ", <-emailChan)
	if len(errChan) == len(req.Emails) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No users found with the given emails",
		})
	}

	for email := range emailChan {
		if email == <-errChan {
			log.Println("skipping")
			continue
		}
		wi := &models.WorkspaceInvitation{
			Email:       email,
			WorkspaceId: req.WorkspaceId,
			AddedBy:     uint64(userId),
		}
		if err := wi.RenewDuration(); err != nil {
			return err
		}
		if _, err := database.PsqlDb.Model(wi).Insert(ctx); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitations sent",
		"errors":  <-errChan,
	})
}
