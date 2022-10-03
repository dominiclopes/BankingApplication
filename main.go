package main

import (
	"os"

	"github.com/urfave/cli"

	"example.com/banking/app"
	"example.com/banking/config"
	"example.com/banking/db"
	"example.com/banking/server"
)

func main() {
	config.Load()
	app.Init()
	defer app.Close()

	cliApp := cli.NewApp()
	cliApp.Name = "GoLang Banking App"
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Action: func(c *cli.Context) {
				server.StartApiServer()
			},
		},
		{
			Name:  "create_migration",
			Usage: "create migration files",
			Action: func(c *cli.Context) {
				db.CreateMigrationFile(c.Args().Get(0))
			},
		},
		{
			Name:  "migrate",
			Usage: "run db migrations",
			Action: func(c *cli.Context) error {
				err := db.RunMigrations()
				return err
			},
		},
		{
			Name:  "rollback",
			Usage: "rollback db migrations",
			Action: func(c *cli.Context) error {
				err := db.RollbackMigration(c.Args().Get(0))
				return err
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
