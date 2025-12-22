package model

import (
	"strings"
	"testing"
	"time"
)

func TestModuleItem_Title(t *testing.T) {
	item := ModuleItem{
		Module: Module{
			ModuleName: "Test Crypto Module",
		},
	}

	got := item.Title()
	want := "Test Crypto Module"

	if got != want {
		t.Errorf("Title() = %q, want %q", got, want)
	}
}

func TestModuleItem_Description(t *testing.T) {
	tests := []struct {
		name     string
		item     ModuleItem
		contains []string
	}{
		{
			name: "active module",
			item: ModuleItem{
				Module: Module{
					VendorName: "Acme Corp",
					ModuleType: "Hardware",
					Status:     StatusActive,
				},
			},
			contains: []string{"Acme Corp", "Hardware", "Active"},
		},
		{
			name: "historical module",
			item: ModuleItem{
				Module: Module{
					VendorName: "Old Vendor",
					ModuleType: "Software",
					Status:     StatusHistorical,
				},
			},
			contains: []string{"Old Vendor", "Software", "Historical"},
		},
		{
			name: "in process module",
			item: ModuleItem{
				Module: Module{
					VendorName: "New Vendor",
					ModuleType: "Firmware",
					Status:     StatusInProcess,
				},
			},
			contains: []string{"New Vendor", "Firmware", "In Process"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.Description()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Description() = %q, should contain %q", got, want)
				}
			}
		})
	}
}

func TestModuleItem_FilterValue(t *testing.T) {
	item := ModuleItem{
		Module: Module{
			CertificateNumber: "1234",
			ModuleName:        "OpenSSL FIPS Provider",
			VendorName:        "OpenSSL Foundation",
			ModuleType:        "Software",
			ValidationDate:    time.Now(),
		},
	}

	got := item.FilterValue()

	// Should contain searchable fields
	if !strings.Contains(got, "1234") {
		t.Error("FilterValue() should contain certificate number")
	}
	if !strings.Contains(got, "OpenSSL FIPS Provider") {
		t.Error("FilterValue() should contain module name")
	}
	if !strings.Contains(got, "OpenSSL Foundation") {
		t.Error("FilterValue() should contain vendor name")
	}
}

func TestModuleItem_FilterValue_EmptyFields(t *testing.T) {
	item := ModuleItem{
		Module: Module{
			ModuleName: "Test Module",
		},
	}

	// Should not panic with empty fields
	got := item.FilterValue()
	if !strings.Contains(got, "Test Module") {
		t.Error("FilterValue() should contain module name even with empty fields")
	}
}
