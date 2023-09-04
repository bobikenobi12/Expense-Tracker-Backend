package models

import "time"

type Currency struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	IsoCode string `json:"iso_code"`
	AddedBy int64  `json:"added_by"`
}
type Workspace struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerId   int64  `json:"owner_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CurrencyWorkspace struct {
	ID          int64 `json:"id"`
	CurrencyId  int64 `json:"currency_id"`
	WorkspaceId int64 `json:"workspace_id"`
}

type WorkspaceMember struct {
	ID          int64 `json:"id"`
	UserId      int64 `json:"user_id"`
	WorkspaceId int64 `json:"workspace_id"`
}

type CurrencyUser struct {
	ID         int64 `json:"id"`
	CurrencyId int64 `json:"currency_id"`
	UserId     int64 `json:"user_id"`
}

func (w *Workspace) BeforeInsert() error {
	w.CreatedAt = time.Now().UTC().String()
	w.UpdatedAt = time.Now().UTC().String()
	return nil
}
