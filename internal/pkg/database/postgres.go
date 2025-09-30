package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config contém as configurações do banco de dados
type Config struct {
	URL            string
	MaxConnections int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

// Connection mantém a conexão com o banco
type Connection struct {
	DB *sql.DB
}

// NewConnection cria uma nova conexão com o PostgreSQL
func NewConnection(cfg Config) (*Connection, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com o banco: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close fecha a conexão com o banco
func (c *Connection) Close() error {
	return c.DB.Close()
}

// Health verifica se a conexão está saudável
func (c *Connection) Health() error {
	return c.DB.Ping()
}
