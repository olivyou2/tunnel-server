package main

import (
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type SessManager struct {
	tunnelSessManager *TunnelSessionManager
	tunnelServer      *TunnelServer

	frontServerManager *FrontServerManager
}

func createUuid() string {
	uuid := uuid.New()
	return uuid.String()
}

func main() {
	var address string = "localhost:4000"

	(&cli.App{
		Name:  "molar",
		Usage: "a tcp tunnel fits perfectly with you",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Value:       address,
				Usage:       "A address for tunneling",
				Destination: &address,
			},
		},
		Action: func(context *cli.Context) error {
			sm := new(SessManager)

			tunnelSessManager := newTunnelSessionManager()
			sm.tunnelSessManager = tunnelSessManager

			sm.frontServerManager = newFrontServerManager()

			tunnelServer := createTunnelServer(sm, address)

			sm.tunnelServer = tunnelServer

			go tunnelServer.accept()
			select {}
		},
	}).Run(os.Args)
}
