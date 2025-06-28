package deliveryDash

// # Delivery Dash: A Dynamic Maze Delivery Challenge 🚗📦

// Introducing Delivery Dash, a fast-paced maze game where you play as a delivery driver navigating through a chaotic, ever-changing city. This isn't your typical maze game - the walls shift beneath your wheels and your delivery destination keeps moving!

// ## Key Features
// - **Dynamic Environment**: The maze walls constantly shift and change, creating a unique challenge every time
// - **Moving Target**: Your delivery destination moves along the bottom of the maze, requiring quick thinking and adaptable strategy
// - **Precision Controls**: One-press-one-move mechanics that reward careful planning and tactical movement
// - **Time Challenge**: Race against the clock to make your delivery as quickly as possible
// - **Permanent Obstacles**: Strategic permanent walls in the center and diagonals create consistent navigation challenges

// ## Gameplay
// Players start at the top of the maze and must navigate to a moving delivery point at the bottom. The challenge comes from:
// - Planning routes through shifting walls
// - Tracking the moving delivery point
// - Making precise movements under pressure
// - Finding the optimal path while the maze changes around you

// ## Technical Highlights
// - Built with Ebitengine for smooth 2D graphics
// - Implements smart pathfinding to prevent player entrapment
// - Efficient maze generation and update algorithms
// - Clean, modular code design for easy maintenance and future enhancements

// ## How to Play
// - Use arrow keys or WASD to move
// - Each key press moves one cell
// - Green square marks the start
// - Blue square marks the moving delivery point
// - Timer starts when you enter the maze
// - Press ESC to exit at any time

// ## Command Line Usage
// ```bash
// # Start the game using any of these commands:
// go-games delivery-dash
// go-games dd
// ```

// ## Development Notes
// - Generated with AI assistance from Cursor on April 24, 2025
// - Collaborative development with @emmahsax
// - Focused on balancing challenge and fairness
// - Implements robust wall protection system to prevent soft-locks

// Perfect for players who enjoy dynamic puzzle-solving and quick-thinking challenges! 🎮

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/spf13/cobra"
)

