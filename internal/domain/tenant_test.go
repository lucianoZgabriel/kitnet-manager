package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTenant(t *testing.T) {
	t.Run("should create valid tenant", func(t *testing.T) {
		tenant, err := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"11987654321",
			"joao@example.com",
		)

		require.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "João da Silva", tenant.FullName)
		assert.Equal(t, "123.456.789-00", tenant.CPF)
		assert.Equal(t, "11987654321", tenant.Phone)
		assert.Equal(t, "joao@example.com", tenant.Email)
		assert.NotEqual(t, "", tenant.ID.String())
	})

	t.Run("should create tenant without email", func(t *testing.T) {
		tenant, err := NewTenant(
			"Maria Santos",
			"987.654.321-00",
			"11912345678",
			"",
		)

		require.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "", tenant.Email)
	})

	t.Run("should fail with empty name", func(t *testing.T) {
		tenant, err := NewTenant(
			"",
			"123.456.789-00",
			"11987654321",
			"",
		)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrInvalidFullName, err)
	})

	t.Run("should fail with invalid CPF format", func(t *testing.T) {
		tenant, err := NewTenant(
			"João da Silva",
			"12345678900", // sem formatação
			"11987654321",
			"",
		)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrInvalidCPF, err)
	})

	t.Run("should fail with empty phone", func(t *testing.T) {
		tenant, err := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"",
			"",
		)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrInvalidPhone, err)
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		tenant, err := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"11987654321",
			"email-invalido",
		)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Equal(t, ErrInvalidEmail, err)
	})
}

func TestTenant_ValidateCPF(t *testing.T) {
	tests := []struct {
		name        string
		cpf         string
		shouldError bool
	}{
		{"valid CPF", "123.456.789-00", false},
		{"valid CPF 2", "987.654.321-99", false},
		{"invalid format - no dots", "12345678900", true},
		{"invalid format - no dash", "123.456.78900", true},
		{"invalid format - letters", "123.456.789-AB", true},
		{"invalid format - incomplete", "123.456.789-0", true},
		{"empty CPF", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant := &Tenant{CPF: tt.cpf}
			err := tenant.ValidateCPF()

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTenant_FormatPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"mobile with 11 digits", "11987654321", "(11) 98765-4321"},
		{"landline with 10 digits", "1133334444", "(11) 3333-4444"},
		{"already formatted mobile", "(11) 98765-4321", "(11) 98765-4321"},
		{"with spaces and dashes", "11 9 8765-4321", "(11) 98765-4321"},
		{"invalid format", "123", "123"}, // retorna original se inválido
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant := &Tenant{Phone: tt.phone}
			result := tenant.FormatPhone()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTenant_UpdateInfo(t *testing.T) {
	t.Run("should update tenant info", func(t *testing.T) {
		tenant, _ := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"11987654321",
			"joao@example.com",
		)

		err := tenant.UpdateInfo(
			"João Silva Junior",
			"11999887766",
			"joao.junior@example.com",
			"RG",
			"12345678",
		)

		require.NoError(t, err)
		assert.Equal(t, "João Silva Junior", tenant.FullName)
		assert.Equal(t, "11999887766", tenant.Phone)
		assert.Equal(t, "joao.junior@example.com", tenant.Email)
		assert.Equal(t, "RG", tenant.IDDocumentType)
		assert.Equal(t, "12345678", tenant.IDDocumentNumber)
	})

	t.Run("should fail with invalid email on update", func(t *testing.T) {
		tenant, _ := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"11987654321",
			"",
		)

		err := tenant.UpdateInfo("", "", "email-invalido", "", "")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidEmail, err)
	})

	t.Run("should not update if field is empty", func(t *testing.T) {
		tenant, _ := NewTenant(
			"João da Silva",
			"123.456.789-00",
			"11987654321",
			"joao@example.com",
		)

		originalName := tenant.FullName
		err := tenant.UpdateInfo("", "", "", "", "")

		require.NoError(t, err)
		assert.Equal(t, originalName, tenant.FullName) // mantém o original
	})
}

func TestTenant_String(t *testing.T) {
	tenant, _ := NewTenant(
		"João da Silva",
		"123.456.789-00",
		"11987654321",
		"",
	)

	result := tenant.String()

	assert.Equal(t, "João da Silva (CPF: 123.456.789-00)", result)
}
