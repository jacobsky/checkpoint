# checkpoint

My own dedication to the internet checkpoints of the internet. Come join and relax, leave a comment and chat with anyone that happens to be sitting by when this happens.

- [x] Initial shared view prototype with datastar
    - [ ] Stream some music and show a static gif of a campfire
- [ ] Single playing embedded youtube video set to the same shared timestamp for everyone (curated playlist of internet checkpoints with links to the site)
- [x] Realtime ephemeral chat
- [ ] Comments section
    - [x] Base DB Schema
    - [ ] List Comments
    - [ ] Post Comment
    - [ ] Pinned comments (mostly for admin)
- [ ] Error modal that can be used for an ephemeral "Something went wrong" message when a user does something silly or there is a backend error (patched via signals)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
