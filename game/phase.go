package game

import (
	"github.com/CuteReimu/FengSheng/protos"
)

// DrawPhase 摸牌阶段
type DrawPhase struct {
	Player IPlayer
}

func (dp *DrawPhase) Resolve() (next Fsm, continueResolve bool) {
	game := dp.Player.GetGame()
	if !dp.Player.IsAlive() {
		return &NextTurn{player: dp.Player}, true
	}
	logger.Info(dp.Player, "的回合开始了")
	for _, p := range game.GetPlayers() {
		p.NotifyDrawPhase(dp.Player)
	}
	dp.Player.Draw(3)
	return &MainPhaseIdle{Player: dp.Player}, true
}

// MainPhaseIdle 出牌阶段空闲时点
type MainPhaseIdle struct {
	Player IPlayer
}

func (mp *MainPhaseIdle) Resolve() (next Fsm, continueResolve bool) {
	game := mp.Player.GetGame()
	if !mp.Player.IsAlive() {
		return &NextTurn{player: mp.Player}, true
	}
	for _, p := range game.GetPlayers() {
		p.NotifyMainPhase(mp.Player, 30)
	}
	return mp, false
}

// SendPhaseStart 情报传递阶段开始时，选择传递一张情报
type SendPhaseStart struct {
	Player IPlayer
}

func (sp *SendPhaseStart) Resolve() (next Fsm, continueResolve bool) {
	game := sp.Player.GetGame()
	if sp.Player.IsAlive() {
		if len(sp.Player.GetCards()) == 0 {
			logger.Info(sp.Player, "没有情报可传，输掉了游戏")
			game.GetDeck().Discard(sp.Player.DeleteAllMessageCards()...)
			sp.Player.SetLose(true)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(sp.Player.Location(), true)
			}
		}
	}
	if !sp.Player.IsAlive() {
		return &NextTurn{player: sp.Player}, true
	}
	for _, p := range game.GetPlayers() {
		p.NotifySendPhaseStart(sp.Player, 20)
	}
	return sp, false
}

// OnSendCard 选择了将要传递一张情报时
type OnSendCard struct {
	WhoseTurn     IPlayer
	MessageCard   ICard
	Dir           protos.Direction
	TargetPlayer  IPlayer
	LockedPlayers []IPlayer
}

func (sc *OnSendCard) Resolve() (next Fsm, continueResolve bool) {
	game := sc.WhoseTurn.GetGame()
	logger.Info(sc.WhoseTurn, "传出了", sc.MessageCard, "，方向是", protos.Direction_name[int32(sc.Dir)], "传给了", sc.TargetPlayer)
	sc.WhoseTurn.DeleteCard(sc.MessageCard.GetId())
	for _, p := range game.GetPlayers() {
		p.NotifySendMessageCard(sc.WhoseTurn, sc.TargetPlayer, sc.LockedPlayers, sc.MessageCard, sc.Dir)
	}
	logger.Info("情报到达", sc.TargetPlayer, "面前")
	return &SendPhaseIdle{
		WhoseTurn:     sc.WhoseTurn,
		MessageCard:   sc.MessageCard,
		Dir:           sc.Dir,
		InFrontOfWhom: sc.TargetPlayer,
		LockedPlayers: sc.LockedPlayers,
	}, true
}

// SendPhaseIdle 情报传递阶段空闲时点
type SendPhaseIdle struct {
	WhoseTurn           IPlayer
	MessageCard         ICard
	Dir                 protos.Direction
	InFrontOfWhom       IPlayer
	LockedPlayers       []IPlayer
	IsMessageCardFaceUp bool
}

func (sp *SendPhaseIdle) Resolve() (next Fsm, continueResolve bool) {
	game := sp.WhoseTurn.GetGame()
	for _, p := range game.GetPlayers() {
		p.NotifySendPhase(sp.WhoseTurn, sp.InFrontOfWhom, sp.LockedPlayers, sp.MessageCard, sp.Dir, sp.IsMessageCardFaceUp, 20)
	}
	return sp, false
}

// MessageMoveNext 情报传递阶段，情报移到下一个人
type MessageMoveNext struct {
	SendPhase *SendPhaseIdle
}

func (mm *MessageMoveNext) Resolve() (next Fsm, continueResolve bool) {
	game := mm.SendPhase.WhoseTurn.GetGame()
	if mm.SendPhase.Dir == protos.Direction_Up {
		if mm.SendPhase.WhoseTurn.IsAlive() {
			mm.SendPhase.InFrontOfWhom = mm.SendPhase.WhoseTurn
			logger.Info("情报到达", mm.SendPhase.InFrontOfWhom, "面前")
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
				logger.Info("情报到达", mm.SendPhase.InFrontOfWhom, "面前")
				return mm.SendPhase, true
			} else if mm.SendPhase.WhoseTurn == mm.SendPhase.InFrontOfWhom {
				return &NextTurn{player: mm.SendPhase.WhoseTurn}, true
			}
		}
	}
}

