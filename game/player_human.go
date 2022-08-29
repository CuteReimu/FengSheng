package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type HumanPlayer struct {
	interfaces.BasePlayer
	cellnet.Session
	Seq    uint32
	Timer  *time.Timer
	logger logrus.FieldLogger
}

func (r *HumanPlayer) String() string {
	return strconv.Itoa(r.Location()) + "号[玩家]"
}

func (r *HumanPlayer) Init(game interfaces.IGame, location int, identity protos.Color, secretTask protos.SecretTask) {
	r.logger = logrus.WithField("human_player", r.Location())
	r.BasePlayer.Init(game, location, identity, secretTask)
	msg := &protos.InitToc{
		PlayerCount: uint32(len(r.GetGame().GetPlayers())),
		Identity:    identity,
		SecretTask:  secretTask,
	}
	r.Send(msg)
	r.Seq++
}

func (r *HumanPlayer) NotifyAddHandCard(location int, unknownCount int, cards ...interfaces.ICard) {
	msg := &protos.AddCardToc{
		PlayerId:         r.GetAlternativeLocation(location),
		UnknownCardCount: uint32(unknownCount),
	}
	for _, card := range cards {
		msg.Cards = append(msg.Cards, card.ToPbCard())
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyDrawPhase() {
	playerId := r.GetAlternativeLocation(r.GetGame().GetWhoseTurn())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Draw_Phase,
		WaitingPlayerId: playerId,
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyMainPhase(waitSecond uint32) {
	playerId := r.GetAlternativeLocation(r.GetGame().GetWhoseTurn())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Main_Phase,
		WaitingPlayerId: playerId,
		WaitingSecond:   waitSecond,
	}
	if r.Location() == r.GetGame().GetWhoseTurn() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					Post(r.GetGame().SendPhaseStart)
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifySendPhaseStart(waitSecond uint32) {
	playerId := r.GetAlternativeLocation(r.GetGame().GetWhoseTurn())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Send_Start_Phase,
		WaitingPlayerId: playerId,
		WaitingSecond:   waitSecond,
	}
	if r.Location() == r.GetGame().GetWhoseTurn() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					autoSendMessageCard(r, false)
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifySendPhase(waitSecond uint32, isFirstTime bool) {
	playerId := r.GetAlternativeLocation(r.GetGame().GetWhoseTurn())
	if isFirstTime {
		msg := &protos.SendMessageCardToc{
			PlayerId:       playerId,
			TargetPlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn()),
			CardDir:        r.GetGame().GetMessageCardDirection(),
		}
		if r.GetGame().GetWhoseTurn() == r.Location() {
			msg.CardId = r.GetGame().GetCurrentMessageCard().GetId()
		}
		for _, id := range r.GetGame().GetLockPlayers() {
			msg.LockPlayerIds = append(msg.LockPlayerIds, r.GetAlternativeLocation(id))
		}
		r.Send(msg)
	}
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Send_Phase,
		MessagePlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn()),
		MessageCardDir:  r.GetGame().GetMessageCardDirection(),
		WaitingPlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn()),
		WaitingSecond:   waitSecond,
	}
	if r.GetGame().IsMessageCardFaceUp() {
		msg.MessageCard = r.GetGame().GetCurrentMessageCard().ToPbCard()
	}
	if r.Location() == r.GetGame().GetWhoseSendTurn() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					if func(r interfaces.IPlayer) bool {
						for _, p := range r.GetGame().GetLockPlayers() {
							if r.Location() == p {
								return true
							}
						}
						return r.GetGame().GetWhoseSendTurn() == r.GetGame().GetWhoseTurn()
					}(r) {
						Post(r.GetGame().OnChooseReceiveCard)
					} else {
						Post(r.GetGame().MessageMoveNext)
					}
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyChooseReceiveCard() {
	r.Send(&protos.ChooseReceiveToc{PlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn())})
}

