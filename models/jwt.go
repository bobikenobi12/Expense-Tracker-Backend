package models

import "time"

type JwtWhitelist struct {
	ID        int64 `pg:",pk"`
	Jwt       string
	ExpiresAt time.Time
}

type JwtBlacklist struct {
	ID  int64 `pg:",pk"`
	Jwt string
}

func (j *JwtWhitelist) BeforeInsert() error {
	j.ExpiresAt = time.Now()
	return nil
}