// OnChooseReceiveCard 选择接收情报时
type OnChooseReceiveCard struct {
	WhoseTurn           IPlayer
	MessageCard         ICard
	InFrontOfWhom       IPlayer
	IsMessageCardFaceUp bool
}

func (cr *OnChooseReceiveCard) Resolve() (next Fsm, continueResolve bool) {
	game := cr.WhoseTurn.GetGame()
	logger.Info(cr.InFrontOfWhom, "选择接收情报")
	for _, p := range game.GetPlayers() {
		p.NotifyChooseReceiveCard(cr.InFrontOfWhom)
	}
	return &FightPhaseIdle{
		WhoseTurn:           cr.WhoseTurn,
		MessageCard:         cr.MessageCard,
		InFrontOfWhom:       cr.InFrontOfWhom,
		WhoseFightTurn:      cr.InFrontOfWhom,
		IsMessageCardFaceUp: cr.IsMessageCardFaceUp,
	}, true
}

// FightPhaseIdle 争夺阶段空闲时点
type FightPhaseIdle struct {
	WhoseTurn           IPlayer
	MessageCard         ICard
	InFrontOfWhom       IPlayer
	WhoseFightTurn      IPlayer
	IsMessageCardFaceUp bool
}

func (fp *FightPhaseIdle) Resolve() (next Fsm, continueResolve bool) {
	game := fp.WhoseTurn.GetGame()
	for _, p := range game.GetPlayers() {
		p.NotifyFightPhase(fp.WhoseTurn, fp.InFrontOfWhom, fp.WhoseFightTurn, fp.MessageCard, fp.IsMessageCardFaceUp, 20)
	}
	return fp, false
}

// FightPhaseNext 争夺阶段即将询问下一个人时
type FightPhaseNext struct {
	FightPhase *FightPhaseIdle
}

func (fp *FightPhaseNext) Resolve() (next Fsm, continueResolve bool) {
	game := fp.FightPhase.WhoseFightTurn.GetGame()
	whoseFightTurn := fp.FightPhase.WhoseFightTurn.Location()
	for {
		whoseFightTurn = (whoseFightTurn + 1) % len(game.GetPlayers())
		if whoseFightTurn == fp.FightPhase.InFrontOfWhom.Location() {
			return &ReceivePhase{
				WhoseTurn:     fp.FightPhase.WhoseTurn,
				MessageCard:   fp.FightPhase.MessageCard,
				InFrontOfWhom: fp.FightPhase.InFrontOfWhom,
			}, true
		} else if game.GetPlayers()[whoseFightTurn].IsAlive() {
			break
		}
	}
	nextPhase := fp.FightPhase
	nextPhase.WhoseFightTurn = game.GetPlayers()[whoseFightTurn]
	return nextPhase, true
}

// ReceivePhase 情报接收阶段
type ReceivePhase struct {
	WhoseTurn     IPlayer
	MessageCard   ICard
	InFrontOfWhom IPlayer
}

func (rp *ReceivePhase) Resolve() (next Fsm, continueResolve bool) {
	player := rp.InFrontOfWhom
	if player.IsAlive() {
		player.AddMessageCards(rp.MessageCard)
		logger.Info(player, "成功接收情报")
		for _, p := range player.GetGame().GetPlayers() {
			p.NotifyReceivePhase(rp.WhoseTurn, rp.InFrontOfWhom, rp.MessageCard)
		}
		return &ReceivePhaseSenderSkill{
			WhoseTurn:     rp.WhoseTurn,
			InFrontOfWhom: rp.InFrontOfWhom,
		}, true
	}
	return &NextTurn{player: rp.WhoseTurn}, true
}

// ReceivePhaseSenderSkill 接收情报时，传出者的技能
type ReceivePhaseSenderSkill struct {
	WhoseTurn     IPlayer
	InFrontOfWhom IPlayer
}

func (rp *ReceivePhaseSenderSkill) Resolve() (next Fsm, continueResolve bool) {
	return &ReceivePhaseReceiverSkill{
		WhoseTurn:     rp.WhoseTurn,
		InFrontOfWhom: rp.InFrontOfWhom,
	}, true
}

