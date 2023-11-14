package folders

import "github.com/charmbracelet/bubbles/list"

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func itemListToItems(items []item) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = itemToListItem(item)
	}
	return listItems
}

func itemToListItem(item item) list.Item {
	return list.Item(item)
}
