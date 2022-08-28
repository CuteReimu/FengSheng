package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

type WeiBi struct {
	interfaces.BaseCard
}

func (card *WeiBi) String() string {
	return utils.CardColorToString(card.Color...) + "威逼"
}

func (card *WeiBi) GetType() protos.CardType {
	return protos.CardType_Wei_Bi
}

func (card *WeiBi) CanUse(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	target := args[0].(interfaces.IPlayer)
	wantType := args[1].(protos.CardType)
	if g.GetCurrentPhase() != protos.Phase_Main_Phase || g.GetWhoseTurn() != r.Location() || !g.IsIdleTimePoint() {
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
	}
	return true
}

func (card *WeiBi) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	wantType := args[1].(protos.CardType)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	if func(player interfaces.IPlayer, cardType protos.CardType) bool {
		for _, c := range player.GetCards() {
			if c.GetType() == cardType {
				return true
			}
		}
		return false
	}(target, wantType) {
		g.SetCurrentCard(&interfaces.CurrentCard{Card: card, Player: r.Location(), TargetPlayer: target.Location()})
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
								player.Seq++
								if player.Timer != nil {
									player.Timer.Stop()
								}
								card.autoSelect(g, r, player, wantType)
							}
						})
					})
				}
				player.Send(msg)
			}
		}
		if _, ok := target.(*game.RobotPlayer); ok {
			time.AfterFunc(2*time.Second, func() {
				game.Post(func() { card.autoSelect(g, r, target, wantType) })
			})
		}
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
		if _, ok := r.(*game.RobotPlayer); ok {
			time.AfterFunc(2*time.Second, func() {
				game.Post(g.MainPhase)
			})
		} else {
			game.Post(g.MainPhase)
		}
	}
}

func (card *WeiBi) CanUse2(_ interfaces.IGame, _ interfaces.IPlayer, args ...interface{}) bool {
	target := args[0].(interfaces.IPlayer)
	cardId := args[1].(uint32)
	if target.FindCard(cardId) == nil {
		logger.Error("没有这张牌")
		return false
	}
	return true
}

func (card *WeiBi) Execute2(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	cardId := args[1].(uint32)
	c := target.FindCard(cardId)
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
	g.SetCurrentCard(nil)
	g.GetDeck().Discard(card)
	game.Post(g.MainPhase)
}

func (card *WeiBi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}

func (card *WeiBi) autoSelect(g interfaces.IGame, player, target interfaces.IPlayer, cardType protos.CardType) {
	var availableCards []uint32
	for _, c := range target.GetCards() {
		if c.GetType() == cardType {
			availableCards = append(availableCards, c.GetId())
		}
	}
	card.Execute2(g, player, target, availableCards[utils.Random.Intn(len(availableCards))])
}
