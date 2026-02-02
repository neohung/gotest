// Package game
package game

import (
	"math/rand"

	"test2/package/actor"
	"test2/package/world"
)

type Game struct {
	World  *world.World
	Actors []*actor.Actor
	Player *actor.Actor
}

func New(w *world.World, player *actor.Actor) *Game {
	g := Game{
		World:  w,
		Player: player,
	}
	return &g
}

func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func MonsterAI(m *actor.Actor, player *actor.Actor) actor.Action {
	dx := sign(player.X - m.X)
	dy := sign(player.Y - m.Y)
	// 看的見就追
	if abs(dx)+abs(dy) <= 5 {
		return actor.Action{
			Type: actor.ActionMove,
			DX:   dx,
			DY:   dy,
		}
	}
	// 看不見就亂走
	return actor.Action{
		Type: actor.ActionMove,
		DX:   rand.Intn(3) - 1,
		DY:   rand.Intn(3) - 1,
	}
}

/* Move: TryMove(dx, dy)
   Attack: Attack()
-  Pickup: Pickup()
-  Wait: Die()
*/

func (g *Game) ApplyAction(a *actor.Actor, act actor.Action) {
	switch act.Type {
	case actor.ActionMove:
		g.TryMove(a, act.DX, act.DY)
	case actor.ActionAttack:
		g.Attack(a, act.Target)
	case actor.ActionPickup:
		g.Pickup(a)
	}
}

// Update Loop
func (g *Game) Update() {
	for _, a := range g.Actors {
		if a == g.Player {
			continue
		}
		// 定義a的Moster AI跟player的互動方式
		act := MonsterAI(a, g.Player)
		g.ApplyAction(a, act)

	}
}

func (g *Game) ActorAt(x, y int) *actor.Actor {
	for _, a := range g.Actors {
		if a.Alive && a.X == x && a.Y == y {
			return a
		}
	}
	return nil
}

func (g *Game) Pickup(a *actor.Actor) {
}

func (g *Game) RemoveActor(a *actor.Actor) {
	for i, other := range g.Actors {
		if other == a {
			// remove from slice
			g.Actors = append(g.Actors[:i], g.Actors[i+1:]...)
			break
		}
	}
	if a == g.Player {
		// Game over
		g.Player = nil
	}
}

func (g *Game) Attack(attacker, defender *actor.Actor) {
	defender.HP -= 1
	if defender.HP <= 0 {
		g.RemoveActor(defender)
	}
}

func (g *Game) TryMove(a *actor.Actor, dx, dy int) {
	nx := a.X + dx
	ny := a.Y + dy
	if g.World.IsBlocked(nx, ny) {
		// can not move
		return
	}
	if target := g.ActorAt(nx, ny); target != nil {
		// can attack
		g.Attack(a, target)
		return
	}
	a.X = nx
	a.Y = ny
}
