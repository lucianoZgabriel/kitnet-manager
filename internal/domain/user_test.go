package domain

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		role        UserRole
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "valid admin user",
			username: "admin",
			password: "admin123",
			role:     UserRoleAdmin,
			wantErr:  false,
		},
		{
			name:     "valid manager user",
			username: "manager",
			password: "manager123",
			role:     UserRoleManager,
			wantErr:  false,
		},
		{
			name:     "valid viewer user",
			username: "viewer",
			password: "viewer123",
			role:     UserRoleViewer,
			wantErr:  false,
		},
		{
			name:        "username too short",
			username:    "ab",
			password:    "password123",
			role:        UserRoleAdmin,
			wantErr:     true,
			expectedErr: ErrInvalidUsername,
		},
		{
			name:        "password too short",
			username:    "admin",
			password:    "12345",
			role:        UserRoleAdmin,
			wantErr:     true,
			expectedErr: ErrInvalidPassword,
		},
		{
			name:        "invalid role",
			username:    "admin",
			password:    "password123",
			role:        UserRole("invalid"),
			wantErr:     true,
			expectedErr: ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.username, tt.password, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("NewUser() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Error("NewUser() returned nil user")
				return
			}

			// Verificar campos
			if user.Username != tt.username {
				t.Errorf("Username = %v, want %v", user.Username, tt.username)
			}

			if user.Role != tt.role {
				t.Errorf("Role = %v, want %v", user.Role, tt.role)
			}

			if !user.IsActive {
				t.Error("IsActive should be true for new users")
			}

			if user.PasswordHash == "" {
				t.Error("PasswordHash should not be empty")
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	user, _ := NewUser("testuser", "password123", UserRoleAdmin)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "wrong password",
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ValidatePassword(tt.password)

			if tt.wantErr && err == nil {
				t.Error("ValidatePassword() expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("ValidatePassword() unexpected error: %v", err)
			}
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	user, _ := NewUser("testuser", "initialpassword", UserRoleAdmin)

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "valid new password",
			password: "newpassword123",
			wantErr:  false,
		},
		{
			name:        "password too short",
			password:    "12345",
			wantErr:     true,
			expectedErr: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.SetPassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Error("SetPassword() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("SetPassword() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("SetPassword() unexpected error: %v", err)
				return
			}

			// Verificar se a senha foi alterada
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.password))
			if err != nil {
				t.Error("Password was not correctly hashed")
			}
		})
	}
}

func TestUser_RoleMethods(t *testing.T) {
	adminUser, _ := NewUser("admin", "admin123", UserRoleAdmin)
	managerUser, _ := NewUser("manager", "manager123", UserRoleManager)
	viewerUser, _ := NewUser("viewer", "viewer123", UserRoleViewer)

	// Test IsAdmin
	if !adminUser.IsAdmin() {
		t.Error("Admin user should return true for IsAdmin()")
	}
	if managerUser.IsAdmin() || viewerUser.IsAdmin() {
		t.Error("Non-admin users should return false for IsAdmin()")
	}

	// Test IsManager
	if !managerUser.IsManager() {
		t.Error("Manager user should return true for IsManager()")
	}

	// Test IsViewer
	if !viewerUser.IsViewer() {
		t.Error("Viewer user should return true for IsViewer()")
	}

	// Test CanManageUsers
	if !adminUser.CanManageUsers() {
		t.Error("Admin should be able to manage users")
	}
	if managerUser.CanManageUsers() || viewerUser.CanManageUsers() {
		t.Error("Only admin should be able to manage users")
	}

	// Test CanWrite
	if !adminUser.CanWrite() || !managerUser.CanWrite() {
		t.Error("Admin and Manager should have write permissions")
	}
	if viewerUser.CanWrite() {
		t.Error("Viewer should not have write permissions")
	}
}

func TestUser_ActivateDeactivate(t *testing.T) {
	user, _ := NewUser("testuser", "password123", UserRoleAdmin)

	if !user.IsActive {
		t.Error("New user should be active")
	}

	user.Deactivate()
	if user.IsActive {
		t.Error("User should be inactive after Deactivate()")
	}

	user.Activate()
	if !user.IsActive {
		t.Error("User should be active after Activate()")
	}
}

func TestUser_ChangeRole(t *testing.T) {
	user, _ := NewUser("testuser", "password123", UserRoleViewer)

	err := user.ChangeRole(UserRoleManager)
	if err != nil {
		t.Errorf("ChangeRole() unexpected error: %v", err)
	}

	if user.Role != UserRoleManager {
		t.Errorf("Role = %v, want %v", user.Role, UserRoleManager)
	}

	err = user.ChangeRole(UserRole("invalid"))
	if err != ErrInvalidRole {
		t.Errorf("ChangeRole() expected ErrInvalidRole, got %v", err)
	}
}
