package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"strconv"

	"github.com/jrangulo/gif-split/web"
)


func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", uploadFormHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r)
	err := web.UploadFormTemplate().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse form"+err.Error(), http.StatusBadRequest)
		return
	}

    fmt.Println(r)
	file, _, err := r.FormFile("gifFile")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	rows, err := strconv.Atoi(r.FormValue("rows"))
	if err != nil || rows < 1 {
		http.Error(w, "Invalid number of rows", http.StatusBadRequest)
		return
	}

	cols, err := strconv.Atoi(r.FormValue("cols"))
	if err != nil || cols < 1 {
		http.Error(w, "Invalid number of columns", http.StatusBadRequest)
		return
	}

	gridGIFs, err := processGIF(file, rows, cols)
	if err != nil {
		http.Error(w, "Error processing GIF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = web.ImageTableTemplate(gridGIFs).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering result: "+err.Error(), http.StatusInternalServerError)
	}
}

func processGIF(file io.Reader, rows, cols int) ([][]string, error) {
    g, err := gif.DecodeAll(file)
    if err != nil {
        return nil, err
    }

    frameWidth := g.Config.Width
    frameHeight := g.Config.Height
    cellWidth := frameWidth / cols
    cellHeight := frameHeight / rows

    gridGIFs := make([][]string, rows)
    for i := range gridGIFs {
        gridGIFs[i] = make([]string, cols)
    }

    for row := 0; row < rows; row++ {
        for col := 0; col < cols; col++ {
            rect := image.Rect(
                col*cellWidth,
                row*cellHeight,
                (col+1)*cellWidth,
                (row+1)*cellHeight,
            )

            newGIF := &gif.GIF{
                Image:     make([]*image.Paletted, len(g.Image)),
                Delay:     make([]int, len(g.Delay)),
                LoopCount: g.LoopCount,
            }

            for i, srcImg := range g.Image {
                dstImg := image.NewPaletted(image.Rect(0, 0, cellWidth, cellHeight), srcImg.Palette)
                draw.Draw(dstImg, dstImg.Rect, srcImg, rect.Min, draw.Over)
                newGIF.Image[i] = dstImg
                newGIF.Delay[i] = g.Delay[i]
            }

            buf := new(bytes.Buffer)
            err := gif.EncodeAll(buf, newGIF)
            if err != nil {
                return nil, err
            }

            encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
            gridGIFs[row][col] = encoded
        }
    }

    return gridGIFs, nil
}
