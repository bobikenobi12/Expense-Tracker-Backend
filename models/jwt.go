package models

type JwtWhitelist struct {
	ID        int64 `pg:",pk"`
	Jwt       string
	ExpiresAt int64
}

type JwtBlacklist struct {
	ID  int64 `pg:",pk"`
	Jwt string
}
