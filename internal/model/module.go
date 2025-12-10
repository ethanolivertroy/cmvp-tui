package model

import "time"

// ModuleStatus represents the validation status of a module
type ModuleStatus int

const (
	StatusActive ModuleStatus = iota
	StatusHistorical
	StatusInProcess
)

func (s ModuleStatus) String() string {
	switch s {
	case StatusActive:
		return "Active"
	case StatusHistorical:
		return "Historical"
	case StatusInProcess:
		return "In Process"
	default:
		return "Unknown"
	}
}

// Module represents a NIST CMVP cryptographic module
type Module struct {
	CertificateNumber string
	CertificateURL    string
	VendorName        string
	ModuleName        string
	ModuleType        string
	ValidationDate    time.Time
	Status            ModuleStatus

	// Extended fields from certificate detail extraction
	Standard          string
	OverallLevel      int
	SunsetDate        string
	Caveat            string
	Embodiment        string
	Description       string
	Lab               string
	Algorithms        []string
	SecurityPolicyURL string
}