// ReceivePhaseReceiverSkill 接收情报时，接收者的技能
type ReceivePhaseReceiverSkill struct {
	WhoseTurn     IPlayer
	InFrontOfWhom IPlayer
}

func (rp *ReceivePhaseReceiverSkill) Resolve() (next Fsm, continueResolve bool) {
	return &CheckWinOrDie{
		WhoseTurn:       rp.WhoseTurn,
		ReceiveOrder:    []IPlayer{rp.InFrontOfWhom},
		AfterDieResolve: &NextTurn{player: rp.WhoseTurn},
	}, true
}

// CheckWinOrDie 判断是否胜利或者有人濒死
type CheckWinOrDie struct {
	WhoseTurn       IPlayer
	ReceiveOrder    []IPlayer
	AfterDieResolve Fsm
}

func (cw *CheckWinOrDie) Resolve() (next Fsm, continueResolve bool) {
	game := cw.WhoseTurn.GetGame()
	var stealer IPlayer
	var redPlayers, bluePlayers []IPlayer
	for _, p := range game.GetPlayers() {
		if p.IsAlive() && !p.HasNoIdentity() && !p.IsLose() {
			identity, secretTask := p.GetIdentity()
			switch identity {
			case protos.Color_Black:
				if secretTask == protos.SecretTask_Stealer {
					stealer = p
				}
			case protos.Color_Red:
				redPlayers = append(redPlayers, p)
			case protos.Color_Blue:
				bluePlayers = append(bluePlayers, p)
			}
		}
	}
	var declareWinner []IPlayer
	var winner []IPlayer
	var dyingQueue []IPlayer
	var redWin, blueWin bool
	for i := range cw.ReceiveOrder {
		player := cw.ReceiveOrder[i]
		var red, blue, black int
		for _, card := range player.GetMessageCards() {
			for _, color := range card.GetColors() {
				switch color {
				case protos.Color_Black:
					black++
				case protos.Color_Red:
					red++
				case protos.Color_Blue:
					blue++
				}
			}
		}
		identity, secretTask := player.GetIdentity()
		switch identity {
		case protos.Color_Black:
			if secretTask == protos.SecretTask_Collector && (red >= 3 || blue >= 3) {
				declareWinner = append(declareWinner, player)
				winner = append(winner, player)
			}
		case protos.Color_Red:
			if red >= 3 {
				declareWinner = append(declareWinner, player)
				redWin = true
			}
		case protos.Color_Blue:
			if blue >= 3 {
				declareWinner = append(declareWinner, player)
				blueWin = true
			}
		}
		if black >= 3 {
			dyingQueue = append(dyingQueue, player)
		}
	}
	if redWin {
		winner = append(winner, redPlayers...)
	}
	if blueWin {
		winner = append(winner, bluePlayers...)
	}
	if declareWinner != nil && stealer != nil && cw.WhoseTurn.Location() == stealer.Location() {
		declareWinner = []IPlayer{stealer}
		winner = []IPlayer{stealer}
	}
	if declareWinner != nil {
		logger.Info(declareWinner, "宣告胜利，胜利者有", winner)
		for _, p := range game.GetPlayers() {
			p.NotifyWin(declareWinner, winner)
		}
		return nil, false
	}
	return &StartWaitForChengQing{
		WhoseTurn:       cw.WhoseTurn,
		DyingQueue:      dyingQueue,
		AfterDieResolve: cw.AfterDieResolve,
	}, true
}

// StartWaitForChengQing 判断是否需要濒死求澄清
type StartWaitForChengQing struct {
	WhoseTurn       IPlayer
	DyingQueue      []IPlayer
	DiedQueue       []IPlayer
	AfterDieResolve Fsm
}

func (cq *StartWaitForChengQing) Resolve() (next Fsm, continueResolve bool) {
	if len(cq.DyingQueue) == 0 {
		return &CheckKillerWin{
			WhoseTurn:       cq.WhoseTurn,
			DiedQueue:       cq.DiedQueue,
			AfterDieResolve: cq.AfterDieResolve,
		}, true
	}
	whoDie := cq.DyingQueue[0]
	logger.Info(whoDie, "濒死")
	cq.DyingQueue = cq.DyingQueue[1:]
	fsm := &WaitForChengQing{
		WhoseTurn:       cq.WhoseTurn,
		WhoDie:          whoDie,
		AskWhom:         cq.WhoseTurn,
		DyingQueue:      cq.DyingQueue,
		DiedQueue:       cq.DiedQueue,
		AfterDieResolve: cq.AfterDieResolve,
	}
	if whoDie.IsAlive() {
		return fsm, true
	}
	return &WaitNextForChengQing{WaitForChengQing: fsm}, true
}

