package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
)

type DrawPhase struct {
	player interfaces.IPlayer
}

func (dp *DrawPhase) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := dp.player.GetGame()
	if !dp.player.IsAlive() {
		return &NextTurn{player: dp.player}, true
	}
	logger.Info(dp.player, "的回合开始了")
	for _, p := range game.GetPlayers() {
		p.NotifyDrawPhase()
	}
	dp.player.Draw(3)
	return &MainPhaseIdle{player: dp.player}, true
}

type MainPhaseIdle struct {
	player interfaces.IPlayer
}

func (mp *MainPhaseIdle) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := mp.player.GetGame()
	if !mp.player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	for _, p := range game.GetPlayers() {
		p.NotifyMainPhase(30)
	}
	return mp, false
}

type SendPhaseStart struct {
	player interfaces.IPlayer
}

func (sp *SendPhaseStart) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := sp.player.GetGame()
	if sp.player.IsAlive() {
		if len(sp.player.GetCards()) == 0 {
			logger.Info(sp.player, "没有情报可传，输掉了游戏")
			game.GetDeck().Discard(sp.player.DeleteAllMessageCards()...)
			sp.player.SetLose(true)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(sp.player.Location(), true)
			}
		}
	}
	if !sp.player.IsAlive() {
		return &NextTurn{player: sp.player}, true
	}
	for _, p := range game.GetPlayers() {
		p.NotifySendPhaseStart(20)
	}
	return sp, false
}

type OnSendCard struct {
	WhoseTurn     interfaces.IPlayer
	MessageCard   interfaces.ICard
	Dir           protos.Direction
	TargetPlayer  interfaces.IPlayer
	LockedPlayers []interfaces.IPlayer
}

func (sc *OnSendCard) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := sc.WhoseTurn.GetGame()
	logger.Info(sc.WhoseTurn, "传出了", sc.MessageCard, "，方向是", protos.Direction_name[int32(sc.Dir)], "传给了", sc.TargetPlayer)
	sc.WhoseTurn.DeleteCard(sc.MessageCard.GetId())
	for _, p := range game.GetPlayers() {
		p.NotifySendMessageCard()
	}
	return &SendPhaseIdle{
		WhoseTurn:     sc.WhoseTurn,
		MessageCard:   sc.MessageCard,
		Dir:           sc.Dir,
		InFrontOfWhom: sc.TargetPlayer,
		LockedPlayers: sc.LockedPlayers,
	}, true
}

type SendPhaseIdle struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	Dir                 protos.Direction
	InFrontOfWhom       interfaces.IPlayer
	LockedPlayers       []interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (sp *SendPhaseIdle) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := sp.WhoseTurn.GetGame()
	logger.Info("情报到达", sp.InFrontOfWhom, "面前")
	for _, p := range game.GetPlayers() {
		p.NotifySendPhase(20)
	}
	return sp, false
}

type MessageMoveNext struct {
	SendPhase *SendPhaseIdle
}

func (mm *MessageMoveNext) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := mm.SendPhase.WhoseTurn.GetGame()
	if mm.SendPhase.Dir == protos.Direction_Up {
		if mm.SendPhase.WhoseTurn.IsAlive() {
			mm.SendPhase.InFrontOfWhom = mm.SendPhase.WhoseTurn
			return mm.SendPhase, true
		} else {
			return &NextTurn{player: mm.SendPhase.WhoseTurn}, true
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
				return mm.SendPhase, true
			} else if mm.SendPhase.WhoseTurn == mm.SendPhase.InFrontOfWhom {
				return &NextTurn{player: mm.SendPhase.WhoseTurn}, true
			}
		}
	}
}

type OnChooseReceiveCard struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	InFrontOfWhom       interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (cr *OnChooseReceiveCard) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := cr.WhoseTurn.GetGame()
	logger.Info(cr.InFrontOfWhom, "选择接收情报")
	for _, p := range game.GetPlayers() {
		p.NotifyChooseReceiveCard()
	}
	return &FightPhaseIdle{
		WhoseTurn:           cr.WhoseTurn,
		MessageCard:         cr.MessageCard,
		InFrontOfWhom:       cr.InFrontOfWhom,
		WhoseFightTurn:      cr.InFrontOfWhom,
		IsMessageCardFaceUp: cr.IsMessageCardFaceUp,
	}, true
}

type FightPhaseIdle struct {
	WhoseTurn           interfaces.IPlayer
	MessageCard         interfaces.ICard
	InFrontOfWhom       interfaces.IPlayer
	WhoseFightTurn      interfaces.IPlayer
	IsMessageCardFaceUp bool
}

func (fp *FightPhaseIdle) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := fp.WhoseTurn.GetGame()
	for _, p := range game.GetPlayers() {
		p.NotifyFightPhase(20)
	}
	return fp, false
}

type NextTurn struct {
	player interfaces.IPlayer
}

func (nt *NextTurn) Resolve() (next interfaces.Fsm, continueResolve bool) {
	game := nt.player.GetGame()
	game.SetMessageCardFaceUp(false)
	game.SetCurrentMessageCard(nil)
	game.SetLockPlayers(nil)
	whoseTurn := nt.player.Location()
	for {
		whoseTurn = (whoseTurn + 1) % len(game.GetPlayers())
		player := game.GetPlayers()[whoseTurn]
		if player.IsAlive() {
			return &DrawPhase{player: player}, true
		}
	}
}
