# Rocket 🚀 [![GoDoc](https://godoc.org/github.com/ubclaunchpad/rocket?status.svg)](https://godoc.org/github.com/ubclaunchpad/rocket) [![Build Status](https://travis-ci.org/ubclaunchpad/rocket.svg?branch=master)](https://travis-ci.org/ubclaunchpad/rocket) [![Coverage Status](https://coveralls.io/repos/github/ubclaunchpad/rocket/badge.svg?branch=master)](https://coveralls.io/github/ubclaunchpad/rocket?branch=master)

Rocket is the management and onboarding system for UBC Launch Pad. More information can be found in the [Wiki](https://github.com/ubclaunchpad/rocket/wiki).

## Architecture

### Rocket

The [Bot](bot/bot.go) holds references to structures that we use to communicate with our external dependencies (Slack, GitHub, and Postgres). It also contains logic for handling Slack messages. The `commands` property maps from command name to command handler.

[server.go](server/server.go) defines some handlers for HTTP requests. Our website will make requests to `/api/teams` and `/api/members` to display information about our teams and members. Note that content is served over HTTPS using `acme/autocert` to get TLS certificates from LetsEncrypt.

#### Plugins

A plugin is intended to be a standalone component of Rocket. A Rocket Plugin is simply any type that implements the [Plugin](plugin/plugin.go) interface:

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

The prime example of a plugin is Rocket's [Core](core/core.go) plugin which provides basic Launch Pad administration commands for managing teams and users on both Slack and GitHub. An even simpler plugin example is the [WelcomePlugin](welcome/welcome.go) that welcomes users when they join our Slack workspace:

```go
// Start starts the welcome plugin.
func (wp *WelcomePlugin) Start() error { return nil }

// Commands returns an empty list of commands, because this plugin has no commands.
func (wp *WelcomePlugin) Commands() []*cmd.Command { return []*cmd.Command{} }

// EventHandlers returns a map from event type to event handler.
func (wp *WelcomePlugin) EventHandlers() map[string]bot.EventHandler {
	return map[string]bot.EventHandler{"team_join": wp.handleTeamJoin}
}
```

You can use the `Start` method of your plugin to start any background tasks you need to. Any `Commands` and `EventHandlers` you expose to Rocket in your implementation of the Plugin interface will be automatically registered with the `Bot`. See the Slack's [API Event Types](https://api.slack.com/events) for a list of events and their names if you implement your own `EventHandler`s for your plugin.

When creating a new plugin, make a new package for your plugin at the same level as the `core` package, create your type that implements the `Plugin` interface, and register your plugin in [plugin.RegisterPlugins](plugin/plugin.go). It is recommended that you place any commands you write for your plugin in their own separate files under your plugin's package.

#### Commands

The command framework can be found in `cmd`. It defines a set of data structures and functions for parsing, validating, and automatically documenting Rocket commands. All commands are defined in the `bot` package.

New commands should go in their own files in the `bot` package. When creating a new command you must define the following properties:

* `Name`: The command name. Rocket will use this to assign a Slack message to a specific command handler in [bot/bot.go:handleMessageEvent](bot/bot.go).
* `HelpText`: A description of what the command does. You don't need to describe the options here as you'll do that in the `HelpText` field of the `Option` struct.
* `Options`: A mapping of option key to option. The key for a given option in the `Options` map should always match the `key` field in that option.
* `HandleFunc`: The `CommandHandler` that executes the command. It should take `cmd.Context` as it's only argument and return a `string` response message with `slack.PostMessageParameters`.

#### Options

When specifying an option for a command you'll need to fill in the following fields:

* `Key`: The key that identifies this option. Of course, keys for different options under the same command should always be unique. For exmaple, one might create a command with one option who's key is `name`. In this case the user would assign a value to this key in their Slack command with `name={myvalue}`.
* `HelpText`: A description of what the option is used for.
* `Format`: A `regexp.Regexp` object that specifies the required format of a value for an option. The `cmd` framework will enforce that this format is met when a user enters a value for a given option, and will return an appropriate error response if this is not the case. Commonly used format `Regex`s can be found in [bot/util.go](bot/util.go).
* `Required`: Whether or not a value for this option is required when a user uses this command. The `cmd` framework will enforce that a value is set for each required option when a user enters a command, and will return an appropriate error if this is not the case.

#### Querying the DB

We use the [go-pg](https://github.com/go-pg/pg) for querying our Postgres database from Rocket. The `dal` package provides an interface to querying our database. The `model` package holds all our data structures that are used by the `dal` package in our queries.

### Postgres

#### Schema

Our schema is defined in [tables.go](schema/tables.sql). If you're starting the database for the first time you'll need to execute that script:

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

If you're updating the DB schema because you want to store a new resource or update an existing one:
 you'll need to create a migration script under [schema/migrations](schema/migrations) and run it against the DB the same way you would run `schema.sql`. Don't forget to update `schema.sql` to include any changes you apply in your migrations.

## Docker Setup

_This section is for reference or for when moving Rocket to a new server. On the current Google Cloud server, Docker is already setup._

We use [Docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/) to run Rocket and the Postgres database that it relies on. In order for Rocket to access the database the rocket container (called "rocket" in `docker-compose.yml`) needs to be running on the same Docker network as the Postgres container (called "postgres" in `docker-compose.yml`). Starting both containers with `docker-compose up` will create a Docker container network called `rocket_default`. Once this is done Rocket will be able to access the DB with the host name `postgres`.

### Deployment

Before deploying you will have to create two config files using the templates provided in `.app.env.example` and `.db.env.exmaple`. Copy these files and add the relevant values to them. Here are the recommended settings with passwords an security tokens omitted:

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
