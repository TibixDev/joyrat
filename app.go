package main

import (
	"fmt"
	"math"

	"github.com/go-vgo/robotgo"
	"github.com/veandco/go-sdl2/sdl"
)

type Stick struct {
	X int16
	Y int16
}

type Settings struct {
	MOUSE_SPEED       int
	MOUSE_SPEED_LOW   int
	MOUSE_SPEED_HIGH  int
	SCROLL_SPEED      int
	JOYSTICK_DEADZONE int16
	AXIS_LT           uint8
	AXIS_RT           uint8
	AXIS_LS_X         uint8
	AXIS_LS_Y         uint8
	AXIS_RS_X         uint8
	AXIS_RS_Y         uint8
}

// Constants
const I16_MIN = -32768
const LEFT_MOUSE_BUTTON = 0
const RIGHT_MOUSE_BUTTON = 1

var DefaultCfg = Settings{
	MOUSE_SPEED:       10,
	MOUSE_SPEED_LOW:   3,
	MOUSE_SPEED_HIGH:  20,
	SCROLL_SPEED:      1,
	JOYSTICK_DEADZONE: 4000,
	AXIS_LT:           2,
	AXIS_RT:           5,
	AXIS_LS_X:         0,
	AXIS_LS_Y:         1,
	AXIS_RS_X:         3,
	AXIS_RS_Y:         4,
}

// Settings
var cfg = Settings{
	MOUSE_SPEED:       10,
	MOUSE_SPEED_LOW:   3,
	MOUSE_SPEED_HIGH:  20,
	SCROLL_SPEED:      1,
	JOYSTICK_DEADZONE: 4000,
	AXIS_LT:           2,
	AXIS_RT:           5,
	AXIS_LS_X:         0,
	AXIS_LS_Y:         1,
	AXIS_RS_X:         3,
	AXIS_RS_Y:         4,
}

// Variables
var stick Stick = Stick{0, 0}
var scrollStick Stick = Stick{0, 0}
var speed = cfg.MOUSE_SPEED
var running = true

func main() {
	fmt.Println("Joyrat init!")

	// Get pid
	pid := robotgo.GetPid()
	fmt.Println("Joyrat PID:", pid)

	// Load config
	LoadCfg(&cfg)

	// Init SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic("Failed to initialize SDL")
	}
	defer sdl.Quit()

	// SDL Version
	ver := sdl.Version{}
	sdl.GetVersion(&ver)
	fmt.Printf("SDL Version: %d.%d.%d\n", ver.Major, ver.Minor, ver.Patch)

	// Enumerate joysticks
	joystickCount := sdl.NumJoysticks()
	fmt.Printf("Joystick count: %d\n", joystickCount)
	if joystickCount < 1 {
		fmt.Println("No joystick found")
		return
	}

	// Open joystick
	joystick := sdl.JoystickOpen(0)
	if joystick == nil {
		fmt.Println("Failed to open joystick")
		return
	}
	// defer joystick.Close()

	// Start mouse mover
	go mouseMover()

	// Start scroll mover
	go scrollMover()

	// Wait for events
	sdlEventLoop := func() {
		for running {
			event := sdl.WaitEvent()
			switch e := event.(type) {
			case *sdl.JoyAxisEvent:
				// fmt.Printf("Axis: %d, Value: %d\n", e.Axis, e.Value)

				// Update move stick
				if e.Axis == cfg.AXIS_LS_X {
					stick.X = e.Value
				} else if e.Axis == cfg.AXIS_LS_Y {
					stick.Y = e.Value
				}

				// Update scroll stick
				if e.Axis == cfg.AXIS_RS_Y {
					scrollStick.Y = e.Value
				} else if e.Axis == cfg.AXIS_RS_X {
					scrollStick.X = e.Value
				}

				// Handle triggers
				if e.Axis == cfg.AXIS_LT {
					if e.Value > I16_MIN {
						speed = cfg.MOUSE_SPEED_LOW
					} else {
						speed = cfg.MOUSE_SPEED
					}
				} else if e.Axis == cfg.AXIS_RT {
					if e.Value > I16_MIN {
						speed = cfg.MOUSE_SPEED_HIGH
					} else {
						speed = cfg.MOUSE_SPEED
					}
				}

			case *sdl.JoyButtonEvent:
				// fmt.Printf("Button: %d, State: %d\n", e.Button, e.State)

				// Handle left mouse button
				if e.Button == LEFT_MOUSE_BUTTON {
					if e.State == 1 {
						robotgo.MouseDown("left")
					} else {
						robotgo.MouseUp("left")
					}
				}

				// Handle right mouse button
				if e.Button == RIGHT_MOUSE_BUTTON {
					if e.State == 1 {
						robotgo.MouseDown("right")
					} else {
						robotgo.MouseUp("right")
					}
				}
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				return
			}
		}
	}

	// Start SDL event loop
	go sdlEventLoop()

	// Create GUI
	_, w := CreateGui(&cfg)

	// Set close event
	w.SetOnClosed(func() {
		running = false
		joystick.Close()
		sdl.Quit()
		fmt.Println("Joyrat is exiting...")
	})

	w.ShowAndRun()
}

