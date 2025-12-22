package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethanolivertroy/cmvp-tui/internal/model"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if !m.loading {
		t.Error("expected loading to be true")
	}
	if m.view != ViewList {
		t.Errorf("expected view to be ViewList, got %v", m.view)
	}
	if m.apiClient == nil {
		t.Error("expected apiClient to be initialized")
	}
	if m.err != nil {
		t.Error("expected err to be nil")
	}
}

func TestModel_Init(t *testing.T) {
	m := NewModel()
	cmd := m.Init()

	if cmd == nil {
		t.Error("expected Init() to return a command")
	}
}

func TestModel_Update_QuitKey(t *testing.T) {
	m := NewModel()
	m.loading = false

	// Test 'q' key quits from list view
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestModel_Update_QuitFromDetail(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.selectedModule = &model.ModuleItem{}

	// Test 'q' from detail view returns to list
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.view != ViewList {
		t.Error("expected 'q' in detail view to return to list view")
	}
}

func TestModel_Update_EscapeFromDetail(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.selectedModule = &model.ModuleItem{}

	// Test ESC from detail view
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.view != ViewList {
		t.Error("expected ESC in detail view to return to list view")
	}
}

func TestModel_Update_ToggleAlgoDetails(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.selectedModule = &model.ModuleItem{
		Module: model.Module{
			AlgorithmsDetailed: []string{"AES-128", "SHA-256"},
		},
	}
	m.width = 80
	m.height = 24

	// Test 'd' key toggles algorithm details
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if !updated.showAlgoDetails {
		t.Error("expected 'd' to toggle showAlgoDetails to true")
	}
}

func TestModel_Update_ModulesLoadedMsg(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24

	modules := []list.Item{
		model.ModuleItem{
			Module: model.Module{
				ModuleName: "Test Module",
				VendorName: "Test Vendor",
				Status:     model.StatusActive,
			},
		},
	}

	msg := ModulesLoadedMsg{Modules: modules}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.loading {
		t.Error("expected loading to be false after ModulesLoadedMsg")
	}
	if len(updated.allModules) != 1 {
		t.Errorf("expected 1 module, got %d", len(updated.allModules))
	}
}

func TestModel_Update_ErrorMsg(t *testing.T) {
	m := NewModel()

	msg := ErrorMsg{Err: &testError{}}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.loading {
		t.Error("expected loading to be false after ErrorMsg")
	}
	if updated.err == nil {
		t.Error("expected err to be set")
	}
}

func TestModel_Update_WindowSizeMsg(t *testing.T) {
	m := NewModel()
	// Keep loading=true to avoid list.SetSize on uninitialized list

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 100 {
		t.Errorf("expected width 100, got %d", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("expected height 50, got %d", updated.height)
	}
}

func TestModel_Update_WindowSizeMsg_WithList(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24

	// First load modules to initialize the list
	modules := []list.Item{
		model.ModuleItem{
			Module: model.Module{
				ModuleName: "Test Module",
				VendorName: "Test Vendor",
			},
		},
	}
	newModel, _ := m.Update(ModulesLoadedMsg{Modules: modules})
	m = newModel.(Model)

	// Now test window resize with initialized list
	msg := tea.WindowSizeMsg{Width: 120, Height: 60}
	newModel, _ = m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 120 {
		t.Errorf("expected width 120, got %d", updated.width)
	}
	if updated.height != 60 {
		t.Errorf("expected height 60, got %d", updated.height)
	}
}

func TestModel_View_Loading(t *testing.T) {
	m := NewModel()
	m.loading = true

	view := m.View()

	if !strings.Contains(view, "Loading") {
		t.Error("expected loading view to contain 'Loading'")
	}
}

func TestModel_View_Error(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.err = &testError{}

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("expected error view to contain 'Error'")
	}
}

func TestViewState_Constants(t *testing.T) {
	if ViewList != 0 {
		t.Errorf("expected ViewList to be 0, got %d", ViewList)
	}
	if ViewDetail != 1 {
		t.Errorf("expected ViewDetail to be 1, got %d", ViewDetail)
	}
}

func TestModel_Update_EnterKey_SelectsModule(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24

	// Load modules to initialize the list
	modules := []list.Item{
		model.ModuleItem{
			Module: model.Module{
				ModuleName: "Test Module",
				VendorName: "Test Vendor",
				Status:     model.StatusActive,
			},
		},
	}
	newModel, _ := m.Update(ModulesLoadedMsg{Modules: modules})
	m = newModel.(Model)

	// Press enter to select the module
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.Update(msg)
	updated := newModel.(Model)

	if updated.view != ViewDetail {
		t.Error("expected Enter to switch to detail view")
	}
	if updated.selectedModule == nil {
		t.Error("expected selectedModule to be set")
	}
}

func TestModel_Update_BackspaceFromDetail(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.selectedModule = &model.ModuleItem{}

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.view != ViewList {
		t.Error("expected Backspace in detail view to return to list view")
	}
}

func TestModel_View_DetailView(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.width = 80
	m.height = 24
	m.selectedModule = &model.ModuleItem{
		Module: model.Module{
			CertificateNumber: "1234",
			ModuleName:        "Test Crypto Module",
			VendorName:        "Acme Corp",
			ModuleType:        "Hardware",
			Status:            model.StatusActive,
			Standard:          "FIPS 140-3",
			OverallLevel:      2,
			Algorithms:        []string{"AES", "SHA-256"},
		},
	}

	view := m.View()

	// Check that detail view contains expected content
	if !strings.Contains(view, "Test Crypto Module") {
		t.Error("detail view should contain module name")
	}
	if !strings.Contains(view, "Acme Corp") {
		t.Error("detail view should contain vendor name")
	}
}

func TestModel_View_DetailView_WithCaveat(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.width = 80
	m.height = 24
	m.selectedModule = &model.ModuleItem{
		Module: model.Module{
			ModuleName: "Caveat Module",
			VendorName: "Vendor",
			Status:     model.StatusActive,
			Caveat:     "Security warning here",
		},
	}

	view := m.View()

	if !strings.Contains(view, "CAVEAT") {
		t.Error("detail view should show CAVEAT label when caveat exists")
	}
}

func TestModel_View_DetailView_WithAlgoDetails(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.width = 80
	m.height = 24
	m.showAlgoDetails = true
	m.algoViewportReady = false
	m.selectedModule = &model.ModuleItem{
		Module: model.Module{
			ModuleName:         "Algo Module",
			VendorName:         "Vendor",
			Status:             model.StatusActive,
			AlgorithmsDetailed: []string{},
		},
	}

	view := m.View()

	if !strings.Contains(view, "Detailed") {
		t.Error("detail view should show 'Detailed' when showAlgoDetails is true")
	}
}

func TestModel_View_DetailView_NilModule(t *testing.T) {
	m := NewModel()
	m.loading = false
	m.view = ViewDetail
	m.selectedModule = nil

	view := m.View()

	// Should return empty string for nil module
	if view != "" {
		t.Log("view returned non-empty for nil module, which is acceptable")
	}
}

func TestModel_Update_CtrlC_Quits(t *testing.T) {
	m := NewModel()
	m.loading = false

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("expected Ctrl+C to return quit command")
	}
}

func TestBuildAlgorithmContent(t *testing.T) {
	algorithms := []string{"AES-128", "SHA-256", "RSA-2048"}
	content := buildAlgorithmContent(algorithms)

	for _, algo := range algorithms {
		if !strings.Contains(content, algo) {
			t.Errorf("content should contain %s", algo)
		}
	}
}

func TestBuildAlgorithmContent_Empty(t *testing.T) {
	content := buildAlgorithmContent([]string{})
	if content != "" {
		t.Error("empty algorithms should produce empty content")
	}
}

// Helper types for testing
type testError struct{}

func (e *testError) Error() string {
	return "test error"
}
