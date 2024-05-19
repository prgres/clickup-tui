package spaceslist

import (
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

func spaceListToItems(spaces []clickup.Space) []listitem.Item {
	items := make([]listitem.Item, len(spaces))
	for i, space := range spaces {
		items[i] = spaceToItem(space)
	}
	return items
}

func spaceToItem(space clickup.Space) listitem.Item {
	return listitem.NewItem(
		space.Name,
		space.Id,
	)
}