// SECTION: INPUT/OUTPUT

// Clamp the stick value to -1.0 to 1.0 range
func clamp16(value int16) float64 {
	return float64(value) / 32767.0
}

// Calculate move for mouse or scrollwheel
func calculateMove(_stick Stick, _speed int) (int, int) {
	// Clamp
	x := clamp16(_stick.X)
	y := clamp16(_stick.Y)

	// Normalize the clamped values (if necessary)
	magnitude := math.Sqrt(x*x + y*y)
	if magnitude > 1 {
		x /= magnitude
		y /= magnitude
	}

	// Return
	return int(x * float64(_speed)), int(y * float64(_speed))
}

// Move the mouse continuously if the L stick is not in the center
func mouseMover() {
	for running {
		// Detect if stick is in the center
		if stick.X == 0 && stick.Y == 0 {
			robotgo.MilliSleep(10)
			continue
		}

		// Detect if stick is in the deadzone
		isXInDeadzone := stick.X < cfg.JOYSTICK_DEADZONE && stick.X > -cfg.JOYSTICK_DEADZONE
		isYInDeadzone := stick.Y < cfg.JOYSTICK_DEADZONE && stick.Y > -cfg.JOYSTICK_DEADZONE
		if isXInDeadzone && isYInDeadzone {
			robotgo.MilliSleep(10)
			continue
		}

		// Calculate move
		x, y := calculateMove(stick, speed)

		// Move mouse
		mx, my := robotgo.Location()
		robotgo.Move(x+mx, y+my)
		fmt.Println("Moving", x, y)

		// Sleep
		robotgo.MilliSleep(10)
	}
}

// Scroll the mouse continuously if the R stick is not in the center
func scrollMover() {
	for running {
		// Detect if stick is in the center
		if scrollStick.X == 0 && scrollStick.Y == 0 {
			robotgo.MilliSleep(10)
			continue
		}

		// Detect if stick is in the deadzone
		isXInDeadzone := scrollStick.X < cfg.JOYSTICK_DEADZONE && scrollStick.X > -cfg.JOYSTICK_DEADZONE
		isYInDeadzone := scrollStick.Y < cfg.JOYSTICK_DEADZONE && scrollStick.Y > -cfg.JOYSTICK_DEADZONE
		if isXInDeadzone && isYInDeadzone {
			robotgo.MilliSleep(10)
			continue
		}

		// Calculate move
		x, y := calculateMove(scrollStick, cfg.SCROLL_SPEED)

		// Scroll mouse
		robotgo.Scroll(-x, -y, 100)
		fmt.Println("Scrolling", -x, -y)

		// Sleep
		robotgo.MilliSleep(10)
	}
}
