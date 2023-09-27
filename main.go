package main

import (
	"image/color"

	"github.com/fogleman/gg"
)

func main() {
	i, err := gg.LoadImage("template1.png")
	if err != nil {
		panic(err)
	}
	dc := gg.NewContextForImage(i)
	dc.SetColor(color.White)
	err = dc.LoadFontFace("fonts/OpenSans-VariableFont_wdth,wght.ttf", 18)
	if err != nil {
		panic(err)
	}

	dc.DrawString("TESTE DE ESCRITA", 100, 200)

	dc.SavePNG("resultado.png")
}
