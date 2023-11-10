package viewtabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx                *context.UserContext
	views              map[string][]clickup.View
	SelectedView       string
	SelectedViewStruct clickup.View
	SelectedFolder     string
	// SelectedSpace      string
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:          ctx,
		views:        map[string][]clickup.View{},
		SelectedView: SPACE_SRE_LIST_COOL,
		// SelectedSpace: SPACE_SRE,
	}
}
func nextView(views []clickup.View, SelectedView string) string {
	for i, view := range views {
		if view.Id == SelectedView {
			if i+1 < len(views) {
				return views[i+1].Id
			}
			return views[0].Id
		}
	}
	return views[0].Id
}

func prevView(views []clickup.View, SelectedView string) string {
	for i, view := range views {
		if view.Id == SelectedView {
			if i-1 >= 0 {
				return views[i-1].Id
			}
			return views[len(views)-1].Id
		}
	}
	return views[0].Id
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "h", "left":
			m.SelectedView = prevView(m.views[m.SelectedFolder], m.SelectedView)
			return m, ViewChangedCmd(m.SelectedView)
		case "l", "right":
			m.SelectedView = nextView(m.views[m.SelectedFolder], m.SelectedView)
			return m, ViewChangedCmd(m.SelectedView)

		default:
			return m, nil
		}

	case common.FolderChangeMsg:
		m.ctx.Logger.Infof("ViewsView received FolderChangeMsg: %s", string(msg))
		m.SelectedFolder = string(msg)

		views, err := m.getViewsFromFolder(string(msg))
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m.ctx.Logger.Infof("----GOT %d VIEWS FORM FOLDER: %s", len(views), string(msg))

		if len(views) == 0 {
			return m, tea.Batch(
				ViewChangedCmd(m.SelectedView),
			)
		}

		m.views[m.SelectedFolder] = views
		m.SelectedView = m.views[m.SelectedFolder][0].Id
		viewsIds := []string{}
		for _, view := range views {
			m.ctx.Logger.Infof("----VIEW: %s; %s", view.Name, view.Id)
			if view.Id == m.SelectedView {
				continue
			}
			viewsIds = append(viewsIds, view.Id)
		}

		return m, tea.Batch(
			ViewChangedCmd(m.SelectedView),
			FetchViewsCmd(viewsToIdList(views)),
		)

	// case common.SpaceChangeMsg:
	// 	m.ctx.Logger.Infof("ViewsView received SpaceChangedMsg")
	// 	m.SelectedSpace = string(msg)
	// 	views, err := m.getViews(string(msg))
	// 	if err != nil {
	// 		return m, common.ErrCmd(err)
	// 	}
	// 	if len(views) == 0 {
	// 		return m, tea.Batch(
	// 			ViewChangedCmd(m.SelectedView),
	// 		)
	// 	}

	// 	m.views[m.SelectedSpace] = views
	// 	m.SelectedView = m.views[m.SelectedSpace][0].Id
	// 	viewsIds := []string{}
	// 	for _, view := range views {
	// 		if view.Id == m.SelectedView {
	// 			continue
	// 		}
	// 		viewsIds = append(viewsIds, view.Id)
	// 	}

	// 	return m, tea.Batch(
	// 		ViewChangedCmd(m.SelectedView),
	// 		FetchViewsCmd(viewsToIdList(views)),
	// 	)

	case common.FocusMsg:
		m.ctx.Logger.Info("ViewsView received FocusMsg")
		return m, nil

	case ViewsListLoadedMsg:
		m.ctx.Logger.Info("ViewsView received ViewsListLoadedMsg")
		m.views[m.SelectedView] = []clickup.View(msg)
		// m.views[m.SelectedSpace] = []clickup.View(msg)
		return m, nil

	case ViewChangedMsg:
		m.ctx.Logger.Info("ViewsView received ViewChangedMsg")
		viewName := string(msg)
		allViews := m.views[m.SelectedFolder]
		// allViews := m.views[m.SelectedSpace]
		for _, view := range allViews {
			if view.Id == viewName {
				m.SelectedViewStruct = view
				return m, tea.Batch(common.ViewLoadedCmd(view))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func removeView(views []clickup.View, s int) []clickup.View {
	return append(views[:s], views[s+1:]...)
}

func viewsToIdList(views []clickup.View) []string {
	ids := []string{}
	for _, view := range views {
		ids = append(ids, view.Id)
	}
	return ids
}

func (m Model) View() string {
	viewsNames := []string{}
	for _, view := range m.views[m.SelectedFolder] {
		// for _, view := range m.views[m.SelectedSpace] {
		t := ""
		if m.SelectedView == view.Id {
			t = activeTabStyle.Render(view.Name)
		} else {
			t = inactiveTabStyle.Render(view.Name)
		}
		viewsNames = append(viewsNames, t)
	}

	return strings.Join(viewsNames, " | ")
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing component: TabsView")
	return common.FolderChangeCmd(FOLDER_INITIATIVE)
}

func (m Model) getViewsFromSpaceCmd(space string) tea.Cmd {
	return func() tea.Msg {
		views, err := m.getViewsFromSpace(space)
		if err != nil {
			return common.ErrMsg(err)
		}
		return ViewsListLoadedMsg(views)
	}
}

func (m Model) getViewsFromFolderCmd(folder string) tea.Cmd {
	return func() tea.Msg {
		views, err := m.getViewsFromFolder(folder)
		if err != nil {
			return common.ErrMsg(err)
		}
		return ViewsListLoadedMsg(views)
	}
}

func (m Model) getViewsFromFolder(folder string) ([]clickup.View, error) {
	m.ctx.Logger.Infof("Getting views for folder: %s", folder)

	// data, ok := m.ctx.Cache.Get("views", folder)
	// if ok {
	// 	m.ctx.Logger.Infof("Views found in cache")
	// 	var views []clickup.View
	// 	if err := m.ctx.Cache.ParseData(data, &views); err != nil {
	// 		return nil, err
	// 	}

	// 	return views, nil
	// }
	m.ctx.Logger.Info("Views not found in cache")

	m.ctx.Logger.Info("Fetching views from API")
	client := m.ctx.Clickup

	m.ctx.Logger.Infof("Getting views from folder: %s", folder)
	views, err := client.GetViewsFromFolder(folder)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d views in folder %s", len(views), folder)

	m.ctx.Logger.Info("Caching views")
	m.ctx.Cache.Set("views", folder, views)

	return views, nil
}

func (m Model) getViewsFromSpace(space string) ([]clickup.View, error) {
	m.ctx.Logger.Infof("Getting views for space: %s", space)

	data, ok := m.ctx.Cache.Get("views", space)
	if ok {
		m.ctx.Logger.Infof("Views found in cache")
		var views []clickup.View
		if err := m.ctx.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}
	m.ctx.Logger.Info("Views not found in cache")

	m.ctx.Logger.Info("Fetching views from API")
	client := m.ctx.Clickup

	m.ctx.Logger.Infof("Getting views from space: %s", space)
	views, err := client.GetViewsFromSpace(space)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d views in space %s", len(views), space)

	m.ctx.Logger.Info("Caching views")
	m.ctx.Cache.Set("views", space, views)

	return views, nil
}
