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
5. Create a directory in `games/` that's the name of your game (e.g. `games/superFunTimeGame`)
6. Create a file in that directory called the name of your game that ends in `.go` (e.g. `games/superFunTimeGame/superFunTimeGame.go`)
7. Copy the contents of `games/yourGame/yourGame.go` into your new file
8. Change the package name on line 1 to be the name of your game (e.g. change `yourGame` to `superFunTimeGame`)
9. Change the usage line of your game on line 9 to be the string people will use on the command-line to call your game (e.g. change `your-game` to `super-fun-time-game`)
10. Code your game using [Ebitengine](https://github.com/hajimehoshi/ebiten) (use the `deliveryDash` game as an example of using Ebitengine to code a fun 2D game)
11. Optionally write some instructions on the top of the file so people know how to play your game
12. Change the package name on lines 8 and 21 of `main.go` to be the same as your package name (e.g. change `yourGame` to `superFunTimeGame`)
    * The line numbers in this instruction are guidelines, if multiple people have added different games, these lines could change
13. Run your game
    ```sh
    # Example where game is named superFunTimeGame and ensure your command-line string matches what you set in Step 9
    go run main.go super-fun-time-game
    ```

## Helpful Hints

Tidy and format code:

```sh
go fmt ./...
```

Test your code:

```sh
go test ./...
```

Install dependencies and clean unnecessary dependencies:

```sh
go mod tidy
```

Build a binary you can run repeatedly on your machine:

```sh
go build -o bin/go-games .
```
