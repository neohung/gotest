// Package actor
package actor

import (
	"test2/package/framebuffer"
	"test2/package/world"

	"github.com/gdamore/tcell/v2"
)

type (
	ActorID   int
	ActorKind int
)

const (
	ActorPlayer ActorKind = iota
	ActorMonster
	ActorBOSS
	ActorNPC
	ActorItem
)

// ===============================
// Action
type ActionType int

const (
	ActionMove ActionType = iota
	ActionAttack
	ActionPickup
)

type Action struct {
	Type   ActionType
	DX, DY int
	Target *Actor
}

//===================================

type AIController interface {
	Update(a *Actor, w *world.World)
}

type Actor struct {
	ID   ActorID
	Kind ActorKind
	// ---- Position ----
	X, Y int
	// ---- Rendering ----
	CH    rune
	Color tcell.Color
	Z     int // layer order（玩家 > 怪物 > 地圖）
	// ---- Visibility ----
	BlocksMove  bool
	BlocksSight bool

	// ---- Gameplay ----
	HP    int
	MaxHP int

	// ---- Behavior ----
	AI AIController // nil for player / items

	// ---- State flags ----
	Alive bool

	FOV int
}

func RenderPlayer(fb *framebuffer.Framebuffer, player *Actor, w *world.World) {
	sx := player.X - w.Cam.X
	sy := player.Y - w.Cam.Y

	if sx >= 0 && sx < w.Cam.W && sy >= 0 && sy < w.Cam.H {
		fb.View[sy][sx] = player.CH
	}
}

// RenderActorsLayer render to actorFB
func RenderActorsLayer(fb *framebuffer.Framebuffer, actors []*Actor, w *world.World) {
	for _, a := range actors {
		if !a.Alive {
			continue
		}
		if !w.Visible[a.Y][a.X] {
			continue
		}
		sx := a.X - w.Cam.X
		sy := a.Y - w.Cam.Y

		if sx < 0 || sy < 0 || sx >= fb.W || sy >= fb.H {
			continue
		}
		fb.View[sy][sx] = a.CH
	}
}
