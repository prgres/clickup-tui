package workspaceslist

import (
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

func workspaceListToItems(workspaces []clickup.Team) []listitem.Item {
	items := make([]listitem.Item, len(workspaces))
	for i, workspace := range workspaces {
		items[i] = workspaceToItem(workspace)
	}
	return items
}

func workspaceToItem(workspace clickup.Team) listitem.Item {
	return listitem.NewItem(
		workspace.Name,
		workspace.Id,
	)
}
