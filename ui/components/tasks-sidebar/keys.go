package taskssidebar

import "github.com/charmbracelet/bubbles/viewport"

type KeyMap struct {
	viewport.KeyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		KeyMap: viewport.DefaultKeyMap(),
	}
}
