package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type HumanPlayer struct {
	BasePlayer
	cellnet.Session
	Seq    uint32
	Timer  *time.Timer
	logger logrus.FieldLogger
}

func (r *HumanPlayer) String() string {
	if r.roleSkillsData.FaceUp {
		return strconv.Itoa(r.Location()) + "号[" + r.roleSkillsData.Name + "]"
	}
	return strconv.Itoa(r.Location()) + "号[玩家]"
}

func (r *HumanPlayer) Init(game *Game, location int, identity protos.Color, secretTask protos.SecretTask, roleSkillsData *RoleSkillsData) {
	r.logger = logrus.WithField("human_player", r.Location())
	r.BasePlayer.Init(game, location, identity, secretTask, roleSkillsData)
	msg := &protos.InitToc{
		PlayerCount: uint32(len(r.GetGame().GetPlayers())),
		Identity:    identity,
		SecretTask:  secretTask,
	}
	l := location
	for {
		if l < len(RoleCache) {
			msg.Roles = append(msg.Roles, RoleCache[l].Role)
		} else {
			msg.Roles = append(msg.Roles, protos.Role_unknown)
		}
		l = (l + 1) % len(game.GetPlayers())
		if l == location {
			break
		}
	}
	r.Send(msg)
	r.Seq++
}

func (r *HumanPlayer) IncrSeq() {
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
}

