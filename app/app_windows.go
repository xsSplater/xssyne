//go:build windows

// xssyne/app/app_windows.go

package app

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/scheduler"
)

// xssyne
const notificationTemplate = `$title = %q
$content = %q
$iconPath = %q
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$toastXml.GetElementsByTagName("image")[0].SetAttribute("src", $iconPath) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%q).Show($toast);`

// Вариант без иконки: использует ToastText02, чтобы не оставлять
// в шаблоне <image src="file:///">, который некоторые версии Windows
// не могут распарсить.
const notificationNoIconTemplate = `$title = %q
$content = %q
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%q).Show($toast);`

const scheduledNotificationTemplate = `$title = %q
$content = %q
$iconPath = %q
$id = %q
$delivery = [DateTimeOffset]::Parse(%q)
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$toastXml.GetElementsByTagName("image")[0].SetAttribute("src", $iconPath) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$scheduled = [Windows.UI.Notifications.ScheduledToastNotification]::new($xml, $delivery)
$scheduled.Tag = $id
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%q).AddToSchedule($scheduled);`

// xssyne
const scheduledNotificationNoIconTemplate = `$title = %q
$content = %q
$id = %q
$delivery = [DateTimeOffset]::Parse(%q)
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$scheduled = [Windows.UI.Notifications.ScheduledToastNotification]::new($xml, $delivery)
$scheduled.Tag = $id
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%q).AddToSchedule($scheduled);`

const cancelScheduledNotificationTemplate = `$id = %q
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%q)
foreach ($s in $notifier.GetScheduledToastNotifications()) {
    if ($s.Tag -eq $id) { $notifier.RemoveFromSchedule($s) }
}`

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// xssyne
var scriptNum atomic.Int64

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	title := n.Title
	content := n.Content
	iconFilePath := a.cachedIconPath()
	appID := a.notificationAppID()

	var script string
	if iconFilePath != "" {
		script = fmt.Sprintf(notificationTemplate, title, content, iconFilePath, appID)
	} else {
		script = fmt.Sprintf(notificationNoIconTemplate, title, content, appID)
	}
	go runScript("notify", script)
}

func (a *fyneApp) ScheduleNotification(n *fyne.Notification, when time.Time) (*fyne.ScheduledNotification, error) {
	if !when.After(time.Now()) {
		return nil, errors.New("scheduled delivery time must be in the future")
	}

	id, err := scheduler.NewID()
	if err != nil {
		return nil, err
	}

	title := n.Title
	content := n.Content
	iconFilePath := a.cachedIconPath()
	delivery := when.UTC().Format(time.RFC3339)
	appID := a.notificationAppID()

	var script string
	if iconFilePath != "" {
		script = fmt.Sprintf(scheduledNotificationTemplate, title, content, iconFilePath, id, delivery, appID)
	} else {
		script = fmt.Sprintf(scheduledNotificationNoIconTemplate, title, content, id, delivery, appID)
	}
	go runScript("schedule", script)
	return fyne.NewScheduledNotification(id, n, when), nil
}

func (a *fyneApp) CancelScheduledNotification(id string) error {
	appID := a.notificationAppID()
	script := fmt.Sprintf(cancelScheduledNotificationTemplate, id, appID)
	go runScript("cancel", script)
	return nil
}

func (a *fyneApp) notificationAppID() string {
	appID := a.UniqueID()
	if appID == "" || strings.Index(appID, "missing-id") == 0 {
		appID = a.Metadata().Name
	}
	return appID
}

// xssyne
// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default, this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if d, ok := a.Driver().(systrayDriver); ok {
		d.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	if d, ok := a.Driver().(systrayDriver); ok {
		d.SetSystemTrayIcon(icon)
	}
}

// SetSystemTrayWindow assigns a window to be shown with the system tray menu is tapped.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayWindow(w fyne.Window) {
	if d, ok := a.Driver().(systrayDriver); ok {
		d.SetSystemTrayWindow(w)
	}
}

func runScript(name, script string) {
	num := scriptNum.Add(1)
	var appID string
	if app := fyne.CurrentApp(); app != nil {
		appID = app.UniqueID()
	}
	fileName := fmt.Sprintf("xssyne-%s-%s-%d.ps1", appID, name, num)

	tmpFilePath := filepath.Join(os.TempDir(), fileName)
	err := os.WriteFile(tmpFilePath, []byte(script), 0o600)
	if err != nil {
		fyne.LogError("Could not write script to show notification", err)
		return
	}
	defer os.Remove(tmpFilePath)

	// %q в Go-формате даёт корректную кавычечную литерал для PowerShell
	// (двойные кавычки, экранирование бэкслэшей/кавычек).
	launch := fmt.Sprintf("(Get-Content -Encoding UTF8 -Path %q -Raw) | Invoke-Expression", tmpFilePath)
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Failed to launch windows notify script", err)
	}
}

func watchTheme(s *settings) {
	go internalapp.WatchTheme(func() {
		fyne.Do(s.setupTheme)
	})
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
