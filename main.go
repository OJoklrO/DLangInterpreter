package main

import (
	"github.com/OJoklrO/Interpreter/drawer"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
)

var d = drawer.NewDrawer()

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(static.Serve("/", static.LocalFile("./source/", false)))

	r.GET("/", GetStatic)
	r.POST("/cmd", ExecuteCmd)

	log.Fatal(r.Run(":8080"))
}

func GetStatic(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ExecuteCmd(c *gin.Context){
	input := c.PostForm("cmd")

	l := NewLexer()
	l.Init()

	t := l.Input(input + ";")

	p := NewParser(t)
	if !p.Parse() {
		p.LogError()
		return
	}

	ret := "1"

	switch p.StmtType {
	case ORIGINSTMT:
		ret = d.SetOrigin(p.O.Pos)
	case SCALESTMT:
		ret = d.SetScale(p.S.Scale)
	case ROTSTMT:
		ret = d.SetRot(p.R.Rot)
	case FORSTMT:
		d.Draw(p.F.Points)
		ret = d.Save("./source/")
	case RESETSTMT:
		d = d.NewDrawer()
		ret = "Reset"
	}

	c.String(200, ret)
}