func (r *HumanPlayer) NotifyAddHandCard(location int, unknownCount int, cards ...ICard) {
	msg := &protos.AddCardToc{
		PlayerId:         r.GetAlternativeLocation(location),
		UnknownCardCount: uint32(unknownCount),
	}
	for _, card := range cards {
		msg.Cards = append(msg.Cards, card.ToPbCard())
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyDrawPhase(player IPlayer) {
	playerId := r.GetAlternativeLocation(player.Location())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Draw_Phase,
		WaitingPlayerId: playerId,
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyMainPhase(player IPlayer, waitSecond uint32) {
	playerId := r.GetAlternativeLocation(player.Location())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Main_Phase,
		WaitingPlayerId: playerId,
		WaitingSecond:   waitSecond,
	}
	if r.Location() == player.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.IncrSeq()
					r.GetGame().Resolve(&SendPhaseStart{Player: player})
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifySendPhaseStart(player IPlayer, waitSecond uint32) {
	playerId := r.GetAlternativeLocation(player.Location())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Send_Start_Phase,
		WaitingPlayerId: playerId,
		WaitingSecond:   waitSecond,
	}
	if r.Location() == player.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.IncrSeq()
					autoSendMessageCard(r, false)
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifySendMessageCard(player, targetPlayer IPlayer, lockedPlayers []IPlayer, messageCard ICard, dir protos.Direction) {
	msg := &protos.SendMessageCardToc{
		PlayerId:       r.GetAlternativeLocation(player.Location()),
		TargetPlayerId: r.GetAlternativeLocation(targetPlayer.Location()),
		CardDir:        dir,
	}
	if player.Location() == r.Location() {
		msg.CardId = messageCard.GetId()
	}
	for _, p := range lockedPlayers {
		msg.LockPlayerIds = append(msg.LockPlayerIds, r.GetAlternativeLocation(p.Location()))
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifySendPhase(whoseTurn, inFrontOfWhom IPlayer, lockedPlayers []IPlayer, messageCard ICard, dir protos.Direction, isMessageFaceUp bool, waitSecond uint32) {
	playerId := r.GetAlternativeLocation(whoseTurn.Location())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Send_Phase,
		MessagePlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
		MessageCardDir:  dir,
		WaitingPlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
		WaitingSecond:   waitSecond,
	}
	if isMessageFaceUp {
		msg.MessageCard = messageCard.ToPbCard()
	}
	if r.Location() == inFrontOfWhom.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.IncrSeq()
					if func(r IPlayer) bool {
						for _, p := range lockedPlayers {
							if r.Location() == p.Location() {
								return true
							}
						}
						return inFrontOfWhom.Location() == whoseTurn.Location()
					}(r) {
						r.GetGame().Resolve(&OnChooseReceiveCard{
							WhoseTurn:           whoseTurn,
							MessageCard:         messageCard,
							InFrontOfWhom:       inFrontOfWhom,
							IsMessageCardFaceUp: isMessageFaceUp,
						})
					} else {
						r.GetGame().Resolve(&MessageMoveNext{
							SendPhase: &SendPhaseIdle{
								WhoseTurn:           whoseTurn,
								MessageCard:         messageCard,
								Dir:                 dir,
								InFrontOfWhom:       inFrontOfWhom,
								LockedPlayers:       lockedPlayers,
								IsMessageCardFaceUp: isMessageFaceUp,
							},
						})
					}
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyChooseReceiveCard(player IPlayer) {
	r.Send(&protos.ChooseReceiveToc{PlayerId: r.GetAlternativeLocation(player.Location())})
}

func (r *HumanPlayer) NotifyFightPhase(whoseTurn, inFrontOfWhom, whoseFightTurn IPlayer, messageCard ICard, isMessageFaceUp bool, waitSecond uint32) {
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: r.GetAlternativeLocation(whoseTurn.Location()),
		CurrentPhase:    protos.Phase_Fight_Phase,
		MessagePlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
		WaitingPlayerId: r.GetAlternativeLocation(whoseFightTurn.Location()),
		WaitingSecond:   waitSecond,
	}
	if isMessageFaceUp {
		msg.MessageCard = messageCard.ToPbCard()
	}
	if r.Location() == whoseFightTurn.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.IncrSeq()
					r.GetGame().Resolve(&FightPhaseNext{
						FightPhase: &FightPhaseIdle{
							WhoseTurn:           whoseTurn,
							MessageCard:         messageCard,
							InFrontOfWhom:       inFrontOfWhom,
							WhoseFightTurn:      whoseFightTurn,
							IsMessageCardFaceUp: isMessageFaceUp,
						},
					})
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyReceivePhase(whoseTurn, inFrontOfWhom IPlayer, messageCard ICard) {
	r.Send(&protos.NotifyPhaseToc{
		CurrentPlayerId: r.GetAlternativeLocation(whoseTurn.Location()),
		CurrentPhase:    protos.Phase_Receive_Phase,
		MessagePlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
		MessageCard:     messageCard.ToPbCard(),
		WaitingPlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
	})
}

func (r *HumanPlayer) NotifyReceivePhaseWithWaiting(whoseTurn, inFrontOfWhom IPlayer, messageCard ICard, waitingPlayer IPlayer, waitSecond uint32) {
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: r.GetAlternativeLocation(whoseTurn.Location()),
		CurrentPhase:    protos.Phase_Receive_Phase,
		MessagePlayerId: r.GetAlternativeLocation(inFrontOfWhom.Location()),
		MessageCard:     messageCard.ToPbCard(),
		WaitingPlayerId: r.GetAlternativeLocation(waitingPlayer.Location()),
		WaitingSecond:   waitSecond,
	}
	if r.Location() == waitingPlayer.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			if seq == r.Seq {
				r.GetGame().TryContinueResolveProtocol(r, &protos.EndReceivePhaseTos{Seq: seq})
			}
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyDie(location int, loseGame bool) {
	if location == r.Location() {
		var cards []ICard
		for _, card := range r.GetCards() {
			cards = append(cards, card)
		}
		r.GetGame().PlayerDiscardCard(r, cards...)
		r.GetGame().GetDeck().Discard(r.DeleteAllMessageCards()...)
	}
	r.Send(&protos.NotifyDieToc{
		PlayerId: r.GetAlternativeLocation(location),
		LoseGame: loseGame,
	})
}

func (r *HumanPlayer) NotifyWin(declareWinner []IPlayer, winner []IPlayer) {
	msg := &protos.NotifyWinnerToc{
		Identity:    make([]protos.Color, len(r.GetGame().GetPlayers())),
		SecretTasks: make([]protos.SecretTask, len(r.GetGame().GetPlayers())),
	}
	for _, p := range declareWinner {
		msg.DeclarePlayerIds = append(msg.DeclarePlayerIds, r.GetAlternativeLocation(p.Location()))
	}
	for _, p := range winner {
		msg.WinnerIds = append(msg.WinnerIds, r.GetAlternativeLocation(p.Location()))
	}
	for _, p := range r.GetGame().GetPlayers() {
		index := r.GetAlternativeLocation(p.Location())
		msg.Identity[index], msg.SecretTasks[index] = p.GetIdentity()
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyAskForChengQing(whoDie IPlayer, askWhom IPlayer) {
	msg := &protos.WaitForChengQingToc{
		DiePlayerId:     r.GetAlternativeLocation(whoDie.Location()),
		WaitingPlayerId: r.GetAlternativeLocation(askWhom.Location()),
		WaitingSecond:   20,
	}
	if askWhom.Location() == r.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		time.AfterFunc(time.Duration(msg.WaitingSecond+2)*time.Second, func() {
			Post(func() {
				if r.Seq == seq {
					r.IncrSeq()
					r.GetGame().Resolve(&WaitNextForChengQing{
						WaitForChengQing: r.GetGame().GetFsm().(*WaitForChengQing),
					})
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) WaitForDieGiveCard(whoDie IPlayer) {
	msg := &protos.WaitForDieGiveCardToc{
		PlayerId:      r.GetAlternativeLocation(whoDie.Location()),
		WaitingSecond: 30,
		Seq:           r.Seq,
	}
	if whoDie.Location() == r.Location() {
		msg.Seq = r.Seq
		seq := r.Seq
		time.AfterFunc(time.Duration(msg.WaitingSecond+2)*time.Second, func() {
			Post(func() {
				if r.Seq == seq {
					r.IncrSeq()
					fsm := r.GetGame().GetFsm().(*WaitForDieGiveCard)
					r.GetGame().Resolve(&AfterDieGiveCard{DieGiveCard: fsm})
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) onEndMainPhase(pb *protos.EndMainPhaseTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		r.logger.Error("不是你的回合的出牌阶段")
		return
	}
	r.IncrSeq()
	r.GetGame().Resolve(&SendPhaseStart{Player: fsm.Player})
}

func (r *HumanPlayer) onUseShiTan(pb *protos.UseShiTanTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Shi_Tan {
		r.logger.Error("这张牌不是试探，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
		return
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target)
	}
}

func (r *HumanPlayer) onExecuteShiTan(pb *protos.ExecuteShiTanTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	r.GetGame().TryContinueResolveProtocol(r, pb)
}

func (r *HumanPlayer) onUseLiYou(pb *protos.UseLiYouTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Li_You {
		r.logger.Error("这张牌不是利诱，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
		return
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target)
	}
}

func (r *HumanPlayer) onUsePingHeng(pb *protos.UsePingHengTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Ping_Heng {
		r.logger.Error("这张牌不是平衡，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target)
	}
}

func (r *HumanPlayer) onUseWeiBi(pb *protos.UseWeiBiTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Wei_Bi {
		r.logger.Error("这张牌不是威逼，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
		return
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target, pb.WantType) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target, pb.WantType)
	}
}

func (r *HumanPlayer) onWeiBiGiveCard(pb *protos.WeiBiGiveCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	r.GetGame().TryContinueResolveProtocol(r, pb)
}

func (r *HumanPlayer) onUseChengQing(pb *protos.UseChengQingTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Cheng_Qing {
		r.logger.Error("这张牌不是澄清，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
		return
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target, pb.TargetCardId) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target, pb.TargetCardId)
	}
}

func (r *HumanPlayer) onSendMessageCard(pb *protos.SendMessageCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*SendPhaseStart)
	if !ok || r.Location() != fsm.Player.Location() {
		r.logger.Error("不是传递情报的时机")
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if pb.TargetPlayerId == 0 || pb.TargetPlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.TargetPlayerId)
		return
	}
	if pb.CardDir != card.GetDirection() {
		r.logger.Error("方向错误: ", pb.TargetPlayerId)
		return
	}
	if r.FindSkill(SkillIdMingEr) == nil {
		var targetLocation int
		switch pb.CardDir {
		case protos.Direction_Left:
			targetLocation = (r.Location() + len(r.GetGame().GetPlayers()) - 1) % len(r.GetGame().GetPlayers())
			for !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
				targetLocation = (targetLocation + len(r.GetGame().GetPlayers()) - 1) % len(r.GetGame().GetPlayers())
			}
		case protos.Direction_Right:
			targetLocation = (r.Location() + 1) % len(r.GetGame().GetPlayers())
			for !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
				targetLocation++
			}
		}
		if pb.CardDir != protos.Direction_Up && pb.TargetPlayerId != r.GetAlternativeLocation(targetLocation) {
			r.logger.Error("不能传给那个人: ", pb.TargetPlayerId)
			return
		}
	}
	if card.CanLock() {
		if len(pb.LockPlayerId) > 1 {
			r.logger.Error("最多锁定一个目标")
			return
		} else if len(pb.LockPlayerId) == 1 {
			if pb.LockPlayerId[0] >= uint32(len(r.GetGame().GetPlayers())) {
				r.logger.Error("锁定目标错误: ", pb.LockPlayerId[0])
				return
			} else if pb.LockPlayerId[0] == 0 {
				r.logger.Error("不能锁定自己")
				return
			}
		}
	} else {
		if len(pb.LockPlayerId) > 0 {
			r.logger.Error("这张情报没有锁定标记")
			return
		}
	}
	targetLocation := r.GetAbstractLocation(int(pb.TargetPlayerId))
	if !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
		r.logger.Error("目标已死亡")
		return
	}
	var lockPlayers []IPlayer
	for _, lockPlayerId := range pb.LockPlayerId {
		lockPlayer := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(lockPlayerId))]
		if !lockPlayer.IsAlive() {
			r.logger.Error("锁定目标已死亡：", lockPlayer)
			return
		}
		lockPlayers = append(lockPlayers, lockPlayer)
	}
	r.IncrSeq()
	r.GetGame().Resolve(&OnSendCard{
		WhoseTurn:     fsm.Player,
		MessageCard:   card,
		Dir:           pb.CardDir,
		TargetPlayer:  r.GetGame().GetPlayers()[targetLocation],
		LockedPlayers: lockPlayers,
	})
}

func (r *HumanPlayer) onChooseWhetherReceive(pb *protos.ChooseWhetherReceiveTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*SendPhaseIdle)
	if !ok || r.Location() != fsm.InFrontOfWhom.Location() {
		r.logger.Error("不是选择是否接收情报的时机")
		return
	}
	if pb.Receive {
		r.IncrSeq()
		r.GetGame().Resolve(&OnChooseReceiveCard{
			WhoseTurn:           fsm.WhoseTurn,
			MessageCard:         fsm.MessageCard,
			InFrontOfWhom:       fsm.InFrontOfWhom,
			IsMessageCardFaceUp: fsm.IsMessageCardFaceUp,
		})
	} else {
		if r.Location() == fsm.WhoseTurn.Location() {
			r.logger.Error("传出者必须接收")
			return
		}
		if func(e int, lockPlayers []IPlayer) bool {
			for _, a := range lockPlayers {
				if e == a.Location() {
					return true
				}
			}
			return false
		}(r.Location(), fsm.LockedPlayers) {
			r.logger.Error("被锁定，必须接收")
			return
		}
		r.IncrSeq()
		r.GetGame().Resolve(&MessageMoveNext{SendPhase: fsm})
	}
}

func (r *HumanPlayer) onEndFightPhase(pb *protos.EndFightPhaseTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*FightPhaseIdle)
	if !ok || r.Location() != fsm.WhoseFightTurn.Location() {
		r.logger.Error("时机不对")
		return
	}
	r.IncrSeq()
	r.GetGame().Resolve(&FightPhaseNext{FightPhase: fsm})
}

func (r *HumanPlayer) onChengQingSaveDie(pb *protos.ChengQingSaveDieTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*WaitForChengQing)
	if !ok {
		r.logger.Error("现在不是使用澄清的时机")
		return
	}
	if !pb.Use {
		r.IncrSeq()
		r.GetGame().Resolve(&WaitNextForChengQing{WaitForChengQing: fsm})
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Cheng_Qing {
		r.logger.Error("这张牌不是澄清，而是", card)
		return
	}
	target := fsm.WhoDie
	if card.CanUse(r.GetGame(), r, target, pb.TargetCardId) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target, pb.TargetCardId)
	}
}

func (r *HumanPlayer) onDieGiveCard(pb *protos.DieGiveCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	fsm, ok := r.GetGame().GetFsm().(*WaitForDieGiveCard)
	if !ok || r.Location() != fsm.DiedQueue[fsm.DiedIndex].Location() {
		r.logger.Error("你没有死亡")
		return
	}
	if pb.TargetPlayerId == 0 {
		r.IncrSeq()
		r.GetGame().Resolve(&AfterDieGiveCard{DieGiveCard: fsm})
		return
	} else if pb.TargetPlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.TargetPlayerId)
		return
	}
	if len(pb.CardId) == 0 {
		r.logger.Warn("参数似乎有些不对，姑且认为不给牌吧")
		r.IncrSeq()
		r.GetGame().Resolve(&AfterDieGiveCard{DieGiveCard: fsm})
		return
	}
	var cards []ICard
	for _, cardId := range pb.CardId {
		card := r.FindCard(cardId)
		if card == nil {
			r.logger.Error("没有这张牌")
			return
		}
		cards = append(cards, card)
	}
	for _, card := range cards {
		r.DeleteCard(card.GetId())
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.TargetPlayerId))]
	target.AddCards(cards...)
	logger.Info(r, "给了", target, cards)
	for _, p := range r.GetGame().GetPlayers() {
		if player, ok := p.(*HumanPlayer); ok {
			msg := &protos.NotifyDieGiveCardToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
				CardCount:      uint32(len(cards)),
			}
			if p.Location() == r.Location() || p.Location() == target.Location() {
				for _, card := range cards {
					msg.Card = append(msg.Card, card.ToPbCard())
				}
			}
			player.Send(msg)
		}
	}
	r.IncrSeq()
	r.GetGame().Resolve(&AfterDieGiveCard{DieGiveCard: fsm})
}

