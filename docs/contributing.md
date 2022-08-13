# Contributing

This document serves as a guide of how you can contribute code to
userstyles.world. Make sure to report security-related issues to
[security@userstyles.world](mailto:security@userstyles.world).

## Setup

Before you can start making changes, make sure you are able to run the code.

Userstyles.world's codebase is written in the [Go](https://go.dev/) programming
language, make sure you've installed this on your system. When go is installed,
you will need to install [air](https://github.com/cosmtrek/air) by executing
`go install github.com/cosmtrek/air`. Also make sure you've installed
[sassc](https://github.com/sass/sassc) in order to generate the CSS.

Once you've installed all the necessary tooling, you can start by generating the
CSS for USw, by executing `./tools/run sass` after that you can start the server
by executing `./tools/run dev true true`.

## Making changes.

When making changes, please ensure you're adding code to the right directory,
the directory name should be self-explanatory of where code should go. Please
ensure you've checked out the `dev` branch instead of the `main` branch, the
`dev` branch contains the latest code, `main` serves as a "stable" snapshot.
