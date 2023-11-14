package spaces

import (
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

const (
	TEAM_RAMP_NETWORK = "24301226"
	SPACE_SRE         = "48458830"
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
