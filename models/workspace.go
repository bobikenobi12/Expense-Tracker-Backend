package models

import (
	"ExpenseTracker/tools"
	"errors"
	"time"
)

type Currency struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	IsoCode string `json:"iso_code"`
}
type Workspace struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	OwnerId   uint64 `json:"owner_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CurrencyWorkspace struct {
	ID          uint64 `json:"id"`
	CurrencyId  uint64 `json:"currency_id"`
	WorkspaceId uint64 `json:"workspace_id"`
	AddedBy     uint64 `json:"added_by"`
}

type WorkspaceMember struct {
	ID          uint64 `json:"id"`
	UserId      uint64 `json:"user_id"`
	WorkspaceId uint64 `json:"workspace_id"`
	JoinedOn    string `json:"joined_on"`
}

type WorkspaceInvitation struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	WorkspaceId uint64 `json:"workspace_id"`
	AddedBy     uint64 `json:"added_by"`
	Expires     string `json:"expires"`
}

type WorkspaceInviteCode struct {
	ID          uint64 `json:"id"`
	WorkspaceId uint64 `json:"workspace_id"`
	IssuedBy    uint64 `json:"issued_by"`
	Expires     string `json:"expires"`
	Code        string `json:"code"`
}

type CurrencyUser struct {
	ID         uint64 `json:"id"`
	CurrencyId uint64 `json:"currency_id"`
	UserId     uint64 `json:"user_id"`
}

func (w *Workspace) BeforeInsert() error {
	w.CreatedAt = time.Now().UTC().String()
	w.UpdatedAt = time.Now().UTC().String()
	return nil
}

func (w *Workspace) BeforeUpdate() error {
	w.UpdatedAt = time.Now().UTC().String()
	return nil
}

func (w *WorkspaceMember) BeforeInsert() error {
	w.JoinedOn = time.Now().UTC().String()
	return nil
}

func (w *WorkspaceInvitation) RenewDuration() error {
	w.Expires = time.Now().UTC().Add(time.Hour * 24).String()
	return nil
}

func (w *WorkspaceInviteCode) RenewDuration() error {
	w.Expires = time.Now().UTC().Add(time.Hour * 24).String()
	return nil
}

func (w *WorkspaceInviteCode) GenerateCode() error {
	code, err := tools.HashPassword(w.Expires)
	if err != nil {
		return err
	}
	w.Code = code
	return nil
}

func (w *WorkspaceInviteCode) ValidateCode() error {
	// if ok := tools.CheckPasswordHash(w.Code, w.Expires); !ok {
	// 	return errors.New("invalid code")
	// }

	layout := "2006-01-02 15:04:05.999999 -0700 MST"
	parsedTime, err := time.Parse(layout, w.Expires)
	if err != nil {
		return errors.New("invalid code")
	}

	if time.Now().UTC().After(parsedTime) {
		return errors.New("code expired")
	}
	return nil
}
