package postgres

import (
	"database/sql"
	"time"

	"github.com/lucianoZgabriel/kitnet-manager/internal/domain"
)

// Helper functions para conversão de tipos nullable

// Para ponteiros (usado em Payment)
func toNullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromNullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func toNullStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func fromNullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func paymentMethodToStringPtr(pm *domain.PaymentMethod) *string {
	if pm == nil {
		return nil
	}
	s := string(*pm)
	return &s
}

func stringToPaymentMethodPtr(s *string) *domain.PaymentMethod {
	if s == nil {
		return nil
	}
	pm := domain.PaymentMethod(*s)
	return &pm
}
