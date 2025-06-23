package db

import "errors"

var ErrUserExists = errors.New("user already exists")

type Role string

const (
	RoleSuperUser Role = "superuser"
	RoleAdmin     Role = "admin"
	RoleEditor    Role = "editor"
	RoleUser      Role = "user"
)

type Permission string

const (
	PermCreatePool       Permission = "create_pool"
	PermDeletePool       Permission = "delete_pool"
	PermCreateSchema     Permission = "create_schema"
	PermDeleteSchema     Permission = "delete_schema"
	PermCreateCollection Permission = "create_collection"
	PermDeleteCollection Permission = "delete_collection"
	PermRead             Permission = "read"
	PermWrite            Permission = "write"
)

var RolePermissions = map[Role][]Permission{
	RoleSuperUser: {
		PermCreatePool, PermDeletePool,
		PermCreateSchema, PermDeleteSchema,
		PermCreateCollection, PermDeleteCollection,
		PermRead, PermWrite,
	},
	RoleAdmin: {
		PermCreateSchema, PermDeleteSchema,
		PermCreateCollection, PermDeleteCollection,
		PermRead, PermWrite,
	},
	RoleEditor: {
		PermCreateCollection,
		PermRead, PermWrite,
	},
	RoleUser: {
		PermRead,
	},
}

type AuthManager struct {
	db *PostgresDB
}

func NewAuthManager(db *PostgresDB) *AuthManager {
	return &AuthManager{db: db}
}

func (am *AuthManager) RegisterUser(username, password string, role Role) error {
	return am.db.CreateUser(username, password, role)
}

func (am *AuthManager) HasPermission(username string, permission Permission) bool {
	role, err := am.db.GetUserRole(username)
	if err != nil {
		return false
	}

	permissions := RolePermissions[role]
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (am *AuthManager) ValidateUser(username, password string) (Role, error) {
	return am.db.ValidateUser(username, password)
}
