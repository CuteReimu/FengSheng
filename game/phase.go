package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
)

type DrawPhase struct {
	player interfaces.IPlayer
}

func (dp *DrawPhase) Resolve() (finished bool) {
	game := dp.player.GetGame()
	if !dp.player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	logger.Info(dp.player, "的回合开始了")
	for _, p := range game.GetPlayers() {
		p.NotifyDrawPhase()
	}
	dp.player.Draw(3)
	game.PushResolveStackNode(&MainPhaseIdle{player: dp.player})
	return true
}

type MainPhaseIdle struct {
	player interfaces.IPlayer
}

func (mp *MainPhaseIdle) Resolve() (finished bool) {
	game := mp.player.GetGame()
	if !mp.player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	for _, p := range game.GetPlayers() {
		p.NotifyMainPhase(30)
	}
	return false
}

type SendPhaseStart struct {
	player interfaces.IPlayer
}

func (sp *SendPhaseStart) Resolve() (finished bool) {
	game := sp.player.GetGame()
	if sp.player.IsAlive() {
		if len(sp.player.GetCards()) == 0 {
			logger.Info(sp.player, "没有情报可传，输掉了游戏")
			game.GetDeck().Discard(sp.player.DeleteAllMessageCards()...)
			sp.player.SetLose(true)
			sp.player.SetAlive(false)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(sp.player.Location(), true)
			}
		}
	}
	if !sp.player.IsAlive() {
		game.PushResolveStackNode(&NextTurn{player: sp.player})
		return true
	}
	for _, p := range game.GetPlayers() {
		p.NotifySendPhaseStart(20)
	}
	return false
}

type OnSendCard struct {
	WhoseTurn     interfaces.IPlayer
	MessageCard   interfaces.ICard
	Dir           protos.Direction
	TargetPlayer  interfaces.IPlayer
	LockedPlayers []interfaces.IPlayer
}

func (sc *OnSendCard) Resolve() (finished bool) {
	game := sc.WhoseTurn.GetGame()
	logger.Info(sc.WhoseTurn, "传出了", sc.MessageCard, "，方向是", protos.Direction_name[int32(sc.Dir)], "传给了", sc.TargetPlayer)
	sc.WhoseTurn.DeleteCard(sc.MessageCard.GetId())
	game.PushResolveStackNode(&SendPhaseIdle{
		WhoseTurn:     sc.WhoseTurn,
		MessageCard:   sc.MessageCard,
		Dir:           sc.Dir,
		InFrontOfWhom: sc.TargetPlayer,
		LockedPlayers: sc.LockedPlayers,
	})
	for _, p := range game.GetPlayers() {
		p.NotifySendMessageCard()
	}
	return true
}

type SendPhaseIdle struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	Dir                 protos.Direction
	InFrontOfWhom       interfaces.IPlayer
	LockedPlayers       []interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (sp *SendPhaseIdle) Resolve() (finished bool) {
	game := sp.WhoseTurn.GetGame()
	logger.Info("情报到达", sp.InFrontOfWhom, "面前")
	for _, p := range game.GetPlayers() {
		p.NotifySendPhase(20)
	}
	return false
}

type MessageMoveNext struct {
	SendPhase *SendPhaseIdle
}

func (mm *MessageMoveNext) Resolve() (finished bool) {
	game := mm.SendPhase.WhoseTurn.GetGame()
	if mm.SendPhase.Dir == protos.Direction_Up {
		if mm.SendPhase.WhoseTurn.IsAlive() {
			mm.SendPhase.InFrontOfWhom = mm.SendPhase.WhoseTurn
			game.PushResolveStackNode(mm.SendPhase)
		} else {
			game.PushResolveStackNode(&NextTurn{player: mm.SendPhase.WhoseTurn})
		}
	} else {
		inFrontOfWhom := mm.SendPhase.InFrontOfWhom.Location()
		for {
			if mm.SendPhase.Dir == protos.Direction_Left {
				inFrontOfWhom = (inFrontOfWhom + len(game.GetPlayers()) - 1) % len(game.GetPlayers())
			} else {
				inFrontOfWhom = (inFrontOfWhom + 1) % len(game.GetPlayers())
			}
			mm.SendPhase.InFrontOfWhom = game.GetPlayers()[inFrontOfWhom]
			if mm.SendPhase.InFrontOfWhom.IsAlive() {
				game.PushResolveStackNode(mm.SendPhase)
				break
			} else if mm.SendPhase.WhoseTurn == mm.SendPhase.InFrontOfWhom {
				game.PushResolveStackNode(&NextTurn{player: mm.SendPhase.WhoseTurn})
				break
			}
		}
	}
	return true
}

type OnChooseReceiveCard struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	InFrontOfWhom       interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (cr *OnChooseReceiveCard) Resolve() (finished bool) {
	game := cr.WhoseTurn.GetGame()
	logger.Info(cr.InFrontOfWhom, "选择接收情报")
	game.PushResolveStackNode(&FightPhaseIdle{
		WhoseTurn:           cr.WhoseTurn,
		MessageCard:         cr.MessageCard,
		InFrontOfWhom:       cr.InFrontOfWhom,
		WhoseFightTurn:      cr.InFrontOfWhom,
		IsMessageCardFaceUp: cr.IsMessageCardFaceUp,
	})
	for _, p := range game.GetPlayers() {
		p.NotifyChooseReceiveCard()
	}
	return true
}

type FightPhaseIdle struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	InFrontOfWhom       interfaces.IPlayer
	WhoseFightTurn      interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (cr *FightPhaseIdle) Resolve() (finished bool) {
	game := cr.WhoseTurn.GetGame()
	for _, p := range game.GetPlayers() {
		p.NotifyFightPhase(20)
	}
	return false
}

type NextTurn struct {
	player interfaces.IPlayer
}

func (n *NextTurn) Resolve() (finished bool) {
	game := n.player.GetGame()
	game.SetCurrentMessageCard(nil)
	whoseTurn := n.player.Location()
	for {
		whoseTurn = (whoseTurn + 1) % len(game.GetPlayers())
		player := game.GetPlayers()[whoseTurn]
		if player.IsAlive() {
			game.PushResolveStackNode(&DrawPhase{player: player})
			return true
		}
	}
}
