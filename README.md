# Governant

| What          | Where                                                             |
| ------------- | ----------------------------------------------------------------- |
| Documentation | Not yet                                                           |
| Discussion    | [#proj-kiwi-iam](https://skypicker.slack.com/messages/CE7KLCKCP/) |
| Maintainer    | [@simon](https://gitlab.skypicker.com/simon.prochazka/)           |

## Usage

Create `.env.yaml` file and set environment variables. Check `.env-sample.yaml`
to see all the possible variables.

To install dependencies and start the project, make sure your GOPATH is set,
and run:

```
make
make start
```

You can use `make dev` on development to reload the server automatically on file
changes.

# Redis

To run this project you need to have redis installed and running, you can use
the following commands to do so.

```shell
# MacOS
brew install redis
brew services start redis

# Linux
# install using your package manager
systemctl start redis

# Docker
docker run -it --rm -p 6379:6379 redis
```

You can use `redis-cli` to interact with redis (ie. check keys and values).
Useful commands:

```shell
# Show all keys
KEYS *

# Show value for a key
GET <key>
```

# Secret Manager

This service uses Vault for syncing secrets to our app.

This service expects the following structure for Vault:
/secret/governant/app_tokens
/secret/governant/settings

## Contributing

- Run `make test/ci` and `make lint` before pushing changes
- Commit messages should be at most 72 characters long
- Commit messages should start with the scope of the changes introduced by the commit
