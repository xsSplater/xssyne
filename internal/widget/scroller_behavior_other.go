//go:build linux || windows

// xssyne/internal/widget/scroller_behavior_other.go

package widget

func isScrollerPageOnTap() bool {
	return false
}

func scrollBarAlwaysVisible() bool {
	return true
}

func subscribeScrollerStyle(_ func()) uint64 { return 0 }
func unsubscribeScrollerStyle(_ uint64)      {}
