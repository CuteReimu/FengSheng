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
	}
	r.Send(msg)
	if r.Location() == r.GetGame().GetWhoseTurn() {
		seq := r.Seq
		r.Timer = time.AfterFunc(time.Second*time.Duration(waitSecond), func() {
			r.GetGame().Post(func() {
				if seq == r.Seq {
					r.GetGame().Post(r.GetGame().SendPhase)
					r.Seq++
				}
			})
		})
	}
}

func (r *HumanPlayer) onEndMainPhase(pb *protos.EndMainPhaseTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	r.Seq++
	if r.Timer != nil {
		r.Timer.Stop()
	}
	r.GetGame().Post(r.GetGame().SendPhase)
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
