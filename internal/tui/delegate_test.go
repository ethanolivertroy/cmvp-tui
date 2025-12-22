package tui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"

	"github.com/ethanolivertroy/cmvp-tui/internal/model"
)

func TestNewModuleDelegate(t *testing.T) {
	d := NewModuleDelegate()

	if !d.ShowDescription {
		t.Error("expected ShowDescription to be true by default")
	}
}

func TestModuleDelegate_Height(t *testing.T) {
	tests := []struct {
		name            string
		showDescription bool
		want            int
	}{
		{
			name:            "with description",
			showDescription: true,
			want:            2,
		},
		{
			name:            "without description",
			showDescription: false,
			want:            1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewModuleDelegate()
			d.ShowDescription = tt.showDescription
			if got := d.Height(); got != tt.want {
				t.Errorf("Height() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestModuleDelegate_Spacing(t *testing.T) {
	d := NewModuleDelegate()
	if got := d.Spacing(); got != 1 {
		t.Errorf("Spacing() = %d, want 1", got)
	}
}

func TestModuleDelegate_Update(t *testing.T) {
	d := NewModuleDelegate()
	cmd := d.Update(nil, nil)
	if cmd != nil {
		t.Error("Update() should return nil")
	}
}

func TestModuleDelegate_Render(t *testing.T) {
	d := NewModuleDelegate()

	item := model.ModuleItem{
		Module: model.Module{
			CertificateNumber: "1234",
			ModuleName:        "Test Module",
			VendorName:        "Test Vendor",
			ModuleType:        "Hardware",
			Status:            model.StatusActive,
		},
	}

	items := []list.Item{item}
	l := list.New(items, d, 80, 24)

	var buf bytes.Buffer
	d.Render(&buf, l, 0, item)

	output := buf.String()

	// Should contain certificate number and title
	if !strings.Contains(output, "1234") {
		t.Error("Render() should include certificate number")
	}
	if !strings.Contains(output, "Test Module") {
		t.Error("Render() should include module name")
	}
}

func TestModuleDelegate_Render_NoCertificate(t *testing.T) {
	d := NewModuleDelegate()

	item := model.ModuleItem{
		Module: model.Module{
			ModuleName: "In Process Module",
			VendorName: "Vendor",
			Status:     model.StatusInProcess,
		},
	}

	items := []list.Item{item}
	l := list.New(items, d, 80, 24)

	var buf bytes.Buffer
	d.Render(&buf, l, 0, item)

	output := buf.String()

	// Should contain title without certificate bracket
	if !strings.Contains(output, "In Process Module") {
		t.Error("Render() should include module name")
	}
}

func TestModuleDelegate_Render_InvalidItem(t *testing.T) {
	d := NewModuleDelegate()

	// Create a mock item that isn't a ModuleItem
	items := []list.Item{}
	l := list.New(items, d, 80, 24)

	var buf bytes.Buffer
	// Pass a string instead of ModuleItem - should not panic
	d.Render(&buf, l, 0, nil)

	// Should produce no output for invalid item
	if buf.Len() != 0 {
		t.Error("Render() should produce no output for nil item")
	}
}

func TestModuleDelegate_Render_WithoutDescription(t *testing.T) {
	d := NewModuleDelegate()
	d.ShowDescription = false

	item := model.ModuleItem{
		Module: model.Module{
			CertificateNumber: "5678",
			ModuleName:        "Short Module",
			VendorName:        "Vendor",
			Status:            model.StatusActive,
		},
	}

	items := []list.Item{item}
	l := list.New(items, d, 80, 24)

	var buf bytes.Buffer
	d.Render(&buf, l, 0, item)

	output := buf.String()

	// Should contain title
	if !strings.Contains(output, "Short Module") {
		t.Error("Render() should include module name")
	}
	// Should not contain newline (no description)
	if strings.Contains(output, "\n") {
		t.Error("Render() should not include newline when ShowDescription is false")
	}
}
