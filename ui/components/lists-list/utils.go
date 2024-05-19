package listslist

import (
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

func listsListToItems(lists []clickup.List) []listitem.Item {
	items := make([]listitem.Item, len(lists))
	for i, list := range lists {
		items[i] = listToItem(list)
	}
	return items
}

func listToItem(list clickup.List) listitem.Item {
	return listitem.NewItem(
		list.Name,
		list.Id,
	)
}
