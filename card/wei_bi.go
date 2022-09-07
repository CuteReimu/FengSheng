package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"google.golang.org/protobuf/proto"
	"time"
)

type WeiBi struct {
	game.BaseCard
}

func (card *WeiBi) String() string {
	return utils.CardColorToString(card.Color...) + "威逼"
}

func (card *WeiBi) GetType() protos.CardType {
	return protos.CardType_Wei_Bi
}

func (card *WeiBi) CanUse(g *game.Game, r game.IPlayer, args ...interface{}) bool {
	target := args[0].(game.IPlayer)
	wantType := args[1].(protos.CardType)
	fsm, ok := g.GetFsm().(*game.MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		logger.Error("威逼的使用时机不对")
		return false
	}
	if r.Location() == target.Location() {
		logger.Error("威逼不能对自己使用")
		return false
	}
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	switch wantType {
	case protos.CardType_Cheng_Qing, protos.CardType_Jie_Huo, protos.CardType_Diao_Bao, protos.CardType_Wu_Dao:
	default:
		logger.Error("威逼选择的卡牌类型错误：", protos.CardType_name[int32(wantType)])
		return false
	}
	return true
}

func (card *WeiBi) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	wantType := args[1].(protos.CardType)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		if func(player game.IPlayer, cardType protos.CardType) bool {
			for _, c := range player.GetCards() {
				if c.GetType() == cardType {
					return true
				}
			}
			return false
		}(target, wantType) {
			return &executeWeiBi{
				player:   r,
				target:   target,
				card:     card,
				wantType: wantType,
			}, true
		} else {
			logger.Info(target, "向", r, "展示了所有手牌")
			g.GetDeck().Discard(card)
			for _, p := range g.GetPlayers() {
				if player, ok := p.(*game.HumanPlayer); ok {
					msg := &protos.WeiBiShowHandCardToc{
						Card:           card.ToPbCard(),
						PlayerId:       p.GetAlternativeLocation(r.Location()),
						WantType:       wantType,
						TargetPlayerId: p.GetAlternativeLocation(target.Location()),
					}
					if p.Location() == r.Location() {
						for _, c := range target.GetCards() {
							msg.Cards = append(msg.Cards, c.ToPbCard())
						}
					}
					player.Send(msg)
				}
			}
			return &game.MainPhaseIdle{Player: r}, true
		}
	}
	g.Resolve(&game.OnUseCard{
		WhoseTurn:   r,
		Player:      r,
		Card:        card,
		AskWhom:     r,
		ResolveFunc: resolveFunc,
	})
}

type executeWeiBi struct {
	player   game.IPlayer
	target   game.IPlayer
	card     *WeiBi
	wantType protos.CardType
}

func (e *executeWeiBi) Resolve() (next game.Fsm, continueResolve bool) {
	r, target, card, wantType := e.player, e.target, e.card, e.wantType
	g := r.GetGame()
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			msg := &protos.WeiBiWaitForGiveCardToc{
				Card:           card.ToPbCard(),
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
				WantType:       wantType,
				WaitingSecond:  20,
			}
			if p.Location() == target.Location() {
				seq := player.Seq
				msg.Seq = player.Seq
				player.Timer = time.AfterFunc(time.Second*time.Duration(msg.WaitingSecond+2), func() {
					game.Post(func() {
						if player.Seq == seq {
							player.IncrSeq()
							e.autoSelect()
							g.Resolve(&game.MainPhaseIdle{Player: e.player})
						}
					})
				})
			}
			player.Send(msg)
		}
	}
	if _, ok := target.(*game.RobotPlayer); ok {
		time.AfterFunc(2*time.Second, func() {
			game.Post(func() {
				e.autoSelect()
				g.Resolve(&game.MainPhaseIdle{Player: e.player})
			})
		})
	}
	return e, false
}

func (e *executeWeiBi) ResolveProtocol(target game.IPlayer, pb proto.Message) (next game.Fsm, continueResolve bool) {
	msg, ok := pb.(*protos.WeiBiGiveCardTos)
	if !ok {
		logger.Error("现在正在结算威逼")
		return e, false
	}
	if target.Location() != e.target.Location() {
		logger.Error("你不是威逼的目标", e.card)
		return e, false
	}
	cardId := msg.CardId
	c := target.FindCard(cardId)
	if c == nil {
		logger.Error("没有这张牌")
		return e, false
	}
	target.IncrSeq()
	r, card := e.player, e.card
	g := r.GetGame()
	logger.Info(target, "给了", r, "一张", c)
	target.DeleteCard(cardId)
	r.AddCards(c)
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.WeiBiGiveCardToc{
				PlayerId:       player.GetAlternativeLocation(r.Location()),
				TargetPlayerId: player.GetAlternativeLocation(target.Location()),
				Card:           c.ToPbCard(),
			})
		}
	}
	g.GetDeck().Discard(card)
	return &game.MainPhaseIdle{Player: e.player}, true
}

func (e *executeWeiBi) autoSelect() {
	var availableCards []uint32
	for _, c := range e.target.GetCards() {
		if c.GetType() == e.wantType {
			availableCards = append(availableCards, c.GetId())
		}
	}
	e.ResolveProtocol(e.target, &protos.WeiBiGiveCardTos{CardId: availableCards[utils.Random.Intn(len(availableCards))]})
}

func (card *WeiBi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
