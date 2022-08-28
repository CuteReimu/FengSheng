package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"strconv"
	"time"
)

var AIMainPhase = make(map[protos.CardType]func(player interfaces.IPlayer, card interfaces.ICard) bool)
var AISendPhase = make(map[protos.CardType]func(player interfaces.IPlayer, card interfaces.ICard) bool)
var AIFightPhase = make(map[protos.CardType]func(player interfaces.IPlayer, card interfaces.ICard) bool)

type RobotPlayer struct {
	interfaces.BasePlayer
}

func (r *RobotPlayer) String() string {
	return strconv.Itoa(r.Location()) + "号[机器人]"
}

func (r *RobotPlayer) NotifyDrawPhase() {
}

func (r *RobotPlayer) NotifyMainPhase(uint32) {
	if r.Location() != r.GetGame().GetWhoseTurn() {
		return
	}
	cards := r.GetCards()
	if len(cards) > 1 {
		for cardId := range cards {
			card := cards[cardId]
			ai := AIMainPhase[card.GetType()]
			if ai != nil && ai(r, card) {
				return
			}
		}
	}
	time.AfterFunc(2*time.Second, func() {
		Post(r.GetGame().SendPhaseStart)
	})
}

func (r *RobotPlayer) NotifySendPhaseStart(uint32) {
	if r.Location() != r.GetGame().GetWhoseTurn() {
		return
	}
	time.AfterFunc(2*time.Second, func() {
		autoSendMessageCard(r, true)
	})
}

func (r *RobotPlayer) NotifySendPhase(uint32, bool) {
	if r.Location() != r.GetGame().GetWhoseSendTurn() {
		return
	}
	cards := r.GetCards()
	for cardId := range cards {
		card := cards[cardId]
		ai := AISendPhase[card.GetType()]
		if ai != nil && ai(r, card) {
			return
		}
	}
	time.AfterFunc(2*time.Second, func() {
		colors := r.GetGame().GetCurrentMessageCard().GetColor()
		certainlyReceive := r.GetGame().IsMessageCardFaceUp() && len(colors) == 1 && colors[0] != protos.Color_Black
		certainlyReject := r.GetGame().IsMessageCardFaceUp() && len(colors) == 1 && colors[0] == protos.Color_Black
		if certainlyReceive || func(e int, arr []int) bool {
			for _, a := range arr {
				if a == e {
					return true
				}
			}
			return false
		}(r.Location(), r.GetGame().GetLockPlayers()) || r.Location() == r.GetGame().GetWhoseTurn() || (!certainlyReject && utils.Random.Intn((len(r.GetGame().GetPlayers())-1)*2) == 0) {
			Post(r.GetGame().OnChooseReceiveCard)
		} else {
			Post(r.GetGame().MessageMoveNext)
		}
	})
}

func (r *RobotPlayer) NotifyChooseReceiveCard() {
}

func (r *RobotPlayer) NotifyFightPhase(uint32) {
	if r.Location() != r.GetGame().GetWhoseFightTurn() {
		return
	}
	cards := r.GetCards()
	for cardId := range cards {
		card := cards[cardId]
		ai := AIFightPhase[card.GetType()]
		if ai != nil && ai(r, card) {
			return
		}
	}
	time.AfterFunc(2*time.Second, func() {
		Post(r.GetGame().FightPhaseNext)
	})
}

func (r *RobotPlayer) NotifyReceivePhase() {
}

func (r *RobotPlayer) NotifyDie(location int, _ bool) {
	if location == r.Location() {
		r.SetAlive(false)
		var cards []interfaces.ICard
		for _, card := range r.GetCards() {
			cards = append(cards, card)
		}
		r.GetGame().PlayerDiscardCard(r, cards...)
		r.GetGame().GetDeck().Discard(r.DeleteAllMessageCards()...)
	}
}

func (r *RobotPlayer) NotifyWin(interfaces.IPlayer, []interfaces.IPlayer) {
}

