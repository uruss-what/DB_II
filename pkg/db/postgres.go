package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresDB struct {
	db *sql.DB
	sync.RWMutex
}

func NewPostgresDB() (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) CreateUser(username, password string, role Role) error {
	p.Lock()
	defer p.Unlock()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	var exists bool
	err = p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}
	if exists {
		return ErrUserExists
	}

	_, err = p.db.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)",
		username,
		string(hashedPassword),
		role,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (p *PostgresDB) ValidateUser(username, password string) (Role, error) {
	p.RLock()
	defer p.RUnlock()

	var (
		hashedPassword string
		role           Role
	)

	err := p.db.QueryRow(
		"SELECT password_hash, role FROM users WHERE username = $1",
		username,
	).Scan(&hashedPassword, &role)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	} else if err != nil {
		return "", fmt.Errorf("error querying user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	return role, nil
}

func (p *PostgresDB) GetUserRole(username string) (Role, error) {
	p.RLock()
	defer p.RUnlock()

	var role Role
	err := p.db.QueryRow(
		"SELECT role FROM users WHERE username = $1",
		username,
	).Scan(&role)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found")
	} else if err != nil {
		return "", fmt.Errorf("error querying user role: %v", err)
	}

	return role, nil
}
