# Go Games

An easy way to make simple command-line games with Go.

## Getting Started

1. Install Go by following instructions here: https://go.dev/doc/install
2. Fork the repository
3. Clone the repository to your machine
4. Install dependencies
  ```sh
  go mod tidy
  ```
5. Place your game's first file in a package in the `games/` directory (use the `games/yourGame/yourGame.go` as a template)
6. (optional) Rename the directory of the file (e.g. change `games/yourGame` to be `games/superFunTimeGame`)
7. (optional) Rename the file in the directory (e.g. change `games/superFunTimeGame/yourGame.go` to be `games/superFunTimeGame/superFunTimeGame.go`)
8. Change the package name on line 1 to be the name of your game (e.g. change `yourGame` to `superFunTimeGame`)
9. Change the usage line of your game on line 22 to be the string people will use on the command-line to call your game (e.g. change `your-game` to `super-fun-time-game`)
10. Code your game using [Ebitengine](https://github.com/hajimehoshi/ebiten)
  * Use the `deliveryDash` game as an example of using Ebitengine to code a fun 2D game
11. (optional) Write some instructions on the top of the file so people know how to play your game
12. Change the package name on lines 8 and 21 of `main.go` to be the same as your package name (e.g. change `yourGame` to `superFunTimeGame`)
13. Run your game
  ```bash
  # Example where game is named superFunTimeGame
  go run main.go super-fun-time-game
  ```

## Helpful Hints

#### Tidy and format code

```sh
- go mod tidy
- go mod vendor
- go generate ./...
- go fmt ./...
```

#### Build a macOS binary

```sh
go build -o bin/go-games .
```
