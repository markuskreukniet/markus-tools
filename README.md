# markus-tools

An Electron application with Solid

## Recommended IDE Setup

- [VSCode](https://code.visualstudio.com/) + [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) + [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)

## Project Setup

### Install

```bash
$ npm install
```

#### Go

We should have installed Go on our local machine to run the Go code. We can install go from [go.dev](https://go.dev/doc/install).

We can install the [Delve](https://github.com/go-delve/delve) debugger, which this project supports.

### Development

```bash
$ npm run dev
```

<!--
TODO: how to install: IntelliJ IDEA, OpenJDK (if it is used)
TODO: JUnit, how to to name unit tests, given-When-Then
-->

### Unit Testing

#### Go

With `go test`, we can run unit test from 'markus-tools/go' to test the 'go' directory and from 'markus-tools/go/utils' to test the 'go/utils' directory.

### Rules

- We should not use CSS margins. This [video](https://www.youtube.com/watch?v=KVQMoEFUee8) can explain why we should not use it.
- We should not always use `click` on button elements, but instead prefer `mouseDown`. This [video](https://www.youtube.com/watch?v=yaMGtiPckAQ) can explain why we should not use it.

### Build

```bash
# For windows
$ npm run build:win

# For macOS
$ npm run build:mac

# For Linux
$ npm run build:linux
```