const (
	cellSize           = 40                          // Size of each cell in the maze
	mazeWidth          = 20                          // Number of cells wide
	mazeHeight         = 15                          // Number of cells tall
	wallThickness      = 2                           // Thickness of maze walls
	borderSize         = cellSize                    // Size of the border around the maze
	screenWidth        = cellSize * (mazeWidth + 2)  // Add 2 cells for borders
	screenHeight       = cellSize * (mazeHeight + 2) // Add 2 cells for borders
	mazeUpdateInterval = 0.05                        // Seconds between maze updates (20 times per second)
	wallUpdateInterval = 0.25                        // Seconds between wall updates (4 times per second)
	moveCooldown       = 10                          // Frames between allowed movements
	wallChangeChance   = 0.75                        // Chance for each wall to change (75%)
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Game struct {
	car            *Car
	maze           [][]bool // true for walls, false for paths
	startX, startY int
	endX, endY     int
	lastMazeUpdate time.Time
	gameOver       bool
	win            bool
	moveTimer      int                 // Counter for movement cooldown
	startTime      time.Time           // When the player started moving
	hasStarted     bool                // Whether the player has left the start position
	finalTime      time.Duration       // Time taken to complete the maze
	lastWallUpdate time.Time           // Last time walls were updated
	lastKeyState   map[ebiten.Key]bool // Track last key press state
	titleScreen    bool                // Whether to show the title screen
}

type Car struct {
	x, y         float64
	direction    Direction
	sprite       *ebiten.Image
	cellX, cellY int     // Current cell position
	rotation     float64 // Rotation in degrees
}

func NewGame() *Game {
	// Create a car sprite (a simple car shape for now)
	carSprite := ebiten.NewImage(30, 20)
	// Draw a simple car shape
	carSprite.Fill(color.RGBA{0, 0, 0, 0}) // Clear the image
	// Draw car body
	vector.DrawFilledRect(carSprite, 0, 5, 30, 10, color.RGBA{255, 0, 0, 255}, false)
	// Draw windows
	vector.DrawFilledRect(carSprite, 5, 2, 8, 3, color.RGBA{200, 200, 255, 255}, false)
	vector.DrawFilledRect(carSprite, 17, 2, 8, 3, color.RGBA{200, 200, 255, 255}, false)
	// Draw wheels
	vector.DrawFilledRect(carSprite, 3, 0, 4, 5, color.RGBA{50, 50, 50, 255}, false)
	vector.DrawFilledRect(carSprite, 23, 0, 4, 5, color.RGBA{50, 50, 50, 255}, false)
	vector.DrawFilledRect(carSprite, 3, 15, 4, 5, color.RGBA{50, 50, 50, 255}, false)
	vector.DrawFilledRect(carSprite, 23, 15, 4, 5, color.RGBA{50, 50, 50, 255}, false)

	// Initialize maze
	maze := make([][]bool, mazeHeight)
	for i := range maze {
		maze[i] = make([]bool, mazeWidth)
	}

	// Set start and end positions outside the maze
	startX := mazeWidth / 2
	startY := -1 // One cell above the maze
	endX := mazeWidth / 2
	endY := mazeHeight // One cell below the maze

	// Create initial maze layout
	generateMaze(maze, startX, 0, endX, mazeHeight-1) // Adjust path generation to connect to borders

	// Create car at start position
	car := &Car{
		x:         float64((startX+1)*cellSize + cellSize/2), // +1 for border
		y:         float64((startY+1)*cellSize + cellSize/2), // +1 for border
		direction: Down,
		sprite:    carSprite,
		cellX:     startX,
		cellY:     startY,
		rotation:  0, // Start facing down (0 degrees)
	}

	// Initialize timers
	now := time.Now()
	return &Game{
		car:            car,
		maze:           maze,
		startX:         startX,
		startY:         startY,
		endX:           endX,
		endY:           endY,
		lastMazeUpdate: now,
		lastWallUpdate: now,
		moveTimer:      0,
		hasStarted:     false,
		lastKeyState:   make(map[ebiten.Key]bool),
		titleScreen:    true, // Start with title screen
	}
}

func generateMaze(maze [][]bool, startX, startY, endX, endY int) {
	// Initialize all cells as paths
	for y := range maze {
		for x := range maze[y] {
			maze[y][x] = false
		}
	}

	// Add random permanent walls
	// Create a random pattern of walls in the middle section
	centerX := mazeWidth / 2
	centerY := mazeHeight / 2

	// Add random walls in a 5x5 area around the center
	for y := centerY - 2; y <= centerY+2; y++ {
		for x := centerX - 2; x <= centerX+2; x++ {
			if rand.Float32() < 0.6 { // 60% chance of wall
				maze[y][x] = true
			}
		}
	}

	// Add some random diagonal walls
	for i := 0; i < mazeHeight/2; i++ {
		if rand.Float32() < 0.7 { // 70% chance of wall
			maze[i][i] = true
			maze[i][i+1] = true
		}
		if rand.Float32() < 0.7 { // 70% chance of wall
			maze[i][mazeWidth-1-i] = true
			maze[i][mazeWidth-2-i] = true
		}
	}

	// Create a single path from start to end
	currentX, currentY := startX, startY
	maze[currentY][currentX] = false

	// Create a path to the end
	for currentX != endX || currentY != endY {
		// Determine possible moves
		var possibleMoves []Direction
		if currentY > 0 && !maze[currentY-1][currentX] {
			possibleMoves = append(possibleMoves, Up)
		}
		if currentX < mazeWidth-1 && !maze[currentY][currentX+1] {
			possibleMoves = append(possibleMoves, Right)
		}
		if currentY < mazeHeight-1 && !maze[currentY+1][currentX] {
			possibleMoves = append(possibleMoves, Down)
		}
		if currentX > 0 && !maze[currentY][currentX-1] {
			possibleMoves = append(possibleMoves, Left)
		}

		// If no moves are possible, backtrack
		if len(possibleMoves) == 0 {
			break
		}

		// Choose a random move
		move := possibleMoves[rand.Intn(len(possibleMoves))]
		switch move {
		case Up:
			currentY--
		case Right:
			currentX++
		case Down:
			currentY++
		case Left:
			currentX--
		}
		maze[currentY][currentX] = false
	}

	// Add some random walls (only 25% of the cells)
	totalCells := mazeWidth * mazeHeight
	wallsToAdd := totalCells / 4
	for i := 0; i < wallsToAdd; i++ {
		x := rand.Intn(mazeWidth)
		y := rand.Intn(mazeHeight)
		// Protect only the entrance and exit cells
		if (x == startX && y == 0) || (x == endX && y == mazeHeight-1) {
			continue
		}
		maze[y][x] = true
	}
}

func (g *Game) Update() error {
	// Check for escape key to exit (always check this first)
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Handle title screen
	if g.titleScreen {
		// Only accept Space or Enter to start
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.titleScreen = false
			return nil
		}
		return nil
	}

	if g.gameOver || g.win {
		return nil
	}

	// Update movement cooldown
	if g.moveTimer > 0 {
		g.moveTimer--
	}

	// Check for maze updates
	currentTime := time.Now()
	timeSinceLastUpdate := currentTime.Sub(g.lastMazeUpdate)

	if timeSinceLastUpdate >= time.Duration(float64(mazeUpdateInterval)*float64(time.Second)) {
		// Try to move the end position multiple times
		for i := 0; i < 2; i++ { // Try to move up to 2 times per update
			oldEndX := g.endX
			// Try to jump to a random position at the bottom
			newEndX := rand.Intn(mazeWidth)
			g.endX = newEndX
			// If the new position would trap the player, revert the change
			if g.wouldTrapPlayer(g.endX, mazeHeight-1) {
				g.endX = oldEndX
			}
		}

		g.lastMazeUpdate = currentTime
	}

	// Check for wall updates (separate from end position updates)
	timeSinceWallUpdate := currentTime.Sub(g.lastWallUpdate)
	if timeSinceWallUpdate >= time.Duration(float64(wallUpdateInterval)*float64(time.Second)) {
		// Update only 25% of the walls, ensuring player is never trapped
		wallsChanged := 0
		for y := range g.maze {
			for x := range g.maze[y] {
				// Protect only the entrance, exit, and permanent walls
				centerX := mazeWidth / 2
				centerY := mazeHeight / 2
				isProtected := (x == g.startX && y == 0) || // Start position
					(x == g.endX && y == mazeHeight-1) || // End position
					(y < mazeHeight/2 && (x == y || x == y+1)) || // First diagonal
					(y < mazeHeight/2 && (x == mazeWidth-1-y || x == mazeWidth-2-y)) || // Second diagonal
					(x == centerX && y >= centerY-1 && y <= centerY+1) || // Center vertical
					(y == centerY && x >= centerX-1 && x <= centerX+1) // Center horizontal

				if !isProtected && rand.Float32() < wallChangeChance {
					// Temporarily store the current state
					oldState := g.maze[y][x]
					// Try the change
					g.maze[y][x] = !g.maze[y][x]
					// If it would trap the player, revert the change
					if g.wouldTrapPlayer(x, y) {
						g.maze[y][x] = oldState
					} else {
						wallsChanged++
					}
				}
			}
		}
		g.lastWallUpdate = currentTime
	}

	// Handle car movement with improved key detection
	if g.moveTimer == 0 {
		// Check for key press transitions (key just pressed)
		if (ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && !g.lastKeyState[ebiten.KeyUp] && !g.lastKeyState[ebiten.KeyW] {
			g.moveCar(Up)
			g.moveTimer = moveCooldown
		}
		if (ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)) && !g.lastKeyState[ebiten.KeyRight] && !g.lastKeyState[ebiten.KeyD] {
			g.moveCar(Right)
			g.moveTimer = moveCooldown
		}
		if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && !g.lastKeyState[ebiten.KeyDown] && !g.lastKeyState[ebiten.KeyS] {
			g.moveCar(Down)
			g.moveTimer = moveCooldown
		}
		if (ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)) && !g.lastKeyState[ebiten.KeyLeft] && !g.lastKeyState[ebiten.KeyA] {
			g.moveCar(Left)
			g.moveTimer = moveCooldown
		}
	}

	// Update last key states
	g.lastKeyState[ebiten.KeyUp] = ebiten.IsKeyPressed(ebiten.KeyUp)
	g.lastKeyState[ebiten.KeyW] = ebiten.IsKeyPressed(ebiten.KeyW)
	g.lastKeyState[ebiten.KeyRight] = ebiten.IsKeyPressed(ebiten.KeyRight)
	g.lastKeyState[ebiten.KeyD] = ebiten.IsKeyPressed(ebiten.KeyD)
	g.lastKeyState[ebiten.KeyDown] = ebiten.IsKeyPressed(ebiten.KeyDown)
	g.lastKeyState[ebiten.KeyS] = ebiten.IsKeyPressed(ebiten.KeyS)
	g.lastKeyState[ebiten.KeyLeft] = ebiten.IsKeyPressed(ebiten.KeyLeft)
	g.lastKeyState[ebiten.KeyA] = ebiten.IsKeyPressed(ebiten.KeyA)

	// Check for win condition
	if g.car.cellY == mazeHeight && g.car.cellX == g.endX {
		g.win = true
		g.finalTime = time.Since(g.startTime)
	}

	return nil
}