func (r *HumanPlayer) NotifyFightPhase(waitSecond uint32) {
	playerId := r.GetAlternativeLocation(r.GetGame().GetWhoseTurn())
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Fight_Phase,
		MessagePlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn()),
		MessageCardDir:  r.GetGame().GetMessageCardDirection(),
		WaitingPlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseFightTurn()),
		WaitingSecond:   waitSecond,
	}
	if r.GetGame().IsMessageCardFaceUp() {
		msg.MessageCard = r.GetGame().GetCurrentMessageCard().ToPbCard()
	}
	if r.Location() == r.GetGame().GetWhoseFightTurn() {
		msg.Seq = r.Seq
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond+2), func() {
			Post(func() {
				if seq == r.Seq {
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					Post(r.GetGame().FightPhaseNext)
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyReceivePhase() {
	r.Send(&protos.NotifyPhaseToc{
		CurrentPlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseTurn()),
		CurrentPhase:    protos.Phase_Receive_Phase,
		MessagePlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseSendTurn()),
		MessageCardDir:  r.GetGame().GetMessageCardDirection(),
		MessageCard:     r.GetGame().GetCurrentMessageCard().ToPbCard(),
		WaitingPlayerId: r.GetAlternativeLocation(r.GetGame().GetWhoseFightTurn()),
	})
}

func (r *HumanPlayer) NotifyDie(location int, loseGame bool) {
	if location == r.Location() {
		r.SetAlive(false)
		var cards []interfaces.ICard
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

func (r *HumanPlayer) NotifyWin(declareWinner interfaces.IPlayer, winner []interfaces.IPlayer) {
	msg := &protos.NotifyWinnerToc{
		DeclarePlayerId: r.GetAlternativeLocation(declareWinner.Location()),
		Identity:        make([]protos.Color, len(r.GetGame().GetPlayers())),
		SecretTasks:     make([]protos.SecretTask, len(r.GetGame().GetPlayers())),
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

func (r *HumanPlayer) NotifyAskForChengQing(whoDie interfaces.IPlayer, askWhom interfaces.IPlayer) {
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
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					Post(r.GetGame().AskNextForChengQing)
				}
			})
		})
	}
	r.Send(msg)
}

func (r *HumanPlayer) WaitForDieGiveCard(whoDie interfaces.IPlayer) {
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
					r.Seq++
					if r.Timer != nil {
						r.Timer.Stop()
					}
					for _, p := range r.GetGame().GetPlayers() {
						p.NotifyDie(whoDie.Location(), false)
					}
					Post(r.GetGame().AfterChengQing)
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
	if r.Location() != r.GetGame().GetWhoseTurn() || r.GetGame().GetCurrentPhase() != protos.Phase_Main_Phase || !r.GetGame().IsIdleTimePoint() {
		r.logger.Error("不是你的回合的出牌阶段")
		return
	}
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
	Post(r.GetGame().SendPhaseStart)
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r, target)
	}
}

func (r *HumanPlayer) onExecuteShiTan(pb *protos.ExecuteShiTanTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	currentCard := r.GetGame().GetCurrentCard()
	if currentCard == nil || currentCard.Card.GetType() != protos.CardType_Shi_Tan {
		r.logger.Error("现在并不在结算试探", currentCard.Card)
		return
	}
	if currentCard.TargetPlayer != r.Location() {
		r.logger.Error("你不是试探的目标", currentCard.Card)
		return
	}
	if currentCard.Card.CanUse2(r.GetGame(), r, pb.CardId) {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		currentCard.Card.Execute2(r.GetGame(), r, pb.CardId)
	}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r, target, pb.WantType)
	}
}

func (r *HumanPlayer) onWeiBiGiveCard(pb *protos.WeiBiGiveCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	currentCard := r.GetGame().GetCurrentCard()
	if currentCard == nil || currentCard.Card.GetType() != protos.CardType_Wei_Bi {
		r.logger.Error("现在并不在结算威逼")
		return
	}
	if currentCard.TargetPlayer != r.Location() {
		r.logger.Error("你不是威逼的目标", currentCard.Card)
		return
	}
	if currentCard.Card.CanUse2(r.GetGame(), r.GetGame().GetPlayers()[currentCard.Player], r, pb.CardId) {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		currentCard.Card.Execute2(r.GetGame(), r.GetGame().GetPlayers()[currentCard.Player], r, pb.CardId)
	}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r, target, pb.TargetCardId)
	}
}

