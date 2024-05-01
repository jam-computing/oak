package components

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jam-computing/oak/pkg/tcp"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func newItem(title, desc string) item {
    return item{
        title: title,
        desc: desc,
    }
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type SelectorModel struct {
	list list.Model
}

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectorModel) View() string {
	return docStyle.Render(m.list.View())
}

func GetSelectorModel() SelectorModel {
	packet := tcp.NewPacket()
	packet.Status = 200
	packet.Command = 4
	recv := packet.SendRecv()

	jsonData := recv.Data[:recv.Len]
	var animations []tcp.Animation

	var items []list.Item

	err := json.Unmarshal([]byte(jsonData), &animations)

	if err != nil {
		items = []list.Item{}
	} else {
        for _, animation := range animations {
            items = append(items, newItem(animation.Title, animation.Artist))
        }
    }

	m := SelectorModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Animations"

	return m
}
