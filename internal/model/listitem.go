package model

import (
	"fmt"
	"strings"
)

// ModuleItem wraps Module to implement list.DefaultItem interface
type ModuleItem struct {
	Module
}

// Title returns the primary display text (Module Name)
func (m ModuleItem) Title() string {
	return m.ModuleName
}

// Description returns secondary display text (Vendor + Type + Status)
func (m ModuleItem) Description() string {
	return fmt.Sprintf("%s | %s | %s", m.VendorName, m.ModuleType, m.Status.String())
}

// FilterValue returns the string used for fuzzy filtering
// Combines multiple fields for comprehensive search
func (m ModuleItem) FilterValue() string {
	return strings.Join([]string{
		m.ModuleName,
		m.VendorName,
		m.CertificateNumber,
		m.ModuleType,
	}, " ")
}
