# Telegram bot

A simple library written in Go to build a custom Telegram bot using a YAML file.

## Prerequisites

Create a Telegram bot and note down the `TOKEN` (that's easy, just check the [Telegram docs](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)).

## Clone the project

```
$ git clone https://github.com/pippokairos/telegram-bot.git
```

## Set the token

Just copy the [.env.example](.env.example) file to `.env` and paste your bot token.

## Define the triggers

Copy the [triggers.yml.example](triggers.yml.example) to `triggers.yml` and configure the triggers.
At the moment, it's only possible to return a string. But you can include the user's input string following the input, or return a random string from a set.

### Example 1

```yaml
- key: hello
  values: Hello!
```

"Hello there" -> "Hello!"

### Example 2

```yaml
- key: throw a dice
  values:
    - 1
    - 2
    - 3
    - 4
    - 5
    - 6
```

"Please throw a dice" -> "3"

### Example 3

```yaml
- key: say hello to
  values: Hello __input__, how are you?
```

"Yo say hello to John Doe" -> "Hello John Doe, how are you?"

## Run the server

```
go run .
```

Remember to adjust the `PORT` in the `.env` file.

## Configure the address as webhook

Just call

```
https://api.telegram.org/bot<TOKEN>/setWebhook?url=https://my-server.example.com
```

You can check [the official documentation](https://core.telegram.org/bots/api#setwebhook).

## Run tests

```
go test
```

You don't say!

## Contribute...?

This is just a very simple application I wrote to play with Go, there are already more complete solutions out there and I'm not aiming to expand this library's capabilities so there's no guideline for contributing. But hey, if you'd like to make additions, feel free to do so.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