// WaitForChengQing 濒死求澄清
type WaitForChengQing struct {
	WhoseTurn       IPlayer
	WhoDie          IPlayer
	AskWhom         IPlayer
	DyingQueue      []IPlayer
	DiedQueue       []IPlayer
	AfterDieResolve Fsm
}

func (cq *WaitForChengQing) Resolve() (next Fsm, continueResolve bool) {
	logger.Info("正在询问", cq.AskWhom, "是否使用澄清")
	for _, p := range cq.AskWhom.GetGame().GetPlayers() {
		p.NotifyAskForChengQing(cq.WhoDie, cq.AskWhom)
	}
	return cq, false
}

// UseChengQingOnDying 濒死求澄清时，使用了澄清
type UseChengQingOnDying struct {
	WaitForChengQing *WaitForChengQing
}

func (cq *UseChengQingOnDying) Resolve() (next Fsm, continueResolve bool) {
	var count int
	for _, card := range cq.WaitForChengQing.WhoDie.GetCards() {
		for _, color := range card.GetColors() {
			if color == protos.Color_Black {
				count++
				break
			}
		}
	}
	if count >= 3 {
		return cq.WaitForChengQing, true
	}
	return &StartWaitForChengQing{
		WhoseTurn:       cq.WaitForChengQing.WhoseTurn,
		DyingQueue:      cq.WaitForChengQing.DyingQueue,
		DiedQueue:       cq.WaitForChengQing.DiedQueue,
		AfterDieResolve: cq.WaitForChengQing.AfterDieResolve,
	}, true
}

// WaitNextForChengQing 濒死求澄清时，询问下一个人
type WaitNextForChengQing struct {
	WaitForChengQing *WaitForChengQing
}

func (cq *WaitNextForChengQing) Resolve() (next Fsm, continueResolve bool) {
	game := cq.WaitForChengQing.AskWhom.GetGame()
	askWhom := cq.WaitForChengQing.AskWhom.Location()
	for {
		askWhom = (askWhom + 1) % len(game.GetPlayers())
		if askWhom == cq.WaitForChengQing.WhoseTurn.Location() {
			logger.Info("无人拯救，", cq.WaitForChengQing.WhoDie, "已死亡")
			cq.WaitForChengQing.WhoDie.SetAlive(false)
			cq.WaitForChengQing.DiedQueue = append(cq.WaitForChengQing.DiedQueue, cq.WaitForChengQing.WhoDie)
			return &StartWaitForChengQing{
				WhoseTurn:       cq.WaitForChengQing.WhoseTurn,
				DyingQueue:      cq.WaitForChengQing.DyingQueue,
				DiedQueue:       cq.WaitForChengQing.DiedQueue,
				AfterDieResolve: cq.WaitForChengQing.AfterDieResolve,
			}, true
		}
		if game.GetPlayers()[askWhom].IsAlive() {
			cq.WaitForChengQing.AskWhom = game.GetPlayers()[askWhom]
			return cq.WaitForChengQing, true
		}
	}
}

// CheckKillerWin 判断镇压者获胜条件，或者只剩一个人存活
type CheckKillerWin struct {
	WhoseTurn       IPlayer
	DiedQueue       []IPlayer
	AfterDieResolve Fsm
}

func (ck *CheckKillerWin) Resolve() (next Fsm, continueResolve bool) {
	if len(ck.DiedQueue) == 0 {
		return ck.AfterDieResolve, true
	}
	game := ck.WhoseTurn.GetGame()
	var killer IPlayer
	for _, p := range game.GetPlayers() {
		if p.IsAlive() && !p.HasNoIdentity() && !p.IsLose() {
			if identity, task := p.GetIdentity(); identity == protos.Color_Black && task == protos.SecretTask_Killer {
				killer = p
				break
			}
		}
	}
	if killer != nil && ck.WhoseTurn.Location() == killer.Location() {
		for _, whoDie := range ck.DiedQueue {
			var count int
			for _, card := range whoDie.GetMessageCards() {
				for _, color := range card.GetColors() {
					if color != protos.Color_Black {
						count++
						break
					}
				}
			}
			if count >= 2 {
				logger.Info(killer, "宣告胜利，胜利者有", []IPlayer{killer})
				for _, p := range game.GetPlayers() {
					p.NotifyWin([]IPlayer{killer}, []IPlayer{killer})
				}
				return nil, false
			}
		}
	}
	var alivePlayer IPlayer
	for _, p := range game.GetPlayers() {
		if p.IsAlive() {
			if alivePlayer == nil {
				alivePlayer = p
			} else {
				// 至少有2个人存活，游戏继续
				return &DieSkill{
					WhoseTurn:       ck.WhoseTurn,
					DiedQueue:       ck.DiedQueue,
					AskWhom:         ck.DiedQueue[0],
					AfterDieResolve: ck.AfterDieResolve,
				}, true
			}
		}
	}
	// 只剩1个人存活，游戏结束
	winner := []IPlayer{alivePlayer}
	if identity1, _ := alivePlayer.GetIdentity(); identity1 != protos.Color_Black {
		for _, p := range game.GetPlayers() {
			if identity2, _ := p.GetIdentity(); identity2 == identity1 && p.Location() != alivePlayer.Location() {
				winner = append(winner, p)
			}
		}
	}
	logger.Info("只剩下", alivePlayer, "存活，胜利者有", winner)
	for _, p := range game.GetPlayers() {
		p.NotifyWin(([]IPlayer)(nil), winner)
	}
	return nil, false
}

