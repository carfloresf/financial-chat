# Financial Chat challenge
- This repository is a chat application that allows users to send messages in real time.
  - It needs to run both the server and the bot to work correctly, both are in the same repository, but they are in different folders/packages/main.go files.

# How to start Dependencies?
- For Rabbitmq you can start a docker container with the following command: `make rabbitmq`

# How to run server and bot?
- Make sure you have a Rabbitmq running: `make rabbitmq`
- Run `make run-server` to run the server, it will be available on http://localhost:8080
- Run `make run-bot` to run bot, it executes the commands.

# How to test locally?
- Make sure you ran:
  - `make rabbitmq`
  - `make run-server`
  - `make run-bot`
- Open http://localhost:8080 on your browser, you should see the login page.
- You can login with the following credentials, for testing purposes we are setting password to `password`, but only the hash is stored in the DB, and the pepper can be changed in the config file.
  - username: `user`, password: `password`
  - username: `carlos`, password: `password`
  - username: `john`, password: `password`
  - username: `jane`, password: `password`
- Create a channel or join an existing one (hardcoded channels)
- Send messages to the channel, you should see the messages on the chat, and only people connected to that channel should see the messages.
- You can also send commands to the bot, the existing ones are:
  - `/stock=[stock_code]`, for example `/stock=aapl.us` or `/stock=meta.us`, the bot will reply with the stock price
    - If the stock code is not found, the bot will reply with a message saying that the stock code was not found.
  - `/help` command will show the available commands
- If a command is not recognized, the bot will reply with a message saying that the command is not recognized.
- Logout after you are done testing.

# Not Answered Questions - sent via email
- Bonus: Handle messages that are not understood or any exceptions raised within the bot.
  - How do I have to handle these messages? I'm sending "commands" to the bot and the ones that are not understood are responded with a generic message, is that ok? is that what is expected in this line? Also this line mentions "handle messages" not commands, is the bot supposed to read all commands? Or is it ok that it only gets commands (messages starting with "/") from rabbitmq, because I was trying to decouple it using queues for request/response.

# Considerations for the future
- Do we have to store and recover channel messages?
  - Right now, all messages are being sent using websockets, we don't have a DB to store messages, but it could be a nice feature to add.
- Keep confidential information secure.
  - I'm using a pepper to hash the password, and storing just the hash in the DB and a salt secret for session storage. Both can be changed in the application config file, or using environment variables, but that would require to registering users again.
  - Also, for ease of testing, the migration database loads the password hash (hashed with the default pepper value), but in a real world scenario, we would have to hash the password again.
- There is no screen for user registration, but we can register them calling the endpoint directly (described down below).

# How to run tests?
- Run `make lint` to run linter.
- Run `make test` to execute unit tests.

# Implementation details
- This app uses:
  - [Go 1.19](https://golang.org/doc/devel/release.html)
  - [Gin] to handle requests, responses, websockets for chat.
  - [SQLite3] to store data.
  - [Migrate] to manage database migrations.
  - [Mockgen & httpmock] to mock some dependencies in unit-tests.
  - [Golangci-lint](https://github.com/golangci/golangci-lint) to lint code.
  - [Go-rabbitmq] used for RabbitMQ connection, consume and publish.
  - [Melody] used for websocket connection, sessions and messages.

# Endpoints
## RegisterUser - There isn't frontend for this endpoint, but it can be used to register users.
### Description:
Method: POST
>```
>http://localhost:8080/v1/user
>```
### Body (**raw**)

```json
{
  "username":"user",
  "password":"12345"
}
```

### Response (**raw**)

```json
{
  "message": "registration success"
}
```
_________________________________________________
Author: [Carlos Flores](https://github.com/carfloresf)
