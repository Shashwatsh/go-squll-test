package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"os/exec"
	"strconv"
	"fmt"
	"time"
)

type Squll struct {
	Version float32
	Name    string
}
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	API_KEY   string `json:"api_key"`
}

type App struct {
	Mem float32
	CPU float32
	img string
	port map[string]int
	res string
}

func main() {
	app := iris.New()
	// Adapt the "httprouter", faster,
	// but it has limits on named path parameters' validation,
	// you can adapt "gorillamux" if you need regexp path validation!
	app.Adapt(httprouter.New())

	app.HandleFunc("GET", "/", func(ctx *iris.Context) {
		var squll Squll
		squll.Version = 0.1
		squll.Name = "Squll"
		ctx.JSON(iris.StatusOK, squll)
	})
	app.HandleFunc("POST", "/create-app", func(ctx *iris.Context){
		//
		var app App
		CPU_string := ctx.FormValue("CPU")
		CPU_value, err := strconv.ParseFloat(CPU_string, 32)
		if err != nil {
			fmt.Println("Error: error converting string to float32 at cpu")
		}

		app.CPU = float32(CPU_value) // float32(ctx.FormValue("CPU"))

		Mem_value, err := strconv.ParseFloat(ctx.FormValue("Mem"), 32)
		if err != nil {
			fmt.Println("Error: error converting string to float32 at mem")
		}
		app.Mem = float32(Mem_value)

		app.img = ctx.FormValue("img")

		app.res = "Container built & started!"

		memory := strconv.FormatFloat(float64(app.Mem), 'f', -1, 32)
		cpu := strconv.FormatFloat(float64(app.CPU), 'f', -1, 32)
		cmd := exec.Command("/bin/rkt", "run", "--memory=" + string(memory) + "M", "--cpu=" + string(cpu), "--insecure-options=image", string(app.img))
		time.Sleep(2000)
		stdout, err := cmd.Output()

		if err != nil {
			ctx.JSON(iris.StatusOK, err.Error())
			return
		}

		ctx.JSON(iris.StatusOK, string(stdout))

		ctx.JSON(iris.StatusOK, app)

	})

	app.Listen(":8010")
}