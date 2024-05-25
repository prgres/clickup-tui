package listitem

type Item struct {
	title string
	desc  string
	data  interface{}
}

func NewItem(title, desc string, data interface{}) Item {
	return Item{
		title: title,
		desc:  desc,
		data:  data,
	}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }
func (i Item) Data() interface{}   { return i.data }
