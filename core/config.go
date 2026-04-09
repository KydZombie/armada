package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Config struct {
	ScreenWidth, ScreenHeight int32
	VSync                     bool
	Resizeable                bool

	Debug bool
}

func (c *Config) updateWindow() {
	rl.SetWindowSize(int(c.ScreenWidth), int(c.ScreenHeight))
}

func (c *Config) SetScreenWidth(s int32) {
	c.ScreenWidth = s
	c.updateWindow()
}

func (c *Config) SetScreenHeight(s int32) {
	c.ScreenHeight = s
	c.updateWindow()
}
