package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

func draw_guy(g guy) {
	termbox.SetCell(g.x, g.y, '☺', termbox.Attribute(15), termbox.Attribute(7))
	termbox.SetCell(g.x, g.y+1, '✝', termbox.Attribute(15), termbox.Attribute(7))
	termbox.SetCell(g.x, g.y+2, '∩', termbox.Attribute(15), termbox.Attribute(7))
}

func draw_world() {
	world := [30]string{}

	world[0] = "                                                                                "
	world[1] = "                                                                                "
	world[2] = "                                                                                "
	world[3] = "                                                                               !"
	world[4] = "                                                                      ----------"
	world[5] = "                                                                     |          "
	world[6] = "                                                                     |          "
	world[7] = "                                                                     |          "
	world[8] = "                                                               -----            "
	world[9] = "                                                             |                  "
	world[10] = "                                                            |                   "
	world[11] = "                                                            |                   "
	world[12] = "                                                       -----                    "
	world[13] = "                                                      |                         "
	world[14] = "                                                      |                         "
	world[15] = "                                                      |                         "
	world[16] = "                                                -----                           "
	world[17] = "                                              |                                 "
	world[18] = "                                              |                                 "
	world[19] = "                                              |                                 "
	world[20] = "                                              |                                 "
	world[21] = "                                           ---                                  "
	world[22] = "                           |              |                                     "
	world[23] = "                           |              |                                     "
	world[24] = "                           |              |                                     "
	world[25] = "            ------------------------------                                      "
	world[26] = "                                                                                "
	world[27] = "                                                                                "
	world[28] = "                                                                                "
	world[29] = "--------------------------------------------------------------------------------"

	colour := func(x int, y int, world [30]string) int {
		switch world[y][x] {
		case ' ':
			return 7
		case '-':
			return 3
		case '|':
			return 3
		case '!':
			return 2
		}
		return 0
	}
	colour(0, 0, world)
	width, height := 80, 30
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			termbox.SetCell(x, y, rune(world[y][x]), termbox.Attribute(colour(x, y, world)), termbox.Attribute(colour(x, y, world)))
		}
	}
}

func draw_all(g guy) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	draw_world()
	draw_guy(g)

	termbox.Flush()
}

type guy struct {
	x          int
	y          int
	jumpCycles int
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	quit := make(chan string)

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	redrawRoutine := make(chan guy)
	go func(r chan guy) {
		for {
			select {
			case g := <-r:
				draw_all(g)
			}
		}

	}(redrawRoutine)

	go func(t *time.Ticker, gs chan guy) {
		g := guy{x: 10, y: 0, jumpCycles: 0}

		for {
			select {
			case <-t.C:
				if g.jumpCycles == 0 {
					w, h := termbox.Size()
					if g.x+(g.y+3)*w < w*h {
						c := termbox.CellBuffer()[g.x+(g.y+3)*w]
						if c.Ch == 32 {
							g.y++
							gs <- g
						}
					}
				} else if g.jumpCycles > 0 {
					w, _ := termbox.Size()
					if g.x+(g.y-1)*w > 0 {
						c := termbox.CellBuffer()[g.x+(g.y-1)*w]
						if c.Ch == 32 {
							g.y--
							gs <- g
						}
					}
					g.jumpCycles--
				}
			case ev := <-event_queue:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						quit <- "Goodbye!"
					case termbox.KeyArrowLeft:
						w, h := 80, 30
						if g.x-1+(g.y+2)*g.x < w*h {
							c1 := termbox.CellBuffer()[g.x-1+g.y*w]
							c2 := termbox.CellBuffer()[g.x-1+(g.y+1)*w]
							c3 := termbox.CellBuffer()[g.x-1+(g.y+2)*w]

							if g.x > 0 && c1.Ch == 32 && c2.Ch == 32 && c3.Ch == 32 {
								g.x--
								gs <- g
							}
						}
					case termbox.KeyArrowRight:
						w, h := termbox.Size()
						if g.x-1+(g.y+2)*g.x < w*h {
							c1 := termbox.CellBuffer()[g.x+1+g.y*w]
							c2 := termbox.CellBuffer()[g.x+1+(g.y+1)*w]
							c3 := termbox.CellBuffer()[g.x+1+(g.y+2)*w]

							if g.x < 79 && c1.Ch == 32 && c2.Ch == 32 && (c3.Ch == 32 || c3.Ch == 33) {
								g.x++
								gs <- g
							}
						}
					case termbox.KeyArrowUp:
						w, h := termbox.Size()
						if g.x+(g.y+3)*w < w*h {
							c := termbox.CellBuffer()[g.x+(g.y+3)*w]
							if c.Ch != 32 && g.jumpCycles == 0 {
								g.jumpCycles = 5
							}
						}
					}
				}
			}

			if g.y == 1 && g.x == 79 {
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				t.Stop()
				quit <- "You WIN!"
			}
		}
	}(ticker, redrawRoutine)

	m := <-quit
	termbox.Close()
	fmt.Println(m)
}
