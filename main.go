package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	mainMenu = []string{
		"Command 1",
		"Command 2",
		"Command 3",
		"Command 4",
		"Command 5",
		"Command 6",
		"Command 7",
	}

	subMenuDepth = 5
	screen       tcell.Screen
)

func drawMenu(menuItems []string, selected int, offsetX int) {
	width := 0
	for _, item := range menuItems {
		if len(item) > width {
			width = len(item)
		}
	}

	for row, text := range menuItems {
		style := tcell.StyleDefault
		if row == selected {
			style = style.Reverse(true)
		}

		screen.SetContent(offsetX, row+1, []rune(text)[0], nil, style)
		for col, ch := range text[1:] {
			screen.SetContent(offsetX+1+col, row+1, ch, nil, style)
		}
	}

	// Draw menu borders and title
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	for row := 0; row <= len(menuItems)+1; row++ {
		screen.SetContent(offsetX, row, '║', nil, borderStyle)
		screen.SetContent(offsetX+width+2, row, '║', nil, borderStyle)
	}

	for col := 0; col <= width+2; col++ {
		screen.SetContent(offsetX+col, 0, '═', nil, borderStyle)
		screen.SetContent(offsetX+col, len(menuItems)+1, '═', nil, borderStyle)
	}

	screen.SetContent(offsetX, 0, '╤', nil, borderStyle)
	screen.SetContent(offsetX+width+2, 0, '╤', nil, borderStyle)
	screen.SetContent(offsetX, len(menuItems)+1, '╧', nil, borderStyle)
	screen.SetContent(offsetX+width+2, len(menuItems)+1, '╧', nil, borderStyle)

	title := "Commands"
	for col, ch := range title {
		screen.SetContent(offsetX+1+col, 0, ch, nil, borderStyle)
	}

	screen.Show()
}

func showMenu(depth int, offsetX int) {
	selected := 0
	drawMenu(mainMenu, selected, offsetX)

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				selected--
				if selected < 0 {
					selected = len(mainMenu) - 1
				}
			case tcell.KeyDown:
				selected++
				if selected >= len(mainMenu) {
					selected = 0
				}
			case tcell.KeyEnter:
				if depth < subMenuDepth {
					showMenu(depth+1, offsetX+len(mainMenu[0])+4)
				} else {
					fmt.Printf("Selected: %s\n", mainMenu[selected])
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			}

			drawMenu(mainMenu, selected, offsetX)
		}
	}
}

func main() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.Init()
	screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	screen.Clear()

	showMenu(1, 0)

	screen.Fini()
}
