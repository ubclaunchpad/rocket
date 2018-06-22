# Rocket 🚀 [![GoDoc](https://godoc.org/github.com/ubclaunchpad/rocket?status.svg)](https://godoc.org/github.com/ubclaunchpad/rocket) [![Build Status](https://travis-ci.org/ubclaunchpad/rocket.svg?branch=master)](https://travis-ci.org/ubclaunchpad/rocket) [![Coverage Status](https://coveralls.io/repos/github/ubclaunchpad/rocket/badge.svg?branch=master)](https://coveralls.io/github/ubclaunchpad/rocket?branch=master)

Rocket is the management and onboarding system for UBC Launch Pad. More information can be found in the [Wiki](https://github.com/ubclaunchpad/rocket/wiki). Rocket is a Slack bot you can talk to at ubclaunchpad.slack.com by messaging `@rocket`. It features GitHub integration, a robust command framework, and a simple interface through which plugins can easily be added.

- [Development](#development)
	- [Creating Your Own Rocket Plugin](#creating-your-own-rocket-plugin)
- [Architecture](#architecture)
	- [Slack Bot](#slack-bot)
	- [Server](#server)
	- [Database](#database)
- [Deployment](#deployment)

<br>

## Development

To get started, make sure you have [Golang](https://golang.org/doc/install#install) installed and download the Rocket codebase:

```bash
$ go get github.com/ubclaunchpad/rocket
$ cd $GOPATH/src/github.com/ubclaunchpad/rocket
$ make                  # install dependencies
$ make test             # run unit tests
```

Additional integration tests can be run if you have `postgres` installed (for Mac users, an easy way is to `brew install postgresql`):

```bash
$ make test-integration  # runs integration tests
```

Make sure you mark integration tests as `-short`-skippable:

```go
func TestMyIntegratedFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	// ...
}
```

### Creating Your Own Rocket Plugin

Features can easily be added to Rocket through Rocket's plugin framework. A Rocket Plugin is simply any type that implements the [Plugin](plugin/plugin.go) interface:

```go
// Plugin is any type that exposes Slack commands and event handlers, and can
// be started.
type Plugin interface {
	// Starts the plugin or returns an error if one occurred.
	// Use this as an opportnity to start background goroutines or do any other
	// additional setup for your plugin.
	Start() error
	// Returns a slice of commands that the plugin handles.
	Commands() []*cmd.Command
	// Returns a mapping from event type to a event handler.
	// See https://api.slack.com/rtm for event types.
	EventHandlers() map[string]bot.EventHandler
}
```

The prime example of a plugin is Rocket's [Core](plugins/core/core.go) plugin which provides basic Launch Pad administration commands for managing teams and users on both Slack and GitHub. An even simpler plugin example is the [welcome plugin](plugins/welcome/welcome.go) that welcomes users when they join our Slack workspace:

```go
// Start starts the welcome plugin.
func (wp *Plugin) Start() error { return nil }

// Commands returns an empty list of commands, because this plugin has no commands.
func (wp *Plugin) Commands() []*cmd.Command { return []*cmd.Command{} }

// EventHandlers returns a map from event type to event handler.
func (wp *Plugin) EventHandlers() map[string]bot.EventHandler {
	return map[string]bot.EventHandler{"team_join": wp.handleTeamJoin}
}
```

You can use the `Start` method of your plugin to start any background tasks you need to. Any `Commands` and `EventHandlers` you expose to Rocket in your implementation of the Plugin interface will be automatically registered with the `Bot`. See the Slack's [API Event Types](https://api.slack.com/events) for a list of events and their names if you implement your own `EventHandler`s for your plugin.

To add your plugin to Rocket, just make a new package for your plugin at the same level as the `core` package (within the `plugins` directory), create your type that implements the `Plugin` interface, and register your plugin in [plugin.RegisterPlugins](plugin/plugin.go). Once you are done, open up a pull request! :tada:

## Architecture

### Slack Bot

The [Bot](bot/bot.go) holds references to structures that we use to communicate with our external dependencies (Slack, GitHub, and Postgres). It also contains logic for handling Slack messages. The `commands` property maps from command name to command handler.

#### Commands

The command framework can be found in `cmd`. It defines a set of data structures and functions for parsing, validating, and automatically documenting Rocket commands. All commands are defined in the `bot` package.

#### Plugins

A Rocket plugin is intended to be a standalone component of Rocket. Rocket's core Slack functionality is implemented as a plugin in [package core](plugins/core).

### Server

[server.go](server/server.go) defines some handlers for HTTP requests. Our website will make requests to `/api/teams` and `/api/members` to display information about our teams and members. Note that content is served over HTTPS using `acme/autocert` to get TLS certificates from LetsEncrypt.

### Database

We use the [go-pg](https://github.com/go-pg/pg) for querying our Postgres database from Rocket. The `dal` package provides an interface to querying our database. The `model` package holds all our data structures that are used by the `dal` package in our queries.

The database schema is defined in [tables.go](schema/tables.sql).

## Deployment

_This section is for reference or for when moving Rocket to a new server. On the current Google Cloud server, Docker is already setup._

We use [Docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/) to run Rocket and the Postgres database that it relies on. In order for Rocket to access the database the rocket container (called "rocket" in `docker-compose.yml`) needs to be running on the same Docker network as the Postgres container (called "postgres" in `docker-compose.yml`). Starting both containers with `docker-compose up` will create a Docker container network called `rocket_default`. Once this is done Rocket will be able to access the DB with the host name `postgres`.

Before deploying, you will have to create two config files using the templates provided in `.app.env.example` and `.db.env.exmaple`. Copy these files and add the relevant values to them. Here are the recommended settings with passwords an security tokens omitted:

#### App Environment Variables

* `ROCKET_HOST`: should essentially always be `0.0.0.0` (bind on all interfaces)
* `ROCKET_PORT`: can be any unreserved port, as long as it is mapped from the container to the host properly in your `docker-compose.yml` under `ports` for the `rocket` service (it's assumed to be port 80 in `docker-compose.yml`)
* `ROCKET_SLACKTOKEN`: get this from Slack
* `ROCKET_GITHUBTOKEN`: get this from Github
* `ROCKET_POSTGRESUSER`: can be anything, but `rocket` is the most sensical choice.
* `ROCKET_POSTGRESPASS`: pick a secure password and make sure it matches `POSTGRES_PASSWORD` in the DB env file
* `ROCKET_POSTGRESDATABASE`: the name of the database to create - it can be anything, but again `rocket` is the most sensical choice

#### DB Environment Variables

* `POSTGRES_USER`: username for the Postgres DB, make sure this matches `ROCKET_POSTGRESUSER`
* `POSTGRES_PASSWORD`: anything secure, make sure this matches `ROCKET_POSTGRESPASS`
* `POSTGRES_DB`: the name of the DB to create, but `rocket` is the most sensical choice - make sure this matches `ROCKET_POSTGRESDATABASE`

These variables are propagated to their respective Docker containers when you do `docker-compose up` via the `env_file` property, so make sure `env_file` for the `rocket` service points to your app environment variables file, and `env_file` for the `postgres` service points to your DB environment variables file. Our `docker-compose.yml` points to `.app.env.example` and `.db.env.example` by default. Note that data is mounted into the Postgres container from a directory called `pgdata` in the `rocket` directory. The first time you do `docker-compose up` this directory will be created for you, the rest of the time it will be re-used. This was if the DB container goes down your data is still on the host machine in the `pgdata` directory.

#### Database Setup

If you're starting the database for the first time you'll need to execute the script defining Rocket's schemas in `schema/tables.sql`:

```bash
# Copy tables.sql into the /tmp folder in the Postgres container
$ docker cp schema/tables.sql <Postgres container ID>:/tmp/
# Run a shell in the Postgres container
$ docker-compose exec postgres bash
# Execute the SQL script against the database
$ psql -U <ROCKET_POSTGRESUSER> -d <ROCKET_POSTGRESDATABASE> < /tmp/tables.sql
# Exit the container
$ exit
```

Note that all the data stored in the DB is mounted into the Postgres container from a directory called `pgdata` in the root folder of this project. This means you can kill the Postgres container and bring it up again and none of your data will be lost.

#### Migrations

If you're updating the DB schema because you want to store a new resource or update an existing one: you'll need to create a migration script under [schema/migrations](schema/migrations) and run it against the DB the same way you would run `schema.sql`. Don't forget to update `schema.sql` to include any changes you apply in your migrations.
