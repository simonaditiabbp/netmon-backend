package repository

import "gorm.io/gorm"

type Database interface {
	GetDB() *gorm.DB
}

type PostgresDatabase struct {
	DB *gorm.DB
}

func NewPostgresDatabase(db *gorm.DB) *PostgresDatabase {
	return &PostgresDatabase{DB: db}
}

func (p *PostgresDatabase) GetDB() *gorm.DB {
	return p.DB
}
