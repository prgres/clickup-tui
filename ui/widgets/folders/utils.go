package folders

import "github.com/prgrs/clickup/pkg/clickup"

const (
	SPACE_SRE = "48458830"
)

func folderListToItems(folders []clickup.Folder) []item {
	items := make([]item, len(folders))
	for i, folder := range folders {
		items[i] = folderToItem(folder)
	}
	return items
}

func folderToItem(folder clickup.Folder) item {
	return item{
		folder.Name,
		folder.Id,
	}
}
