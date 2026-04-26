package core

import rl "github.com/gen2brain/raylib-go/raylib"

type Config struct {
	ScreenWidth, ScreenHeight int32
	VSync                     bool
	Fullscreen                bool
	Resizeable                bool

	MasterVolume float32
	MusicVolume  float32
	SFXVolume    float32

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

func (c *Config) GetEffectiveMusicVolume() float32 {
	return c.MasterVolume * c.MusicVolume
}

func (c *Config) GetEffectiveSFXVolume() float32 {
	return c.MasterVolume * c.SFXVolume
}