func (r *HumanPlayer) onUsePoYi(pb *protos.UsePoYiTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Po_Yi {
		r.logger.Error("这张牌不是破译，而是", card)
		return
	}
	if card.CanUse(r.GetGame(), r) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r)
	}
}

func (r *HumanPlayer) onPoYiShow(pb *protos.PoYiShowTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	r.GetGame().TryContinueResolveProtocol(r, pb)
}

func (r *HumanPlayer) onUseJieHuo(pb *protos.UseJieHuoTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Jie_Huo {
		r.logger.Error("这张牌不是截获，而是", card)
		return
	}
	if card.CanUse(r.GetGame(), r) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r)
	}
}

func (r *HumanPlayer) onUseWuDao(pb *protos.UseWuDaoTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Wu_Dao {
		r.logger.Error("这张牌不是误导，而是", card)
		return
	}
	if pb.TargetPlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.TargetPlayerId)
		return
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.TargetPlayerId))]
	if card.CanUse(r.GetGame(), r, target) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r, target)
	}
}

func (r *HumanPlayer) onUseDiaoBao(pb *protos.UseDiaoBaoTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Diao_Bao {
		r.logger.Error("这张牌不是调包，而是", card)
		return
	}
	if card.CanUse(r.GetGame(), r) {
		r.IncrSeq()
		card.Execute(r.GetGame(), r)
	}
}

func (r *HumanPlayer) onEndReceivePhase(pb *protos.EndReceivePhaseTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	r.GetGame().TryContinueResolveProtocol(r, pb)
}
