package main

import (
	"fmt"
	//"log"
	"os"
//	"strings"

	"github.com/codegangsta/cli"
	//"github.com/daviddengcn/go-colortext"

	"github.com/kunalkushwaha/container-image-manager/lib"
)

func main() {
	app := cli.NewApp()
	app.Name = "ImageManager"
	app.Usage = "Command line tool for LXD image management"
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "Show list of remote image servers",
			Action: func(c *cli.Context) {
				listImageServers(c)
			},
		},
		{
			Name:      "images",
			ShortName: "i",
			Usage:     "Show list of all available images",
			Action: func(c *cli.Context) {
				listImages(c)
			},
		},
	}

	app.Run(os.Args)
}

func listImageServers(c *cli.Context) error {
	var cts []string
	var rm lib.RegistryManager

	rm.InitRegistryManager()
	err := rm.FetchImageServerData()
	if err != nil {
		return err
	}

	cts, err = rm.GetImageServers()

	if err != nil {
		return err
	}

	for _, ct := range cts {
		fmt.Println(ct)
	}
	return err

}

func listImages(c *cli.Context) error {
	var cts []string
	var rm lib.RegistryManager

	rm.InitRegistryManager()
	err := rm.FetchImageServerData()
	if err != nil {
		return err
	}

	cts, err = rm.GetImageList(c.Args().First())

	if err != nil {
		return err
	}

	for _, ct := range cts {
		fmt.Println(ct)
	}
	return err
}
