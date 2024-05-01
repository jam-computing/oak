package components

import tea "github.com/charmbracelet/bubbletea"

type AppState int

const (
	Loading AppState = iota
	ViewingAnimations
	ViewingAnimation
)

type Model struct {
	State    AppState
	Loader   LoaderModel
	Selector SelectorModel
}

func GetModel(state AppState) Model {
	return Model{
		State:    state,
		Loader:   GetLoaderModel(),
		Selector: GetSelectorModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Selector.list, cmd = m.Selector.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.State {
	case Loading:
		return "Loading"
	case ViewingAnimations:
		return m.View()
	default:
		return "Unknown State"
	}
}
