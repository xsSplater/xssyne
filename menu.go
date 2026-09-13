// xssyne/menu.go

package fyne

type systemTrayDriver interface {
	Driver
	SetSystemTrayMenu(*Menu)
	SystemTrayMenu() *Menu
}

// Menu stores the information required for a standard menu.
// A menu can pop down from a [MainMenu] or could be a pop out menu.
type Menu struct {
	Label string
	Items []*MenuItem
	// Icon — опциональная иконка слева от Label в верхнем меню.
	//
	// Since: xssyne
	Icon Resource
	// IsSeparator — если true, Menu не является кликабельным пунктом,
	// а служит вертикальным разделителем в MenuBar. Label, Items, Icon
	// игнорируются. Работает только на desktop (Windows/Linux);
	// на macOS меню рисует AppKit — там сепаратор пропускается при
	// построении структуры.
	//
	// Since: xssyne
	IsSeparator bool
}

// NewMenu creates a new menu given the specified label (to show in a [MainMenu]) and list of items to display.
func NewMenu(label string, items ...*MenuItem) *Menu {
	return &Menu{Label: label, Items: items}
}

// NewMenuWithIcon — как NewMenu, но с иконкой (для верхнего меню).
//
// Since: xssyne
func NewMenuWithIcon(label string, icon Resource, items ...*MenuItem) *Menu {
	return &Menu{Label: label, Items: items, Icon: icon}
}

// NewMenuSeparator создаёт пункт-разделитель для верхнего меню.
//
// Since: xssyne
func NewMenuSeparator() *Menu {
	return &Menu{IsSeparator: true}
}

// Refresh will instruct this menu to update its display.
//
// Since: 2.2
func (m *Menu) Refresh() {
	for _, w := range CurrentApp().Driver().AllWindows() {
		main := w.MainMenu()
		if main != nil {
			for _, menu := range main.Items {
				if menu == m {
					w.SetMainMenu(main)
					break
				}
			}
		}
	}

	if d, ok := CurrentApp().Driver().(systemTrayDriver); ok {
		if m == d.SystemTrayMenu() {
			d.SetSystemTrayMenu(m)
		}
	}
}

// MenuItem is a single item within any menu, it contains a display Label and Action function that is called when tapped.
type MenuItem struct {
	ChildMenu *Menu
	// Since: 2.1
	IsQuit      bool
	IsSeparator bool
	Label       string
	Action      func() `json:"-"`
	// Since: 2.1
	Disabled bool
	// Since: 2.1
	Checked bool
	// Since: 2.2
	Icon Resource
	// Since: 2.2
	Shortcut Shortcut

	// Danger — пункт отображается «опасным» цветом (красный) для действий
	// типа «Удалить», «Выйти без сохранения». Не влияет на поведение,
	// только на цвет текста.
	//
	// Since: xssyne
	Danger bool

	// Header — пункт-заголовок группы. Не кликается, выделен жирным,
	// имеет фон. Используется внутри выпадающего меню для визуальной
	// группировки пунктов.
	//
	// Since: xssyne
	Header bool

	// Subtitle — вторая строка под основным текстом, мелким шрифтом.
	// Если задана, высота пункта увеличивается. Header-пункты Subtitle
	// игнорируют.
	//
	// Since: xssyne
	Subtitle string
	// CheckedIcon — кастомная иконка для отмеченного пункта. Если nil,
	// используется стандартная (IconNameConfirm).
	//
	// Since: xssyne
	CheckedIcon Resource

	// UncheckedIcon — кастомная иконка для неотмеченного пункта. Если nil,
	// ничего не рисуется (поведение Fyne по умолчанию).
	//
	// Since: xssyne
	UncheckedIcon Resource
}

// NewMenuItem creates a new menu item from the passed label and action parameters.
func NewMenuItem(label string, action func()) *MenuItem {
	return &MenuItem{Label: label, Action: action}
}

// NewMenuItemWithIcon creates a new menu item from the passed label, icon, and action parameters.
//
// Since: 2.7
func NewMenuItemWithIcon(label string, icon Resource, action func()) *MenuItem {
	return &MenuItem{Label: label, Icon: icon, Action: action}
}

// NewMenuItemSeparator creates a menu item that is to be used as a separator.
func NewMenuItemSeparator() *MenuItem {
	return &MenuItem{IsSeparator: true, Action: func() {}}
}

// MainMenu defines the data required to show a menu bar (desktop) or other appropriate top level menu.
type MainMenu struct {
	Items []*Menu
}

// NewMainMenu creates a top level menu structure used by fyne.Window for displaying a menubar
// (or appropriate equivalent).
func NewMainMenu(items ...*Menu) *MainMenu {
	return &MainMenu{Items: items}
}

// Refresh will instruct any rendered menus using this struct to update their display.
//
// Since: 2.2
func (m *MainMenu) Refresh() {
	for _, w := range CurrentApp().Driver().AllWindows() {
		menu := w.MainMenu()
		if menu != nil && menu == m {
			w.SetMainMenu(m)
		}
	}
}
