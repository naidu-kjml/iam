# End to End tests

These tests create the Kiwi IAM service and required services (Okta nginx, Redis) and runs tests against it to guarantee that our service is running correctly.

## Running locally

When running locally it is important to first set up environment variables in your shell.

```sh
export GITLAB_USERNAME=tester@kiwi.com
export GITLAB_PASSWORD=your-access-token-from-gitlab
make run-e2e
```
