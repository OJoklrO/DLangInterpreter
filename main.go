package main

import (
	"github.com/OJoklrO/Interpreter/drawer"
	"github.com/OJoklrO/Interpreter/parser"
	"github.com/OJoklrO/Interpreter/scanner"
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

	l := scanner.NewLexer()
	l.Init()

	t := l.Input(input + ";")

	p := parser.NewParser(t)
	p.Parse()

	ret := "1"

	switch p.StmtType {
	case parser.ORIGINSTMT:
		ret = d.SetOrigin(p.O.Pos)
	case parser.SCALESTMT:
		ret = d.SetScale(p.S.Scale)
	case parser.ROTSTMT:
		ret = d.SetRot(p.R.Rot)
	case parser.FORSTMT:
		d.Draw(p.F.Points)
		ret = d.Save("./source/")
	case parser.RESETSTMT:
		d = d.NewDrawer()
		ret = "Reset"
	}

	c.String(200, ret)
}