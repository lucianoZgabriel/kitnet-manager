package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository"
	"github.com/lucianoZgabriel/kitnet-manager/internal/repository/sqlc"
)

// Compile-time check to ensure UserRepository implements repository.UserRepository
var _ repository.UserRepository = (*UserRepository)(nil)

// UserRepository implementa repository.UserRepository usando PostgreSQL via SQLC
type UserRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

// NewUserRepository cria uma nova instância do repositório de usuários
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// Create insere um novo usuário no banco de dados
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	params := sqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         sqlc.UserRole(user.Role),
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	created, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Atualiza o objeto user com dados do banco
	user.CreatedAt = created.CreatedAt
	user.UpdatedAt = created.UpdatedAt

	return nil
}

// GetByID busca um usuário pelo ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Não encontrou, retorna nil sem erro
		}
		return nil, err
	}

	return r.toDomain(dbUser), nil
}

// GetByUsername busca um usuário pelo username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(dbUser), nil
}

// List retorna todos os usuários
func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUsers), nil
}

// ListByRole retorna usuários filtrados por role
func (r *UserRepository) ListByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error) {
	dbUsers, err := r.queries.ListUsersByRole(ctx, sqlc.UserRole(role))
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(dbUsers), nil
}

// Update atualiza um usuário existente (role e is_active)
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	params := sqlc.UpdateUserParams{
		ID:        user.ID,
		Role:      sqlc.UserRole(user.Role),
		IsActive:  user.IsActive,
		UpdatedAt: user.UpdatedAt,
	}

	updated, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	user.UpdatedAt = updated.UpdatedAt
	return nil
}

// UpdatePassword atualiza apenas a senha do usuário
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	params := sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
		UpdatedAt:    time.Now(),
	}

	_, err := r.queries.UpdateUserPassword(ctx, params)
	return err
}

// UpdateLastLogin atualiza o timestamp do último login
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error {
	params := sqlc.UpdateLastLoginParams{
		ID: id,
		LastLoginAt: sql.NullTime{
			Time:  lastLogin,
			Valid: true,
		},
		UpdatedAt: time.Now(),
	}

	_, err := r.queries.UpdateLastLogin(ctx, params)
	return err
}

// Deactivate desativa um usuário
func (r *UserRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	params := sqlc.DeactivateUserParams{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	return r.queries.DeactivateUser(ctx, params)
}

// Activate ativa um usuário
func (r *UserRepository) Activate(ctx context.Context, id uuid.UUID) error {
	params := sqlc.ActivateUserParams{
		ID:        id,
		UpdatedAt: time.Now(),
	}

	return r.queries.ActivateUser(ctx, params)
}

// Delete remove um usuário do banco de dados
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

// ExistsByUsername verifica se já existe um usuário com o username
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.queries.UserExistsByUsername(ctx, username)
}

// Count retorna o total de usuários
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

// CountActive retorna o total de usuários ativos
func (r *UserRepository) CountActive(ctx context.Context) (int64, error) {
	return r.queries.CountActiveUsers(ctx)
}

// toDomain converte sqlc.User para domain.User
func (r *UserRepository) toDomain(dbUser sqlc.User) *domain.User {
	var lastLogin *time.Time
	if dbUser.LastLoginAt.Valid {
		lastLogin = &dbUser.LastLoginAt.Time
	}

	return &domain.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		PasswordHash: dbUser.PasswordHash,
		Role:         domain.UserRole(dbUser.Role),
		IsActive:     dbUser.IsActive,
		LastLoginAt:  lastLogin,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

// toDomainSlice converte []sqlc.User para []*domain.User
func (r *UserRepository) toDomainSlice(dbUsers []sqlc.User) []*domain.User {
	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomain(dbUser)
	}
	return users
}
