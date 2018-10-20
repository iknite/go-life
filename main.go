package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func main() {
	rand.Seed(time.Now().Unix())
	stop := false
	// init screen
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()

	width, height := s.Size()

	// init board
	b := NewBoard(height, width/2)
	b.Init()
	// b.Set(1, 1, [][]bool{
	// 	{false, false, false, true, true, false, false, false},
	// 	{false, false, true, false, false, true, false, false},
	// 	{false, true, false, false, false, false, true, false},
	// 	{true, false, false, false, false, false, false, true},
	// 	{true, false, false, false, false, false, false, true},
	// 	{false, true, false, false, false, false, true, false},
	// 	{false, false, true, false, false, true, false, false},
	// 	{false, false, false, true, true, false, false, false},
	// })

	b.Set(80, 40, [][]bool{
		{false, true, false, false, false, false, false, false},
		{false, false, false, true, false, false, false, false},
		{true, true, false, false, true, true, true, false},
	})

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan struct{})
	stopSwtich := make(chan struct{})
	reset := make(chan struct{})
	step := make(chan struct{})
	clear := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					x, y := ev.Position()
					b.Get(x/2, y).Switch()
				default:
					continue
				}
			case *tcell.EventResize:
				continue
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEnter {
					step <- struct{}{}
				} else if ev.Key() == tcell.KeyEsc || ev.Rune() == 'q' {
					done <- struct{}{}
				} else if ev.Rune() == ' ' {
					stopSwtich <- struct{}{}
				} else if ev.Rune() == 'c' {
					clear <- struct{}{}
				} else if ev.Rune() == 'r' {
					reset <- struct{}{}
				}
			default:
				continue
			}
		}
	}()

	for {
		s.Clear()
		for i, row := range b.State() {
			for j, cell := range row {
				st := tcell.StyleDefault.Background(ThemeBlackAndWhite.Color(cell.LiveTime()))
				if cell.State() {
					s.SetCell(j*2, i, st, ' ')
					s.SetCell(j*2+1, i, st, ' ')
				}
			}
		}

		select {
		case <-reset:
			stop = !stop
			b.Init()
			grid := make([][]bool, height)
			for i := 0; i < height; i++ {
				grid[i] = make([]bool, width)
			}

			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {
					grid[i][j] = rand.Int()%2 == 0
				}
			}
			b.Set(0, 0, grid)
			stop = !stop
		case <-stopSwtich:
			stop = !stop
		case <-step:
			b.Next()
			s.Show()
		case <-done:
			s.Fini()
			os.Exit(0)
		case <-clear:
			b.Init()
			s.Show()
		case <-ticker.C:
			if !stop {
				b.Next()
			}
			s.Show()
		}
	}
}
