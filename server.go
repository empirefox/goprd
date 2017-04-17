package main

import (
	"fmt"
	"net"
	"os"

	arukas "github.com/arukasio/cli"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func main() {
	app := iris.New()
	app.Adapt(httprouter.New())

	app.Get("/", getAddress)

	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	app.Listen(bind)
}

func getAddress(ctx *iris.Context) {
	var parsedContainer []arukas.Container
	client := arukas.NewClientWithOsExitOnErr()

	if err := client.Get(&parsedContainer, "/containers"); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		return
	}

	portMapping := parsedContainer[0].PortMappings[0][0]
	addrs, err := net.LookupHost(portMapping.Host)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		return
	}

	ctx.JSON(iris.StatusOK,
		iris.Map{
			"schema": "tcp",
			"ip":     addrs[0],
			"post":   fmt.Sprintf("%d", portMapping.ServicePort),
			"bind":   fmt.Sprintf("%s:%d", addrs[0], portMapping.ServicePort),
		})
}
