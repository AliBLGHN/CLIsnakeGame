package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/nsf/termbox-go"
)

const snakeBodyValue = "Y"

var playGround *Pixel
var snakeHead *Pixel

var maxWidth int
var maxHeight int

var PTypeMap = map[int]string{
	0:  "-", // empty pixel
	-1: "*", // food
	-2: "O", // border
}

type Pixel struct {
	PType  int
	Right  *Pixel
	Left   *Pixel
	Top    *Pixel
	Bottom *Pixel
}

func isZeroValue(val interface{}) bool {
	if val != nil {
		return reflect.ValueOf(val).IsZero()
	}
	return true
}

func setDefaultValue(pixel *Pixel, currentWidth, currentHeight, maxWidth, maxHeight int) {
	if currentWidth == 0 || currentHeight == 0 || currentWidth == maxWidth-1 || currentHeight == maxHeight-1 {
		pixel.PType = -2
	} else if currentHeight == maxHeight/2 && currentWidth == maxWidth/2 {
		pixel.PType = 4
	} else {
		pixel.PType = 0
	}
}

func createPlayGround(maxWidth, maxHeight int) (returnPx, snakeHead *Pixel) {
	tempFirstPixel := new(Pixel)
	tempLeftPixel := new(Pixel)
	tempTopPixel := new(Pixel)

	for i := 0; i < maxHeight; i++ {
		for j := 0; j < maxWidth; j++ {
			newPx := new(Pixel)
			setDefaultValue(newPx, j, i, maxWidth, maxHeight)
			if !isZeroValue(*tempLeftPixel) { // İf left px is exist
				tempLeftPixel.Right = newPx
				newPx.Left = tempLeftPixel
				if j != maxWidth-1 {
					tempLeftPixel = tempLeftPixel.Right
				} else {
					tempLeftPixel = new(Pixel)
				}
			} else {
				if i == 0 && j == 0 {
					returnPx = newPx
				}
				tempFirstPixel = newPx
				tempLeftPixel = newPx
			}
			if !isZeroValue(*tempTopPixel) {
				tempTopPixel.Bottom = newPx
				newPx.Top = tempTopPixel
				if j != maxWidth-1 {
					tempTopPixel = tempTopPixel.Right
				} else {
					tempTopPixel = tempFirstPixel
				}
			}
			if newPx.PType == 4 {
				snakeHead = newPx
				snakeHead.Left.PType = 3
				snakeHead.Left.Left.PType = 2
				snakeHead.Left.Left.Left.PType = 1
			}
		}
		tempTopPixel = tempFirstPixel
	}
	return
}

func printPlayGround() {
	FirstPixel := playGround
	tempFirstPixel := FirstPixel
	for {

		if FirstPixel.PType > 0 {
			fmt.Print(snakeBodyValue)

		} else {
			fmt.Print(PTypeMap[FirstPixel.PType])

		}
		if !isZeroValue(FirstPixel.Right) {
			FirstPixel = FirstPixel.Right
		} else if !isZeroValue(tempFirstPixel.Bottom) {
			fmt.Println()
			FirstPixel = tempFirstPixel.Bottom
			tempFirstPixel = tempFirstPixel.Bottom
		} else {
			fmt.Println()
			break
		}
	}
}

func findDependency(snakeTail *Pixel) *Pixel {
	if snakeTail.Top.PType == snakeTail.PType+1 {
		return snakeTail.Top

	} else if snakeTail.Right.PType == snakeTail.PType+1 {
		return snakeTail.Right

	} else if snakeTail.Bottom.PType == snakeTail.PType+1 {
		return snakeTail.Bottom

	} else if snakeTail.Left.PType == snakeTail.PType+1 {
		return snakeTail.Left

	}
	return snakeTail
}

func moveSnake(direction *string, snakeHead, snakeTail *Pixel) (sh, st *Pixel) {
	var nextPixel *Pixel
	st = snakeTail

	switch *direction {
	case "top":
		nextPixel = snakeHead.Top
	case "left":
		nextPixel = snakeHead.Left
	case "right":
		nextPixel = snakeHead.Right
	case "bottom":
		nextPixel = snakeHead.Bottom
	}

	if -1 <= nextPixel.PType && nextPixel.PType <= 0 {

		if nextPixel.PType != -1 {
			st = findDependency(snakeTail)
			snakeTail.PType = 0
		} else {
			placeFood()
		}
		nextPixel.PType = snakeHead.PType + 1
		sh = nextPixel

	} else {
		*direction = "esc"
	}
	return
}

func waitKey(direction *string) {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Termbox başlatılamadı:", err)
		return
	}

	defer termbox.Close()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc {
				*direction = "esc"
			} else if ev.Key == termbox.KeyArrowRight && *direction != "left" {
				*direction = "right"
			} else if ev.Key == termbox.KeyArrowLeft && *direction != "right" {
				*direction = "left"
			} else if ev.Key == termbox.KeyArrowUp && *direction != "bottom" {
				*direction = "top"
			} else if ev.Key == termbox.KeyArrowDown && *direction != "top" {
				*direction = "bottom"
			}
		}
		time.Sleep(time.Second / 4)
	}
}

func findAvailablePixelForFood(h, w int) (pixel *Pixel) {
	pixel = playGround.Bottom.Right
	for i := 1; i < h; i++ {
		pixel = pixel.Bottom
	}
	for i := 1; i < w; i++ {
		pixel = pixel.Right
	}
	return
}

func placeFood() {
	for {
		h := rand.Intn(maxHeight-1) + 1
		w := rand.Intn(maxWidth-1) + 1

		foodPx := findAvailablePixelForFood(h, w)
		if foodPx.PType == 0 {
			foodPx.PType = -1
			break
		}
	}
}

func startGame() {
	clearScreen := "\033[H\033[2J" // Konsolu temizlemek için ANSI kaçış dizisi
	direction := "right"
	snakeTail := snakeHead.Left.Left.Left

	go waitKey(&direction)

	placeFood()

	for {
		if direction == "esc" {
			fmt.Println(clearScreen)
			break
		} else {
			snakeHead, snakeTail = moveSnake(&direction, snakeHead, snakeTail)
		}
		printPlayGround()

		time.Sleep(time.Second / 4)
		fmt.Println(clearScreen)
	}
}

func main() {
	maxWidth = 33
	maxHeight = 11

	playGround, snakeHead = createPlayGround(maxWidth, maxHeight)
	startGame()
}
