# Listener
To monitor and debug responses, such as webhooks from a mail service, a gateway is essential, especially when working on a local development environment where direct access to external responses may not be possible.

The gateway serves to intercept and read responses locally, allowing developers to inspect them for debugging purposes. Once intercepted, these responses can be forwarded, for example, through Postman for further analysis or testing.

## Usage

Send any request to 'localhost/listen' to catch request and see what it contains. To see what requests were caught 
go to `licahost:8080` to see UI for this listener.

## Tech Stack
* GO for handling HTTP requests
* MongoDB for storing information
* Vue.js (Typescript) for presentation

## Building
Create `.env.local` for variables
To build docker run `docker compose --env-file .env.local build --no-cache`
Tu run project `docker compose up -d`

## Initial work
Greg, Lord of Mailgun Messages, Master of Redis Realms, Messenger of Symfony Secrets, Champion of Domain-Driven Development, Guardian of PHPUnit Proclamations, and Sage of PHP Sorcery
