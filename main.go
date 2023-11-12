// Copyright 2015 Tamás Gulácsi
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Package main of fronius gets the data from Solar.Web
package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	stdlog "log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var logger = log.NewLogfmtLogger(os.Stderr)
	stdlog.SetOutput(log.NewStdlibAdapter(logger))

	var (
		postgresUri = "postgresql://postgres:password@192.168.178.50:5432/fronius?sslmode=disable"
		servePath   = "/solarapi/v1/current/"
	)

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "accept push from the fronius datalogger",
		Run: func(_ *cobra.Command, args []string) {
			postgresClient := newPostgresClient(postgresUri, logger)

			http.Handle(servePath, solarAPIAccept{postgresClient})
			addr := ":15015"
			if len(args) > 0 {
				addr = args[0]
			}
			level.Info(logger).Log("msg", "Start listening", "address", addr, "path", servePath)
			http.ListenAndServe(addr, nil)
		},
	}

	flags := serveCmd.Flags()
	flags.StringVar(&servePath, "serve.path", servePath, "HTTP endpoint to publish")
	flags.StringVar(&postgresUri, "server", postgresUri, "Postgres URI to connect to")

	mainCmd := &cobra.Command{
		Use: "fronius",
		Run: func(_ *cobra.Command, args []string) {
			serveCmd.Run(serveCmd, args)
		},
	}

	if _, _, err := mainCmd.Find(os.Args[1:]); err != nil && strings.HasPrefix(err.Error(), "unknown command") {
		mainCmd.SetArgs(append([]string{"serve"}, os.Args[1:]...))
	}
	mainCmd.Execute()
}
