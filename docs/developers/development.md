# Developing Hivelocity Cloud Controller

Developing our cloud controller manager is quite easy. First, you need to install some base requirements. Second, you need to follow the quickstart documents to set up everything related to Hivelocity.

## Installation

Install Go. Use `sudo snap install --classic go` or follow the official [Install Go Docs](https://go.dev/doc/install)

Get the source code:
```
> git clone git@github.com:hivelocity/hivelocity-cloud-controller-manager.git
```

Download the dependencies:
```
> go mod tidy
```

Run the unit tests:
```
> make test
```

## Hivelocity API key.

You need an API key to access the Hivelocity API.

You find more docs about accessing Hivelocity with Go in the [Hivelocity Go Client](https://github.com/hivelocity/hivelocity-client-go).


## Submitting PRs and testing

Pull requests and issues are highly encouraged! For more information, please have a look in the [Contribution Guidelines](../../CONTRIBUTING.md)

There are two important commands that you should make use of before creating the PR.

With `make verify` you can run all linting checks and others. Make sure that all of these checks pass - otherwise the PR cannot be merged. Note that you need to commit all changes for the last checks to pass. 

With `make test` all unit tests are triggered.

