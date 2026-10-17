// xssyne/canvas/svg_colorize.go

package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/svg"
)

// ColorizeThemedResource возвращает новый StaticResource с содержимым
// SVG-ресурса, перекрашенным в col.
//
// Зачем это нужно:
//
// Внутри canvas.Image колоризация SVG выполняется через
// cache.WidgetTheme(i) — Fyne ищет выше по дереву виджет с темой
// и берёт цвет через res.ThemeColorName(). Если Image живёт в
// обычном контейнере (container.NewVBox → SetContent), а не внутри
// fyne.Widget, cache.WidgetTheme возвращает nil, colorize не
// запускается, и SVG отрисовывается с fallback-цветом currentColor
// (обычно foreground темы). Именно так ведут себя, например,
// placeholder-иконки в карточке описания.
//
// Этот хелпер обходит WidgetTheme: цвет берётся у вызывающего
// (то есть у app-кода, который знает свою тему) и применяется
// напрямую. Итог — обычный StaticResource, готовый к передаче в
// canvas.NewImageFromResource.
//
// res может быть любым fyne.Resource с SVG-содержимым (текущие
// theme.*Icon() подходят). Если содержимое не SVG или svg.Colorize
// вернул ошибку — возвращается исходный res как есть, чтобы
// приложение не падало из-за декоративной иконки.
func ColorizeThemedResource(res fyne.Resource, col color.Color) fyne.Resource {
	if res == nil {
		return nil
	}
	colored, err := svg.Colorize(res.Content(), col)
	if err != nil {
		return res
	}
	return fyne.NewStaticResource(res.Name(), colored)
}
