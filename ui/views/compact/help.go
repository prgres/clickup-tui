package compact

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	switch m.state {
	case m.widgetNavigator.Id():
		return m.widgetNavigator.Help()
	case m.widgetViewsTabs.Id():
		return m.widgetViewsTabs.Help()
	case m.widgetTasks.Id():
		return m.widgetTasks.Help()
	default:
		return common.NewEmptyHelp()
	}
}
