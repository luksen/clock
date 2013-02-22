package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

// commandline args
var color string
var char string
var redraw time.Duration
var version bool
var offset struct{ x, y uint }
var zoom uint
var background string

// data
var font map[rune][][]bool
var colors map[string]string

func init() {
	flag.StringVar(&color, "color", "none", "set the digits' color")
	flag.StringVar(&background, "back", "none", "set background color")
	flag.StringVar(&char, "char", "▀ ", "set the character(s) to use for drawing")
	flag.DurationVar(&redraw, "redraw", 15*time.Second, "set time to wait between redraws")
	flag.BoolVar(&version, "version", false, "print version information and exit")
	flag.UintVar(&offset.x, "x", 2, "x offset (from left)")
	flag.UintVar(&offset.y, "y", 1, "y offset (from top)")
	flag.UintVar(&zoom, "zoom", 1, "display digits X times bigger")
	font = map[rune][][]bool{
		'0': {{true, true, true}, {true, false, true}, {true, false, true}, {true, false, true}, {true, true, true}},
		'1': {{false, false, true}, {false, false, true}, {false, false, true}, {false, false, true}, {false, false, true}},
		'2': {{true, true, true}, {false, false, true}, {true, true, true}, {true, false, false}, {true, true, true}},
		'3': {{true, true, true}, {false, false, true}, {true, true, true}, {false, false, true}, {true, true, true}},
		'4': {{true, false, true}, {true, false, true}, {true, true, true}, {false, false, true}, {false, false, true}},
		'5': {{true, true, true}, {true, false, false}, {true, true, true}, {false, false, true}, {true, true, true}},
		'6': {{true, true, true}, {true, false, false}, {true, true, true}, {true, false, true}, {true, true, true}},
		'7': {{true, true, true}, {false, false, true}, {false, false, true}, {false, false, true}, {false, false, true}},
		'8': {{true, true, true}, {true, false, true}, {true, true, true}, {true, false, true}, {true, true, true}},
		'9': {{true, true, true}, {true, false, true}, {true, true, true}, {false, false, true}, {true, true, true}},
		':': {{false}, {true}, {false}, {true}, {false}},
	}
	colors = map[string]string{
		"black":   "30",
		"red":     "31",
		"green":   "32",
		"yellow":  "33",
		"blue":    "34",
		"magenta": "35",
		"cyan":    "36",
		"white":   "37",
		"none":    "0",
	}
}

func clear() {
	fmt.Printf("\033[2J")
}

func hideCursor() {
	fmt.Printf("\033[?25l")
}

func showCursor() {
	fmt.Printf("\033[?25h")
}

func setAt(x, y int, char rune) {
	fmt.Printf("\033[%d;%dH%c", y+1, x+1, char)
}

func setColor(color string) {
	fmt.Printf("\033[%sm", colors[color])
}

func setColorBack(color string) {
	if color == "none" {
		return
	}
	fmt.Printf("\033[%sm", strings.Replace(colors[color], "3", "4", 1))
}

func drawNumberAt(x, y int, num rune) {
	for l, line := range font[num] {
		i := 0
		for c, col := range line {
			if col {
				for z := 0; z < int(zoom)*2; z++ {
					for zz := 0; zz < int(zoom); zz++ {
						ch := []rune(char)[i%len([]rune(char))]
						setAt(x+c*2*int(zoom)+z, y+l*int(zoom)+zz, ch)
					}
					i++
				}
			} else {
				i += int(zoom) * 2
			}
		}
	}
}

func drawString(x, y int, str string) {
	for _, char := range str {
		drawNumberAt(x, y, char)
		x += (len(font[char][len(font[char])-1]) + 1) * int(zoom) * 2
	}
}

func drawTime(x, y int) {
	drawString(x, y, time.Now().Format("15:04"))
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("clock.go v1.0")
		return
	}

	hideCursor()
	defer showCursor()
	setColor(color)
	defer setColor("none")
	setColorBack(background)
	csig := make(chan os.Signal, 1)
	signal.Notify(csig, os.Interrupt)

	go func() {
		for {
			clear()
			drawTime(int(offset.x), int(offset.y))
			time.Sleep(redraw)
		}
	}()

	<-csig
}
