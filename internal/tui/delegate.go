package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanolivertroy/cmvp-tui/internal/model"
)

// ModuleDelegate is a custom delegate for rendering module items in the list
type ModuleDelegate struct {
	ShowDescription bool
	Styles          ModuleDelegateStyles
}

// ModuleDelegateStyles contains styles for the module delegate
type ModuleDelegateStyles struct {
	NormalTitle   lipgloss.Style
	NormalDesc    lipgloss.Style
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
	DimmedTitle   lipgloss.Style
	DimmedDesc    lipgloss.Style
}

// NewModuleDelegate creates a new module delegate with default styles
func NewModuleDelegate() ModuleDelegate {
	return ModuleDelegate{
		ShowDescription: true,
		Styles: ModuleDelegateStyles{
			NormalTitle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Padding(0, 0, 0, 2),
			NormalDesc: lipgloss.NewStyle().
				Foreground(SubtleColor).
				Padding(0, 0, 0, 2),
			SelectedTitle: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(PrimaryColor).
				Foreground(PrimaryColor).
				Padding(0, 0, 0, 1),
			SelectedDesc: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(PrimaryColor).
				Foreground(lipgloss.Color("#FAFAFA")).
				Padding(0, 0, 0, 1),
			DimmedTitle: lipgloss.NewStyle().
				Foreground(SubtleColor).
				Padding(0, 0, 0, 2),
			DimmedDesc: lipgloss.NewStyle().
				Foreground(SubtleColor).
				Padding(0, 0, 0, 2),
		},
	}
}

// Height returns the height of each item
func (d ModuleDelegate) Height() int {
	if d.ShowDescription {
		return 2
	}
	return 1
}

// Spacing returns the spacing between items
func (d ModuleDelegate) Spacing() int {
	return 1
}

// Update handles any updates to the delegate
func (d ModuleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render renders a single item
func (d ModuleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	moduleItem, ok := item.(model.ModuleItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	isFiltering := m.FilterState() == list.Filtering

	var titleStyle, descStyle lipgloss.Style
	if isSelected && !isFiltering {
		titleStyle = d.Styles.SelectedTitle
		descStyle = d.Styles.SelectedDesc
	} else if isFiltering && !isSelected {
		titleStyle = d.Styles.DimmedTitle
		descStyle = d.Styles.DimmedDesc
	} else {
		titleStyle = d.Styles.NormalTitle
		descStyle = d.Styles.NormalDesc
	}

	// Render title with certificate number prefix
	var title string
	if moduleItem.CertificateNumber != "" {
		title = fmt.Sprintf("[%s] %s", moduleItem.CertificateNumber, moduleItem.Title())
	} else {
		title = moduleItem.Title()
	}
	fmt.Fprint(w, titleStyle.Render(title))

	if d.ShowDescription {
		fmt.Fprint(w, "\n")
		desc := moduleItem.Description()
		fmt.Fprint(w, descStyle.Render(desc))
	}
}