func (g *Game) wouldTrapPlayer(x, y int) bool {
	// Create a temporary copy of the maze
	tempMaze := make([][]bool, len(g.maze))
	for i := range g.maze {
		tempMaze[i] = make([]bool, len(g.maze[i]))
		copy(tempMaze[i], g.maze[i])
	}

	// Check if there's a path from player to end
	return !g.hasPathToEnd(tempMaze, g.car.cellX, g.car.cellY)
}

func (g *Game) hasPathToEnd(maze [][]bool, startX, startY int) bool {
	visited := make([][]bool, len(maze))
	for i := range maze {
		visited[i] = make([]bool, len(maze[i]))
	}

	var dfs func(x, y int) bool
	dfs = func(x, y int) bool {
		if x < 0 || x >= mazeWidth || y < 0 || y >= mazeHeight || visited[y][x] || maze[y][x] {
			return false
		}
		if x == g.endX && y == mazeHeight-1 {
			return true
		}
		visited[y][x] = true
		return dfs(x+1, y) || dfs(x-1, y) || dfs(x, y+1) || dfs(x, y-1)
	}

	return dfs(startX, startY)
}

func (g *Game) moveCar(dir Direction) {
	// Calculate new cell position
	newCellX, newCellY := g.car.cellX, g.car.cellY

	// Calculate the target position based on direction
	switch dir {
	case Up:
		newCellY--
		g.car.rotation = 180 // Face up (180 degrees from down)
	case Right:
		newCellX++
		g.car.rotation = 270 // Face right (270 degrees from down)
	case Down:
		newCellY++
		g.car.rotation = 0 // Face down (0 degrees)
	case Left:
		newCellX--
		g.car.rotation = 90 // Face left (90 degrees from down)
	}

	// Special case for start position (above maze)
	if g.car.cellY == -1 {
		if newCellY == 0 && !g.maze[0][newCellX] {
			// Allow movement into the maze
			g.car.cellX = newCellX
			g.car.cellY = newCellY
			g.car.direction = dir
			g.car.x = float64((newCellX+1)*cellSize + cellSize/2)
			g.car.y = float64((newCellY+1)*cellSize + cellSize/2)
			// Start the timer when leaving the start position
			if !g.hasStarted {
				g.hasStarted = true
				g.startTime = time.Now()
			}
		}
		return
	}

	// Special case for end position (below maze)
	if g.car.cellY == mazeHeight-1 && dir == Down && newCellX == g.endX {
		// Allow movement to the end position
		g.car.cellX = newCellX
		g.car.cellY = newCellY
		g.car.direction = dir
		g.car.x = float64((newCellX+1)*cellSize + cellSize/2)
		g.car.y = float64((newCellY+1)*cellSize + cellSize/2)
		return
	}

	// Normal maze movement
	if newCellX >= 0 && newCellX < mazeWidth &&
		newCellY >= 0 && newCellY < mazeHeight &&
		!g.maze[newCellY][newCellX] {
		// Update car position
		g.car.cellX = newCellX
		g.car.cellY = newCellY
		g.car.direction = dir
		// Center the car in the new cell
		g.car.x = float64((newCellX+1)*cellSize + cellSize/2)
		g.car.y = float64((newCellY+1)*cellSize + cellSize/2)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{50, 50, 50, 255})

	if g.titleScreen {
		// Draw title screen
		title := "DELIVERY DASH"
		scenario := "You are a delivery driver in a chaotic, ever-changing city.\n" +
			"The walls are shifting beneath your wheels and the customer keeps moving!\n" +
			"Navigate through the maze to deliver your package as fast\n" +
			"as possible, before the walls trap you!\n\n" +
			"Use arrow keys or WASD to move.\n\n" +
			"Press SPACE or ENTER to start, ESC to exit"

		// Draw title
		ebitenutil.DebugPrintAt(screen, title, screenWidth/2-100, screenHeight/2-100)

		// Draw scenario text (split into lines)
		lines := strings.Split(scenario, "\n")
		for i, line := range lines {
			ebitenutil.DebugPrintAt(screen, line, screenWidth/2-200, screenHeight/2-50+i*20)
		}

		return
	}

	// Draw border around the maze
	vector.DrawFilledRect(screen,
		float32(cellSize),
		float32(cellSize),
		float32(mazeWidth*cellSize),
		float32(mazeHeight*cellSize),
		color.RGBA{30, 30, 30, 255},
		false,
	)

	// Draw maze walls as thin lines
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] {
				// Draw a cross of lines for wall cells
				// Vertical line in the middle of the cell
				vector.DrawFilledRect(screen,
					float32((x+1)*cellSize+cellSize/2), // +1 for border
					float32((y+1)*cellSize),            // +1 for border
					float32(wallThickness),
					float32(cellSize),
					color.RGBA{100, 100, 100, 255},
					false,
				)
				// Horizontal line in the middle of the cell
				vector.DrawFilledRect(screen,
					float32((x+1)*cellSize),            // +1 for border
					float32((y+1)*cellSize+cellSize/2), // +1 for border
					float32(cellSize),
					float32(wallThickness),
					color.RGBA{100, 100, 100, 255},
					false,
				)
			} else {
				// Optional: very subtle path indicator
				vector.DrawFilledRect(screen,
					float32((x+1)*cellSize+cellSize/2-1), // +1 for border
					float32((y+1)*cellSize+cellSize/2-1), // +1 for border
					2,
					2,
					color.RGBA{60, 60, 60, 255},
					false,
				)
			}
		}
	}

	// Draw start and end positions
	vector.DrawFilledRect(screen,
		float32((g.startX+1)*cellSize), // +1 for border
		float32((g.startY+1)*cellSize), // +1 for border
		float32(cellSize),
		float32(cellSize),
		color.RGBA{0, 255, 0, 255},
		false,
	)
	vector.DrawFilledRect(screen,
		float32((g.endX+1)*cellSize), // +1 for border
		float32((g.endY+1)*cellSize), // +1 for border
		float32(cellSize),
		float32(cellSize),
		color.RGBA{0, 0, 255, 255},
		false,
	)

	// Draw car with rotation
	op := &ebiten.DrawImageOptions{}
	// Set the rotation center to the middle of the car
	op.GeoM.Translate(-15, -10)                             // Move to center
	op.GeoM.Rotate(float64(g.car.rotation) * math.Pi / 180) // Rotate
	op.GeoM.Translate(g.car.x, g.car.y)                     // Move to position
	screen.DrawImage(g.car.sprite, op)

	// Draw game over or win message
	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over - Press ESC to exit")
	} else if g.win {
		// Show final time
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Total Time: %.2f seconds - Press ESC to exit", g.finalTime.Seconds()))
	} else if g.hasStarted {
		// Show current time while playing
		elapsed := time.Since(g.startTime)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Time: %.2f - Press ESC to exit", elapsed.Seconds()))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delivery-dash",
		Aliases: []string{"dd"},
		Short:   "Escape the chaotic, ever-changing maze to deliver your package to the customer! (alias: dd)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ebiten.SetWindowSize(screenWidth, screenHeight)
			ebiten.SetWindowTitle("Delivery Dash")

			if err := ebiten.RunGame(NewGame()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
