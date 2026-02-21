//go:build windows && gui

package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	colorBackground    = color.NRGBA{R: 27, G: 27, B: 47, A: 255}
	colorSurface       = color.NRGBA{R: 45, G: 45, B: 68, A: 255}
	colorSurfaceLight  = color.NRGBA{R: 55, G: 55, B: 82, A: 255}
	colorPrimary       = color.NRGBA{R: 91, G: 134, B: 229, A: 255}
	colorPrimaryDark   = color.NRGBA{R: 66, G: 103, B: 190, A: 255}
	colorSuccess       = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
	colorWarning       = color.NRGBA{R: 255, G: 193, B: 7, A: 255}
	colorError         = color.NRGBA{R: 255, G: 82, B: 82, A: 255}
	colorTextPrimary   = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	colorTextSecondary = color.NRGBA{R: 158, G: 158, B: 158, A: 255}
	colorDivider       = color.NRGBA{R: 60, G: 60, B: 90, A: 255}
	colorSidebarBg     = color.NRGBA{R: 20, G: 20, B: 38, A: 255}
	colorCardHover     = color.NRGBA{R: 65, G: 65, B: 100, A: 255}
	colorConnected     = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
	colorDisconnected  = color.NRGBA{R: 158, G: 158, B: 158, A: 255}
)

type dnsTheme struct{}

func newDNSTheme() fyne.Theme {
	return &dnsTheme{}
}

func (t *dnsTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colorBackground
	case theme.ColorNameForeground:
		return colorTextPrimary
	case theme.ColorNamePrimary:
		return colorPrimary
	case theme.ColorNameButton:
		return colorPrimary
	case theme.ColorNameDisabled:
		return colorTextSecondary
	case theme.ColorNameError:
		return colorError
	case theme.ColorNameSuccess:
		return colorSuccess
	case theme.ColorNameWarning:
		return colorWarning
	case theme.ColorNameInputBackground:
		return colorSurface
	case theme.ColorNameInputBorder:
		return colorDivider
	case theme.ColorNamePlaceHolder:
		return colorTextSecondary
	case theme.ColorNameMenuBackground:
		return colorSidebarBg
	case theme.ColorNameOverlayBackground:
		return colorSurface
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 80}
	case theme.ColorNameHover:
		return colorCardHover
	case theme.ColorNameSeparator:
		return colorDivider
	case theme.ColorNameHeaderBackground:
		return colorSurfaceLight
	case theme.ColorNameScrollBar:
		return colorDivider
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (t *dnsTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *dnsTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *dnsTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNameSubHeadingText:
		return 17
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputRadius:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}
