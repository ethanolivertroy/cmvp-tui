package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanolivertroy/nist-cmvp-cli/internal/api"
	"github.com/ethanolivertroy/nist-cmvp-cli/internal/model"
)

// ViewState represents the current view of the application
type ViewState int

const (
	ViewList ViewState = iota
	ViewDetail
)

// ModulesLoadedMsg is sent when modules are loaded from the API
type ModulesLoadedMsg struct {
	Modules []list.Item
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Err error
}

// Model is the main application model
type Model struct {
	list           list.Model
	allModules     []list.Item
	spinner        spinner.Model
	loading        bool
	err            error
	width          int
	height         int
	view           ViewState
	selectedModule *model.ModuleItem
	apiClient      *api.Client
}

// NewModel creates a new application model
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(PrimaryColor)

	return Model{
		spinner:   s,
		loading:   true,
		view:      ViewList,
		apiClient: api.NewClient(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchModules(),
	)
}

func (m Model) fetchModules() tea.Cmd {
	return func() tea.Msg {
		modules, err := m.apiClient.FetchAllModules()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		items := make([]list.Item, len(modules))
		for i, mod := range modules {
			items[i] = model.ModuleItem{Module: mod}
		}
		return ModulesLoadedMsg{Modules: items}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle keys while filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == ViewDetail {
				m.view = ViewList
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.view == ViewList && !m.loading {
				if item, ok := m.list.SelectedItem().(model.ModuleItem); ok {
					m.selectedModule = &item
					m.view = ViewDetail
					return m, nil
				}
			}
		case "esc", "backspace":
			if m.view == ViewDetail {
				m.view = ViewList
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.loading {
			m.list.SetSize(msg.Width-4, msg.Height-4)
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case ModulesLoadedMsg:
		m.loading = false
		m.allModules = msg.Modules

		delegate := NewModuleDelegate()
		m.list = list.New(msg.Modules, delegate, m.width-4, m.height-4)
		m.list.Title = "NIST CMVP Modules"
		m.list.SetShowStatusBar(true)
		m.list.SetFilteringEnabled(true)
		m.list.Styles.Title = TitleStyle
		m.list.FilterInput.Prompt = "Filter: "
		m.list.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(PrimaryColor)

		return m, nil

	case ErrorMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil
	}

	// Pass messages to list component when in list view
	if m.view == ViewList && !m.loading && m.err == nil {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	if m.loading {
		return AppStyle.Render(
			fmt.Sprintf("\n\n   %s Loading CMVP modules...\n\n", m.spinner.View()),
		)
	}

	if m.err != nil {
		return AppStyle.Render(
			lipgloss.NewStyle().Foreground(ErrorColor).Render(
				fmt.Sprintf("\n\n   Error: %v\n\n   Press q to quit.", m.err),
			),
		)
	}

	switch m.view {
	case ViewDetail:
		return m.renderDetailView()
	default:
		return AppStyle.Render(m.list.View())
	}
}

func (m Model) renderDetailView() string {
	if m.selectedModule == nil {
		return ""
	}

	mod := m.selectedModule

	var b strings.Builder

	// Title with status badge
	b.WriteString(DetailTitleStyle.Render(mod.ModuleName))
	b.WriteString("  ")
	b.WriteString(StatusBadge(mod.Status))
	b.WriteString("\n\n")

	// Details grid
	details := []struct {
		label string
		value string
	}{
		{"Certificate #:", mod.CertificateNumber},
		{"Vendor:", mod.VendorName},
		{"Module Type:", mod.ModuleType},
	}

	// Add validation date if available
	if !mod.ValidationDate.IsZero() {
		details = append(details, struct {
			label string
			value string
		}{"Validation Date:", mod.ValidationDate.Format("January 2, 2006")})
	}

	// Add URL if available
	if mod.CertificateURL != "" {
		details = append(details, struct {
			label string
			value string
		}{"NIST URL:", mod.CertificateURL})
	}

	for _, d := range details {
		if d.value == "" {
			continue
		}
		b.WriteString(DetailLabelStyle.Render(d.label))
		if d.label == "NIST URL:" {
			b.WriteString(DetailURLStyle.Render(d.value))
		} else {
			b.WriteString(DetailValueStyle.Render(d.value))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("Press ESC or Backspace to return to list"))

	return AppStyle.Render(b.String())
}
