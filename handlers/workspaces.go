package handlers

import (
	"ExpenseTracker/config"
	"ExpenseTracker/database"
	"ExpenseTracker/models"
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

	if err := c.QueryParser(req); err != nil {
		return err
	}

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

	type EmailChan struct {
		Email string
		Msg   string
		Err   string
	}

	emailChan := make(chan *EmailChan, len(req.Emails))

	var wg sync.WaitGroup

	for _, email := range req.Emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			user := &models.User{}
			if err := database.PsqlDb.Model(user).Where("email = ?", email).Select(ctx); err != nil {
				emailChan <- &EmailChan{Email: email, Err: "no user found with this email", Msg: "Could not send invitation"}
				return
			}

			if err := database.PsqlDb.Model(&models.WorkspaceMember{}).Where("user_id = ? AND workspace_id = ?", user.ID, req.WorkspaceId).Select(ctx); err == nil {
				emailChan <- &EmailChan{Email: email, Err: "user is already a member of this workspace", Msg: "Could not send invitation"}
				return
			}
			res, _ := database.PsqlDb.Model(&models.WorkspaceInvitation{}).Where("email = ? AND workspace_id = ?", email, req.WorkspaceId).Delete(ctx)
			if res.RowsAffected() > 0 {
				emailChan <- &EmailChan{Email: email, Msg: "Invitation resent", Err: ""}
				return
			}

			emailChan <- &EmailChan{Email: email, Msg: "", Err: ""}

		}(email)
	}

	go func() {
		wg.Wait()
		close(emailChan)
	}()

	emailRes := &[]EmailChan{}

	for v := range emailChan {
		*emailRes = append(*emailRes, *v)
		if v.Err != "" {
			continue
		}

		wi := &models.WorkspaceInvitation{
			Email:       v.Email,
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
		if v.Msg == "" {
			v.Msg = "Invitation sent"
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Invitations sent",
		// return everything from emailChan
		"emails": emailRes,
	},
	)
}

func IssueInviteCode(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	req := &config.IssueInviteCodeRequest{}

	if err := c.QueryParser(req); err != nil {
		return err
	}

	if err := config.ValidationResponse(req); err != nil {
		return err
	}

	inviteCode := &models.WorkspaceInviteCode{
		WorkspaceId: req.WorkspaceId,
		IssuedBy:    uint64(userId),
	}

	if err := inviteCode.RenewDuration(); err != nil {
		return err
	}

	if err := inviteCode.GenerateCode(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(inviteCode).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"code":   inviteCode.Code,
	})

}

func JoinWorkspaceByCode(c *fiber.Ctx) error {
	ctx := c.Context()

	claimData := c.Locals("jwtClaims").(jwt.MapClaims)

	if claimData == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Jwt was bypassed",
		})
	}

	userId := claimData["id"].(float64)

	req := &config.JoinWorkspaceByCodeRequest{}

	if err := c.QueryParser(req); err != nil {
		return err
	}

	if err := config.ValidationResponse(req); err != nil {
		return err
	}

	inviteCode := &models.WorkspaceInviteCode{}

	if err := database.PsqlDb.Model(inviteCode).Where("code = ?", req.Code).Select(ctx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No code found",
		})
	}

	if err := inviteCode.ValidateCode(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if err := database.PsqlDb.Model(&models.WorkspaceMember{}).Where("user_id = ? AND workspace_id = ?", uint64(userId), inviteCode.WorkspaceId).Select(ctx); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "User is already a member of this workspace",
		})
	}

	wm := &models.WorkspaceMember{
		UserId:      uint64(userId),
		WorkspaceId: inviteCode.WorkspaceId,
	}

	if err := wm.BeforeInsert(); err != nil {
		return err
	}

	if _, err := database.PsqlDb.Model(wm).Insert(ctx); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   wm,
	})
}