func (r *RobotPlayer) NotifyAskForChengQing(_ interfaces.IPlayer, askWhom interfaces.IPlayer) {
	if askWhom.Location() != r.Location() {
		return
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			Post(r.GetGame().AskNextForChengQing)
		})
	})
}

func (r *RobotPlayer) WaitForDieGiveCard(whoDie interfaces.IPlayer) {
	if whoDie.Location() != r.Location() {
		return
	}

	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			if r.Location() != r.GetGame().GetWhoDie() {
				return
			}
			identity1, _ := r.GetIdentity()
			if identity1 != protos.Color_Black {
				for _, p := range r.GetGame().GetPlayers() {
					if identity2, _ := r.GetIdentity(); identity1 == identity2 && p.Location() != r.Location() {
						var cards []interfaces.ICard
						for _, card := range r.GetCards() {
							cards = append(cards, card)
							if len(cards) >= 3 {
								break
							}
						}
						if len(cards) > 0 {
							for _, card := range cards {
								r.DeleteCard(card.GetId())
							}
							target := r.GetGame().GetPlayers()[p.Location()]
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
						}
						break
					}
				}
			}
			for _, p := range r.GetGame().GetPlayers() {
				p.NotifyDie(whoDie.Location(), false)
			}
			Post(r.GetGame().AfterChengQing)
		})
	})
}

func autoSendMessageCard(r interfaces.IPlayer, lock bool) {
	var card interfaces.ICard
	for _, card = range r.GetCards() {
		break
	}
	dir := card.GetDirection()
	var targetLocation int
	var lockLocation, availableLocations []int
	for _, p := range r.GetGame().GetPlayers() {
		if p.Location() != r.Location() && p.IsAlive() {
			availableLocations = append(availableLocations, p.Location())
		}
	}
	if dir != protos.Direction_Up && lock && card.CanLock() && utils.Random.Intn(3) != 0 {
		location := availableLocations[utils.Random.Intn(len(availableLocations))]
		if r.GetGame().GetPlayers()[location].IsAlive() {
			lockLocation = append(lockLocation, location)
		}
	}
	switch dir {
	case protos.Direction_Up:
		targetLocation = availableLocations[utils.Random.Intn(len(availableLocations))]
		if lock && card.CanLock() && utils.Random.Intn(2) != 0 {
			lockLocation = append(lockLocation, targetLocation)
		}
	case protos.Direction_Left:
		targetLocation = (r.Location() + len(r.GetGame().GetPlayers()) - 1) % len(r.GetGame().GetPlayers())
		for !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
			targetLocation = (targetLocation + len(r.GetGame().GetPlayers()) - 1) % len(r.GetGame().GetPlayers())
		}
	case protos.Direction_Right:
		targetLocation = (r.Location() + 1) % len(r.GetGame().GetPlayers())
		for !r.GetGame().GetPlayers()[targetLocation].IsAlive() {
			targetLocation = (targetLocation + 1) % len(r.GetGame().GetPlayers())
		}
	}
	Post(func() { r.GetGame().OnSendCard(card, dir, targetLocation, lockLocation) })
}

type IdlePlayer struct {
	interfaces.BasePlayer
}

func (p *IdlePlayer) NotifyDrawPhase() {
}

func (p *IdlePlayer) NotifyMainPhase(uint32) {
}

func (p *IdlePlayer) NotifySendPhaseStart(uint32) {
}

func (p *IdlePlayer) NotifySendPhase(uint32, bool) {
}

func (p *IdlePlayer) NotifyChooseReceiveCard() {
}

func (p *IdlePlayer) NotifyFightPhase(uint32) {
}

func (p *IdlePlayer) NotifyReceivePhase() {
}

func (p *IdlePlayer) NotifyDie(int, bool) {
}

func (p *IdlePlayer) NotifyWin(interfaces.IPlayer, []interfaces.IPlayer) {
}

func (p *IdlePlayer) NotifyAskForChengQing(interfaces.IPlayer, interfaces.IPlayer) {
}

func (p *IdlePlayer) WaitForDieGiveCard(interfaces.IPlayer) {
}

func (p *IdlePlayer) String() string {
	return ""
}
