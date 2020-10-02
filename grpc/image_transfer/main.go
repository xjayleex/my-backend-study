package main

import (
	clitool "gopkg.in/urfave/cli.v2"
	"github.com/xjayleex/my-backend-study/grpc_study/image_transfer/cli"
	"os"
)
func main() {
	app := &clitool.App{
		Name : "imgtransfer",
		Usage: "send a image",
		Commands: []*clitool.Command{
			&cli.Serve,
			&cli.Send,
		},
		Flags: []clitool.Flag{
			&clitool.BoolFlag{
				Name: "debug",
				Usage: "enables debug logging",
			},
		},
	}

	app.Run(os.Args)
}