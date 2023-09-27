package main

import (
	"bytes"
	"image"
	"image/color"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

type JsonRequest struct {
	UrlImagem     string          `json:"urlImagem"`
	ListaEntradas []EntradaImagem `json:"listaEntradas"`
}

type EntradaImagem struct {
	Text      string  `json:"text"`
	XPosition float64 `json:"xPosition"`
	YPosition float64 `json:"yPosition"`
	FontSize  float64 `json:"fontSize"`
}

func main() {
	//creates a instance of gin
	r := gin.Default()

	//main endpoint to receive the req as a post request and do all the proccess
	r.POST("/GerarImagem", func(ctx *gin.Context) {

		//load the json with all parameters
		var req JsonRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		//loads the template image from web
		resp, err := http.Get(req.UrlImagem)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error loading image from url"})
			return
		}

		//ensure memory cleanup in the end of method
		defer resp.Body.Close()

		//checks the response
		if resp.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching image"})
		}

		//decode the image
		i, _, err := image.Decode(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding image"})
		}

		//creates a intance of the image editor
		dc := gg.NewContextForImage(i)
		dc.SetColor(color.White)

		for _, ei := range req.ListaEntradas {
			if err = dc.LoadFontFace("fonts/OpenSans-VariableFont_wdth,wght.ttf", ei.FontSize); err != nil {
				ctx.JSON(500, gin.H{"error": "Error loading font"})
				return
			}
			dc.DrawString(ei.Text, ei.XPosition, ei.YPosition)
		}

		var imgData bytes.Buffer

		if err := dc.EncodePNG(&imgData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the image. "})
		}

		ctx.Data(http.StatusOK, "image/png", imgData.Bytes())
	})

	r.Run(":8080")

}
