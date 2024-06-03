package navigator

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	switch m.state {
	case m.componentWorkspacesList.Id():
		return m.componentWorkspacesList.Help()
	case m.componentSpacesList.Id():
		return m.componentSpacesList.Help()
	case m.componentFoldersList.Id():
		return m.componentFoldersList.Help()
	case m.componentListsList.Id():
		return m.componentListsList.Help()
	default:
		return common.NewEmptyHelp()
	}
}
