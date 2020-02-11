# Nibbler Sample

A sample web app built atop of [Nibbler](https://github.com/markdicksonjr/nibbler).  It has registration, login, forgot 
password, and email validation features out of the box.

## Requirements
- Go v1.12+ (unless you want to try to build without mod)
- SendGrid API key
- Dep (optional)

## Getting Started

The only thing that must be configured to really play around with this sample is a SendGrid
API Key.  Set it in the SENDGRID_API_KEY environment variable.  To build and run the app:

`cd public/vue && npm install && npm run build && cd -`

`PORT=3000 go run main.go`

## Configuration

By default, [Nibbler](https://github.com/markdicksonjr/nibbler) uses environment variables for core 
configuration, but also allows for ./config.json to drive configuration options (overriding the env 
vars), though no configuration file is used in this sample.  More details can be found in the 
[Nibbler](https://github.com/markdicksonjr/nibbler) README.

Also by default, [Nibbler](https://github.com/markdicksonjr/nibbler) uses the PORT environment variable 
to decide which port it serves on (defaulting to 3000), so it is ready to deploy to PaaS platforms such 
as Heroku.

## Database

By default, an in-memory database driven by sqlite will be used.  To use a persistent SQL database, set 
the DATABASE_URL environment variable appropriately.

## UI

The UI will be served by Nibbler (by default, out of the ./public directory).  It is written with Vue 2.