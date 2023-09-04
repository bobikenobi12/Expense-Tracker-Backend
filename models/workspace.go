package models

import "time"

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
