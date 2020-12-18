package drawer

import (
	"fmt"
	"github.com/OJoklrO/Interpreter/parser"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
)

type Drawer struct {
	origin parser.Vector2
	scale parser.Vector2
	rot float64
	img *image.RGBA
	color color.RGBA

	photoIndex int
}

func NewDrawer() *Drawer {
	return &Drawer{
		origin: parser.Vector2{ X: 0, Y: 0 },
		scale:  parser.Vector2{ X: 1, Y: 1 },
		rot:    0,
		img: image.NewRGBA(image.Rect(0, 0, 1000, 600)),
		color: color.RGBA{
			R: 64,
			G: 169,
			B: 255,
			A: 255,
		},
		photoIndex: 0,
	}
}

func (d *Drawer) NewDrawer() *Drawer {
	res := NewDrawer()
	res.photoIndex = d.photoIndex
	return res
}

func (d *Drawer) SetOrigin(p parser.Vector2) string {
	ret := "Origin is ( "
	x := strconv.Itoa(int(p.X))
	ret = ret + x + ", "
	y := strconv.Itoa(int(p.Y))
	ret = ret + y + " )"
	d.origin.X = p.X
	d.origin.Y = p.Y

	return ret
}

func (d *Drawer) SetScale(p parser.Vector2) string {
	ret := "Scale is ( "
	x := strconv.Itoa(int(p.X))
	ret = ret + x + ", "
	y := strconv.Itoa(int(p.Y))
	ret = ret + y + " )"
	d.scale.X = p.X
	d.scale.Y = p.Y

	return ret
}

func (d *Drawer) SetRot(r float64) string {
	d.rot = r
	return "Rot is " + strconv.FormatFloat(r, 'f', -1, 64)
}

func (d *Drawer) Draw(points []parser.Vector2) {
	for _, p := range points {
		var temp parser.Vector2
		temp.X, temp.Y = p.X, p.Y
		// scale
		temp.X *= d.scale.X
		temp.Y *= d.scale.Y
		// rot
		x, y := temp.X, temp.Y
		temp.X = x * math.Cos(d.rot) + y * math.Sin(d.rot)
		temp.Y = y * math.Cos(d.rot) - x * math.Sin(d.rot)
		// origin
		temp.X, temp.Y = temp.X + d.origin.X, temp.Y + d.origin.Y

		fmt.Println(temp)
		d.img.Set(int(temp.X), int(temp.Y), d.color)
	}
}

func (d *Drawer) Save(path string) string {
	fileName := "p" + strconv.Itoa(d.photoIndex) + ".png"
	file, err := os.Create(path + fileName)
	if err != nil {
		return ""
	}
	err = png.Encode(file, d.img)
	if err != nil {
		return ""
	}

	d.photoIndex++
	return fileName
}