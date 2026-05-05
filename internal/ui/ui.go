package ui

type Component interface {
	Render()
	Width() int
	Height() int
	WriteTo()
}