func (r *HumanPlayer) onSendMessageCard(pb *protos.SendMessageCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	if r.GetGame().GetWhoseTurn() != r.Location() || r.GetGame().GetCurrentPhase() != protos.Phase_Send_Start_Phase || !r.GetGame().IsIdleTimePoint() {
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
	targetLocation = r.GetAbstractLocation(int(pb.TargetPlayerId))
	if !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
		r.logger.Error("目标已死亡")
		return
	}
	var lockLocation []int
	for _, lockPlayerId := range pb.LockPlayerId {
		lockLocation = append(lockLocation, r.GetAbstractLocation(int(lockPlayerId)))
	}
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
	Post(func() { r.GetGame().OnSendCard(card, pb.CardDir, targetLocation, lockLocation) })
}

func (r *HumanPlayer) onChooseWhetherReceive(pb *protos.ChooseWhetherReceiveTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	if r.GetGame().GetWhoseSendTurn() != r.Location() || r.GetGame().GetCurrentPhase() != protos.Phase_Send_Phase || !r.GetGame().IsIdleTimePoint() {
		r.logger.Error("不是选择是否接收情报的时机")
		return
	}
	if pb.Receive {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		Post(r.GetGame().OnChooseReceiveCard)
	} else {
		if r.Location() == r.GetGame().GetWhoseTurn() {
			r.logger.Error("传出者必须接收")
			return
		}
		if func(e int, arr []int) bool {
			for _, a := range arr {
				if e == a {
					return true
				}
			}
			return false
		}(r.Location(), r.GetGame().GetLockPlayers()) {
			r.logger.Error("被锁定，必须接收")
			return
		}
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		Post(r.GetGame().MessageMoveNext)
	}
}

func (r *HumanPlayer) onEndFightPhase(pb *protos.EndFightPhaseTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	if r.GetGame().GetCurrentPhase() != protos.Phase_Fight_Phase {
		r.logger.Error("时机不对")
		return
	}
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
	Post(r.GetGame().FightPhaseNext)
}

func (r *HumanPlayer) onChengQingSaveDie(pb *protos.ChengQingSaveDieTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	if r.GetGame().GetDieState() != interfaces.DieStateWaitForChengQing {
		r.logger.Error("现在不是使用澄清的时候")
		return
	}
	if !pb.Use {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		Post(r.GetGame().AskNextForChengQing)
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
	target := r.GetGame().GetPlayers()[r.GetGame().GetWhoDie()]
	if card.CanUse(r.GetGame(), r, target, pb.TargetCardId) {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r, target, pb.TargetCardId)
	}
}

func (r *HumanPlayer) onDieGiveCard(pb *protos.DieGiveCardTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	if r.GetGame().GetDieState() != interfaces.DieStateDying {
		r.logger.Error("你没有死亡")
		return
	}
	if pb.TargetPlayerId == 0 {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		for _, p := range r.GetGame().GetPlayers() {
			p.NotifyDie(r.GetGame().GetWhoDie(), false)
		}
		Post(r.GetGame().AfterChengQing)
		return
	} else if pb.TargetPlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.TargetPlayerId)
		return
	}
	if len(pb.CardId) == 0 {
		r.logger.Warn("参数似乎有些不对，姑且认为不给牌吧")
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		for _, p := range r.GetGame().GetPlayers() {
			p.NotifyDie(r.GetGame().GetWhoDie(), false)
		}
		Post(r.GetGame().AfterChengQing)
		return
	}
	var cards []interfaces.ICard
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
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
	for _, p := range r.GetGame().GetPlayers() {
		p.NotifyDie(r.GetGame().GetWhoDie(), false)
	}
	Post(r.GetGame().AfterChengQing)
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r)
	}
}

func (r *HumanPlayer) onPoYiShow(pb *protos.PoYiShowTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	currentCard := r.GetGame().GetCurrentCard()
	if currentCard == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if currentCard.Player != r.Location() {
		r.logger.Error("你不是破译的使用者")
		return
	}
	if currentCard.Card.GetType() != protos.CardType_Po_Yi {
		r.logger.Error("这张牌不是误导，而是", currentCard)
		return
	}
	if currentCard.Card.CanUse2(r.GetGame(), r, pb.Show) {
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		currentCard.Card.Execute2(r.GetGame(), r, pb.Show)
	}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
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
		r.Seq++
		if r.Timer != nil {
			r.Timer.Stop()
		}
		card.Execute(r.GetGame(), r)
	}
}