// DieSkill 死亡角色发动技能
type DieSkill struct {
	WhoseTurn       IPlayer
	DiedIndex       int
	DiedQueue       []IPlayer
	AskWhom         IPlayer
	ReceiveOrder    []IPlayer
	AfterDieResolve Fsm
}

func (ds *DieSkill) Resolve() (next Fsm, continueResolve bool) {
	// TODO 有技能应该暂停
	return &DieSkillNext{DieSkill: ds}, true
}

type DieSkillNext struct {
	DieSkill *DieSkill
}

func (ds *DieSkillNext) Resolve() (next Fsm, continueResolve bool) {
	game := ds.DieSkill.WhoseTurn.GetGame()
	askWhom := ds.DieSkill.AskWhom.Location()
	for {
		askWhom = (askWhom + 1) % len(game.GetPlayers())
		if askWhom == ds.DieSkill.DiedQueue[ds.DieSkill.DiedIndex].Location() {
			ds.DieSkill.DiedIndex++
			if ds.DieSkill.DiedIndex > len(ds.DieSkill.DiedQueue) {
				return &WaitForDieGiveCard{
					WhoseTurn:       ds.DieSkill.WhoseTurn,
					DiedQueue:       ds.DieSkill.DiedQueue,
					ReceiveOrder:    ds.DieSkill.ReceiveOrder,
					AfterDieResolve: ds.DieSkill.AfterDieResolve,
				}, true
			}
			ds.DieSkill.AskWhom = ds.DieSkill.DiedQueue[ds.DieSkill.DiedIndex]
			return ds.DieSkill, true
		}
		if game.GetPlayers()[askWhom].IsAlive() {
			ds.DieSkill.AskWhom = game.GetPlayers()[askWhom]
			return ds.DieSkill, true
		}
	}
}

// WaitForDieGiveCard 等待死亡角色给三张牌
type WaitForDieGiveCard struct {
	WhoseTurn       IPlayer
	DiedIndex       int
	DiedQueue       []IPlayer
	ReceiveOrder    []IPlayer
	AfterDieResolve Fsm
}

func (dg *WaitForDieGiveCard) Resolve() (next Fsm, continueResolve bool) {
	if dg.DiedIndex >= len(dg.DiedQueue) {
		return dg.AfterDieResolve, true
	}
	whoDie := dg.DiedQueue[dg.DiedIndex]
	for _, p := range whoDie.GetGame().GetPlayers() {
		p.WaitForDieGiveCard(whoDie)
	}
	return dg, false
}

type AfterDieGiveCard struct {
	DieGiveCard *WaitForDieGiveCard
}

func (dg *AfterDieGiveCard) Resolve() (next Fsm, continueResolve bool) {
	if len(dg.DieGiveCard.ReceiveOrder) > 0 {
		return &CheckWinOrDie{
			WhoseTurn:       dg.DieGiveCard.WhoseTurn,
			ReceiveOrder:    dg.DieGiveCard.ReceiveOrder,
			AfterDieResolve: dg.DieGiveCard.AfterDieResolve,
		}, true
	} else {
		return dg.DieGiveCard.AfterDieResolve, true
	}
}

// NextTurn 即将跳转到下一回合时
type NextTurn struct {
	player IPlayer
}

func (nt *NextTurn) Resolve() (next Fsm, continueResolve bool) {
	game := nt.player.GetGame()
	whoseTurn := nt.player.Location()
	for {
		whoseTurn = (whoseTurn + 1) % len(game.GetPlayers())
		player := game.GetPlayers()[whoseTurn]
		if player.IsAlive() {
			return &DrawPhase{Player: player}, true
		}
	}
}
