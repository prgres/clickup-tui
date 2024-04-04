package folderslist

import (
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

func folderListToItems(folders []clickup.Folder) []listitem.Item {
	items := make([]listitem.Item, len(folders))
	for i, folder := range folders {
		items[i] = folderToItem(folder)
	}
	return items
}

func folderToItem(folder clickup.Folder) listitem.Item {
	return listitem.NewItem(
		folder.Name,
		folder.Id,
	)
}
