package service

import (
	"database/sql"

	"github.com/hako/branca"
)

// 서비스 핵심 로직. REST, GraphQL, RPC API 등 원하는거 사용
type Service struct {
	db    *sql.DB
	codec *branca.Branca
}

//DB와 Codec 생성자
func New(db *sql.DB, codec *branca.Branca) *Service {
	return &Service{
		db:    db,
		codec: codec,
	}
}
