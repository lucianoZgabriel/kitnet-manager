package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Tenant representa um morador/inquilino
type Tenant struct {
	ID               uuid.UUID `json:"id"`
	FullName         string    `json:"full_name"`
	CPF              string    `json:"cpf"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email,omitempty"`
	IDDocumentType   string    `json:"id_document_type,omitempty"`
	IDDocumentNumber string    `json:"id_document_number,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Domain errors
var (
	ErrInvalidFullName  = errors.New("full name cannot be empty")
	ErrInvalidCPF       = errors.New("invalid CPF format")
	ErrInvalidCPFDigits = errors.New("CPF must contain exactly 11 digits")
	ErrInvalidPhone     = errors.New("phone cannot be empty")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrCPFAlreadyExists = errors.New("CPF already registered")
)

// CPF regex pattern: XXX.XXX.XXX-XX
var cpfRegex = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)

// Email regex pattern (simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewTenant cria um novo morador com valores padrão
func NewTenant(fullname, cpf, phone, email string) (*Tenant, error) {
	tenant := &Tenant{
		ID:        uuid.New(),
		FullName:  strings.TrimSpace(fullname),
		CPF:       cpf,
		Phone:     strings.TrimSpace(phone),
		Email:     strings.TrimSpace(email),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Valida o morador
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	return tenant, nil
}

// Validate verifica se o morador possui dados válidos
func (t *Tenant) Validate() error {
	// Valida nome completo
	if strings.TrimSpace(t.FullName) == "" {
		return ErrInvalidFullName
	}

	// Validar CPF
	if err := t.ValidateCPF(); err != nil {
		return err
	}

	// Validar telefone
	if strings.TrimSpace(t.Phone) == "" {
		return ErrInvalidPhone
	}

	// Validar email (se fornecido)
	if t.Email != "" {
		if !emailRegex.MatchString(t.Email) {
			return ErrInvalidEmail
		}
	}

	return nil
}

// ValidateCPF verifica se o CPF está no formato correto
func (t *Tenant) ValidateCPF() error {
	if !cpfRegex.MatchString(t.CPF) {
		return ErrInvalidCPF
	}

	// Extrair apenas os dígitos para validação adicional
	digits := strings.ReplaceAll(strings.ReplaceAll(t.CPF, ".", ""), "-", "")
	if len(digits) != 11 {
		return ErrInvalidCPFDigits
	}

	return nil
}

// FormatPhone formata o número de telefone
// Aceita vários formatos e retorna no formato (XX) XXXXX-XXXX ou (XX) XXXX-XXXX
func (t *Tenant) FormatPhone() string {
	// Remove tudo que não é dígito
	phone := regexp.MustCompile(`\D`).ReplaceAllString(t.Phone, "")

	// Verifica o tamanho
	if len(phone) == 11 { // Celular: (XX) XXXXX-XXXX
		return "(" + phone[0:2] + ") " + phone[2:7] + "-" + phone[7:11]
	} else if len(phone) == 10 { // Fixo: (XX) XXXX-XXXX
		return "(" + phone[0:2] + ") " + phone[2:6] + "-" + phone[6:10]
	}

	// Se não for um formato válido, retorna o original
	return t.Phone
}

// UpdateInfo atualiza as informações do morador
func (t *Tenant) UpdateInfo(fullName, phone, email, idDocType, idDocNumber string) error {
	if fullName != "" {
		t.FullName = strings.TrimSpace(fullName)
	}
	if phone != "" {
		t.Phone = strings.TrimSpace(phone)
	}
	if email != "" {
		t.Email = strings.TrimSpace(email)
	}
	if idDocType != "" {
		t.IDDocumentType = strings.TrimSpace(idDocType)
	}
	if idDocNumber != "" {
		t.IDDocumentNumber = strings.TrimSpace(idDocNumber)
	}

	t.UpdatedAt = time.Now()

	// Valida após atualização
	return t.Validate()
}

// HasActiveContract indica se o morador pode ser removido
// Esta validação será feita no Service com consulta ao repositório
func (t *Tenant) HasActiveContract() bool {
	// Placeholder - será implementado no service
	return false
}

// String retorna uma representação em string do morador
func (t *Tenant) String() string {
	return t.FullName + " (CPF: " + t.CPF + ")"
}
