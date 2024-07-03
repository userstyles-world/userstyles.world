# Contributing

Last updated 2022-08-13.

This document serves as a guide of how you can get started with contribute code
to UserStyles.world project. Make sure to report security-related issues to
[security@userstyles.world](mailto:security@userstyles.world).

## Install dependencies

Before you can start making changes, make sure you're able to set up a
development environment. At the moment, having a Unix-like environment is a
requirement for our `run` tool.

Depending on your Operating System, there are different ways to install/build
dependencies. You'll have to figure out how to install/build them on your own.
It's unlikely that you'll need to build them from scratch, but it's possible.

### Install Go

Our codebase uses the [Go](https://go.dev/) programming language. With Go
installed, you will need to install [air](https://github.com/air-verse/air) by
running `go install github.com/air-verse/air@latest` if you want to use `run
watch` command.

### Install Sass

To compile Sass, you'll need to install [sassc].

[sassc]: https://github.com/sass/sassc

### Install Vips

To convert images, you'll need to install [libvips].

[libvips]: https://github.com/libvips/libvips#install

## Using `run` tool

Once you've installed all the necessary tooling, you need to run `./tools/run
setup` (or `sh tools/run setup`) in order to compile Sass, TypeScript, Go, and
seed the database.

With that complete, run `./tools/run help` and `./tools/run help command` to
familiarize yourself with how to use the `run` tool. Usually, you'll only need
to use `./tools/run watch all` during development, but do read the documentation
for commands if you need something more. To stop the process, press `CTRL-C`.

## Making changes

When making changes, please ensure you're adding code to the right directory.
The directory names should be self-explanatory where code should go. Please
ensure you've checked out the `dev` branch instead of the `main` branch. The
`dev` branch contains the latest code, `main` serves as a "stable" snapshot.

Keep your changes relevant. It's always better to create multiple small pull
requests than a single big one. In addition, try to create a new branch that's
based on the upstream `dev` branch and name it according to what it does.

## Happy hacking!

That should be it. If you have any further questions, don't hesitate to reach
out via any methods we list in the footer of our website.
