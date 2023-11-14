package listitem

import "github.com/charmbracelet/bubbles/list"

type Item struct {
	title string
	desc  string
}

func NewItem(title, desc string) Item {
	return Item{
		title: title,
		desc:  desc,
	}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func ItemListToBubblesItems(items []Item) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = ItemToBubblesItem(item)
	}
	return listItems
}

func ItemToBubblesItem(item Item) list.Item {
	return list.Item(item)
}

func BubblesItemToItem(item list.Item) Item {
	return item.(Item)
}
