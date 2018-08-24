package main

import (
	"flag"
	"os"
	"sync"

	"./graylog"
)

var debug bool
var appName, team, env, baseURL *string

func init() {
	appName = flag.String("app", "", "Application name to create Dashboards and Streams")
	team = flag.String("team", "", "Team role that application belongs to")
	env = flag.String("env", "", "Application environment (stg|prod)")
	baseURL = flag.String("baseURL", "", "Graylog server URL")
	const (
		defaultDebug = false
		usage        = "Debug message on"
	)
	flag.BoolVar(&debug, "debug", defaultDebug, usage)
	flag.BoolVar(&debug, "d", defaultDebug, usage+" (shorthand)")

	flag.Parse()

	if *appName == "" || *team == "" || *env == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	api := graylog.NewAPI(*baseURL, os.Getenv("GRAYLOG_TOKEN"))

	var wg sync.WaitGroup
	wg.Add(3)

	var dash graylog.Dashboard
	go dash.CreateDashboard(&api, &wg, *appName, *env)
	var stream graylog.Stream
	go stream.CreateStream(&api, &wg, *appName, *env)
	var role graylog.Role
	go role.GetPermissions(&api, &wg, *team)
	wg.Wait()
	role.UpdateRole(&api, stream, dash)

	// var role Role
	// role.CleanRole(api)
}
