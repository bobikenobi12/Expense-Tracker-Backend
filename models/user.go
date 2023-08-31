package models

import "time"

type UserSecrets struct {
	ID        int64  `json:"id"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	Email         string       `json:"email"`
	CountryCode   string       `json:"country_code"`
	CreatedAt     time.Time    `json:"created_at" pg:"default:now()"`
	UpdatedAt     time.Time    `json:"updated_at" pg:"default:now()"`
	UserSecretsId int64        `json:"user_secrets_id"`
	UserSecrets   *UserSecrets `json:"user_secrets" pg:"rel:has-one"`
}

func (u *User) BeforeInsert() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}

func (us *UserSecrets) BeforeInsert() error {
	us.CreatedAt = time.Now()
	us.UpdatedAt = time.Now()
	return nil
}

func (us *UserSecrets) BeforeUpdate() error {
	us.UpdatedAt = time.Now()
	return nil
}
