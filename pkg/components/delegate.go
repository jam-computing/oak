package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jam-computing/oak/pkg/tcp"
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.play):
				packet := tcp.NewFullPacket(tcp.NewMetaPacket(), &title, nil)
				packet.Meta.Status = 200
				packet.Meta.Command = 2

				_ = packet.SendRecv()

				return m.NewStatusMessage(statusMessageStyle("Playing: " + title))
			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + title))

			case key.Matches(msg, keys.blank):
				packet := tcp.NewFullPacket(tcp.NewMetaPacket(), nil, nil)
				packet.Meta.Status = 200
				packet.Meta.Command = 13
				_ = packet.SendRecv()
				return m.NewStatusMessage(statusMessageStyle("Stopped all animations"))
			}
		}

		return nil
	}

	help := []key.Binding{keys.play, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	play   key.Binding
	remove key.Binding
	blank  key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.play,
		d.remove,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.play,
			d.remove,
			d.blank,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		play: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
		blank: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "blank tree"),
		),
	}
}
