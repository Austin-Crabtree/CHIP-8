package main

import (
	"flag"
	"math/rand"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/massung/chip-8/chip8"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/sqweek/dialog"
)

var (
	/// True if the ROM should load paused.
	///
	Break bool

	/// True if pausing emulation (single stepping).
	///
	Paused bool

	/// The CHIP-8 virtual machine.
	///
	VM *chip8.CHIP_8

	/// The SDL Window and Renderer.
	///
	Window *sdl.Window
	Renderer *sdl.Renderer
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var err error

	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// parse the command line
	flag.BoolVar(&Break, "b", false, "Start ROM paused.")
	flag.Parse()

	// get the file name of the ROM to load
	file := flag.Arg(0)

	// initialize SDL or panic
	if err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		panic(err)
	}

	// create the main window and renderer or panic
	flags := sdl.WINDOW_OPENGL | sdl.WINDOWPOS_CENTERED
	if Window, Renderer, err = sdl.CreateWindowAndRenderer(614, 380, uint32(flags)); err != nil {
		panic(err)
	}

	// set the icon
	if icon, err := sdl.LoadBMP("data/icon.bmp"); err == nil {
		mask := sdl.MapRGB(icon.Format, 255, 0, 255)

		// create the mask color key and set the icon
		icon.SetColorKey(1, mask)
		Window.SetIcon(icon)
	}

	// set the title
	Window.SetTitle("CHIP-8")

	// show copyright information
	Log("CHIP-8, Copyright 2017 by Jeffrey Massung")
	Log("All rights reserved")

	// create a new CHIP-8 virtual machine
	Load(file)

	// initialize sub-systems
	InitScreen()
	InitAudio()
	InitFont()

	// initially break into debugger?
	Paused = Break

	// set processor speed and refresh rate
	clock := time.NewTicker(time.Millisecond * 2)
	video := time.NewTicker(time.Second / 60)

	// notify that the main loop has started
	Logln("Starting program; press 'H' for help")

	// loop until window closed or user quit
	for ProcessEvents() {
		select {
		case <-video.C:
			Refresh()
		case <-clock.C:
			res := VM.Process(Paused)

			switch res.(type) {
			case chip8.Breakpoint:
				Log()
				Log(res.Error())

				// break the emulation
				Paused = true
			}
		}
	}
}

func LoadDialog() {
	dlg := dialog.File().Title("Load ROM / C8 Assembler")

	// types of files to load
	dlg.Filter("All Files", "*")
	dlg.Filter("C8 Assembler Files", "c8", "chip8")
	dlg.Filter("ROMs", "rom", "")

	// try and load it
	if file, err := dlg.Load(); err == nil {
		Load(file)
	}
}

func Load(file string) {
	if file == "" {
		Logln("Loading PONG... ")
		VM, _ = chip8.LoadROM(chip8.Pong)
	} else {
		base := filepath.Base(file)

		// show the action being taken
		Logln("Loading", base)

		// is this a chip-8 assembly source file?
		if strings.ToUpper(filepath.Ext(base)) == ".C8" {
			asm, err := chip8.Assemble(file)
			if err != nil {
				Log(err.Error())
			}

			// even on error, the assembly is valid
			VM, _ = chip8.LoadAssembly(asm)
		} else {
			VM, _ = chip8.LoadFile(file)
		}
	}
}

func Refresh() {
	Renderer.SetDrawColor(32, 42, 53, 255)
	Renderer.Clear()

	// frame various portions of the app
	Frame(8, 8, 386, 194)
	Frame(8, 208, 386, 164)
	Frame(402, 8, 204, 194)
	Frame(402, 208, 204, 164)

	// update the video screen and copy it
	RefreshScreen()
	CopyScreen(10, 10, 384, 192)

	// debug assembly and virtual registers
	DebugLog(12, 212)
	DebugAssembly(406, 12)
	DebugRegisters(406, 212)

	// show the new frame
	Renderer.Present()
}

func Frame(x, y, w, h int) {
	Renderer.SetDrawColor(0, 0, 0, 255)
	Renderer.DrawLine(x, y, x + w, y)
	Renderer.DrawLine(x, y, x, y + h)

	// highlight
	Renderer.SetDrawColor(95, 112, 120, 255)
	Renderer.DrawLine(x + w, y, x + w, y + h)
	Renderer.DrawLine(x, y + h, x + w, y + h)
}
