package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"strconv"
	"time"
)

var AISkillMainPhase = make(map[SkillId]func(e *MainPhaseIdle, skill ISkill) bool)
var AIMainPhase = make(map[protos.CardType]func(e *MainPhaseIdle, card ICard) bool)
var AISendPhase = make(map[protos.CardType]func(e *SendPhaseIdle, card ICard) bool)
var AIFightPhase = make(map[protos.CardType]func(e *FightPhaseIdle, card ICard) bool)

type RobotPlayer struct {
	BasePlayer
}

func (r *RobotPlayer) String() string {
	if r.roleSkillsData.FaceUp {
		return strconv.Itoa(r.Location()) + "号[" + r.roleSkillsData.Name + "]"
	}
	return strconv.Itoa(r.Location()) + "号[机器人]"
}

func (r *RobotPlayer) NotifyDrawPhase(IPlayer) {
}

func (r *RobotPlayer) NotifyMainPhase(player IPlayer, _ uint32) {
	fsm := r.GetGame().GetFsm().(*MainPhaseIdle)
	if r.Location() != player.Location() {
		return
	}
	for _, skill := range r.GetSkills() {
		ai := AISkillMainPhase[skill.GetSkillId()]
		if ai != nil && ai(fsm, skill) {
			return
		}
	}
	cards := r.GetCards()
	if len(cards) > 1 {
		for cardId := range cards {
			card := cards[cardId]
			ai := AIMainPhase[card.GetType()]
			if ai != nil && ai(fsm, card) {
				return
			}
		}
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			r.GetGame().Resolve(&SendPhaseStart{Player: r})
		})
	})
}

func (r *RobotPlayer) NotifySendPhaseStart(IPlayer, uint32) {
	fsm := r.GetGame().GetFsm().(*SendPhaseStart)
	if r.Location() != fsm.Player.Location() {
		return
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			autoSendMessageCard(r, true)
		})
	})
}

func (r *RobotPlayer) NotifySendMessageCard(IPlayer, IPlayer, []IPlayer, ICard, protos.Direction) {
}

func (r *RobotPlayer) NotifySendPhase(whoseTurn, inFrontOfWhom IPlayer, lockedPlayers []IPlayer, messageCard ICard, _ protos.Direction, isMessageFaceUp bool, _ uint32) {
	fsm := r.GetGame().GetFsm().(*SendPhaseIdle)
	if r.Location() != inFrontOfWhom.Location() {
		return
	}
	cards := r.GetCards()
	for cardId := range cards {
		card := cards[cardId]
		ai := AISendPhase[card.GetType()]
		if ai != nil && ai(fsm, card) {
			return
		}
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			colors := messageCard.GetColors()
			certainlyReceive := isMessageFaceUp && len(colors) == 1 && colors[0] != protos.Color_Black
			certainlyReject := isMessageFaceUp && len(colors) == 1 && colors[0] == protos.Color_Black
			if certainlyReceive || func(e int, lockedPlayers []IPlayer) bool {
				for _, a := range lockedPlayers {
					if a.Location() == e {
						return true
					}
				}
				return false
			}(r.Location(), lockedPlayers) || r.Location() == whoseTurn.Location() || (!certainlyReject && utils.Random.Intn((len(r.GetGame().GetPlayers())-1)*2) == 0) {
				r.GetGame().Resolve(&OnChooseReceiveCard{
					WhoseTurn:           whoseTurn,
					MessageCard:         messageCard,
					InFrontOfWhom:       inFrontOfWhom,
					IsMessageCardFaceUp: isMessageFaceUp,
				})
			} else {
				r.GetGame().Resolve(&MessageMoveNext{SendPhase: fsm})
			}
		})
	})
}

func (r *RobotPlayer) NotifyChooseReceiveCard(IPlayer) {
}

func (r *RobotPlayer) NotifyFightPhase(_, _, whoseFightTurn IPlayer, _ ICard, _ bool, _ uint32) {
	fsm := r.GetGame().GetFsm().(*FightPhaseIdle)
	if r.Location() != whoseFightTurn.Location() {
		return
	}
	cards := r.GetCards()
	for cardId := range cards {
		card := cards[cardId]
		ai := AIFightPhase[card.GetType()]
		if ai != nil && ai(fsm, card) {
			return
		}
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			r.GetGame().Resolve(&FightPhaseNext{FightPhase: fsm})
		})
	})
}

func (r *RobotPlayer) NotifyReceivePhase(IPlayer, IPlayer, ICard) {
}

func (r *RobotPlayer) NotifyReceivePhaseWithWaiting(_, _ IPlayer, _ ICard, waitingPlayer IPlayer, _ uint32) {
	if r.Location() != waitingPlayer.Location() {
		return
	}
	// TODO 需要增加AI
	time.AfterFunc(2*time.Second, func() {
		r.GetGame().TryContinueResolveProtocol(r, &protos.EndReceivePhaseTos{})
	})
}

func (r *RobotPlayer) NotifyDie(location int, _ bool) {
	if location == r.Location() {
		var cards []ICard
		for _, card := range r.GetCards() {
			cards = append(cards, card)
		}
		r.GetGame().PlayerDiscardCard(r, cards...)
		r.GetGame().GetDeck().Discard(r.DeleteAllMessageCards()...)
	}
}

func (r *RobotPlayer) NotifyWin([]IPlayer, []IPlayer) {
}

func (r *RobotPlayer) NotifyAskForChengQing(_ IPlayer, askWhom IPlayer) {
	fsm := r.GetGame().GetFsm().(*WaitForChengQing)
	if askWhom.Location() != r.Location() {
		return
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			r.GetGame().Resolve(&WaitNextForChengQing{WaitForChengQing: fsm})
		})
	})
}

func (r *RobotPlayer) WaitForDieGiveCard(whoDie IPlayer) {
	fsm := r.GetGame().GetFsm().(*WaitForDieGiveCard)
	if whoDie.Location() != r.Location() {
		return
	}
	time.AfterFunc(2*time.Second, func() {
		Post(func() {
			identity1, _ := r.GetIdentity()
			if identity1 != protos.Color_Black {
				for _, p := range r.GetGame().GetPlayers() {
					if identity2, _ := p.GetIdentity(); identity1 == identity2 && p.Location() != r.Location() {
						var cards []ICard
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
							target := p
							target.AddCards(cards...)
							logger.Info(r, "给了", target, cards)
							for _, p := range r.GetGame().GetPlayers() {
								if player, ok := p.(*HumanPlayer); ok {
									msg := &protos.NotifyDieGiveCardToc{
										PlayerId:       p.GetAlternativeLocation(r.Location()),
										TargetPlayerId: p.GetAlternativeLocation(target.Location()),
										CardCount:      uint32(len(cards)),
									}
									if p.Location() == target.Location() {
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
			r.GetGame().Resolve(&AfterDieGiveCard{DieGiveCard: fsm})
		})
	})
}

func autoSendMessageCard(r IPlayer, lock bool) {
	var card ICard
	for _, card = range r.GetCards() {
		break
	}
	fsm := r.GetGame().GetFsm().(*SendPhaseStart)
	dir := card.GetDirection()
	if r.FindSkill(SkillIdLianLuo) != nil {
		dir = protos.Direction(utils.Random.Intn(len(protos.Direction_name)))
	}
	var targetLocation int
	var availableLocations []int
	var lockedPlayers []IPlayer
	for _, p := range r.GetGame().GetPlayers() {
		if p.Location() != r.Location() && p.IsAlive() {
			availableLocations = append(availableLocations, p.Location())
		}
	}
	if dir != protos.Direction_Up && lock && card.CanLock() && utils.Random.Intn(3) != 0 {
		location := availableLocations[utils.Random.Intn(len(availableLocations))]
		if r.GetGame().GetPlayers()[location].IsAlive() {
			lockedPlayers = append(lockedPlayers, r.GetGame().GetPlayers()[location])
		}
	}
	switch dir {
	case protos.Direction_Up:
		targetLocation = availableLocations[utils.Random.Intn(len(availableLocations))]
		if lock && card.CanLock() && utils.Random.Intn(2) != 0 {
			lockedPlayers = append(lockedPlayers, r.GetGame().GetPlayers()[targetLocation])
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
	r.GetGame().Resolve(&OnSendCard{
		WhoseTurn:     fsm.Player,
		MessageCard:   card,
		Dir:           dir,
		TargetPlayer:  r.GetGame().GetPlayers()[targetLocation],
		LockedPlayers: lockedPlayers,
	})
}

type IdlePlayer struct {
	BasePlayer
}

func (p *IdlePlayer) NotifyDrawPhase(IPlayer) {
}

func (p *IdlePlayer) NotifyMainPhase(IPlayer, uint32) {
}

func (p *IdlePlayer) NotifySendPhaseStart(IPlayer, uint32) {
}

func (p *IdlePlayer) NotifySendMessageCard(IPlayer, IPlayer, []IPlayer, ICard, protos.Direction) {
}

func (p *IdlePlayer) NotifySendPhase(IPlayer, IPlayer, []IPlayer, ICard, protos.Direction, bool, uint32) {
}

func (p *IdlePlayer) NotifyChooseReceiveCard(IPlayer) {
}

func (p *IdlePlayer) NotifyFightPhase(IPlayer, IPlayer, IPlayer, ICard, bool, uint32) {
}

func (p *IdlePlayer) NotifyReceivePhase(IPlayer, IPlayer, ICard) {
}

func (p *IdlePlayer) NotifyReceivePhaseWithWaiting(IPlayer, IPlayer, ICard, IPlayer, uint32) {
}

func (p *IdlePlayer) NotifyDie(int, bool) {
}

func (p *IdlePlayer) NotifyWin([]IPlayer, []IPlayer) {
}

func (p *IdlePlayer) NotifyAskForChengQing(IPlayer, IPlayer) {
}

func (p *IdlePlayer) WaitForDieGiveCard(IPlayer) {
}

func (p *IdlePlayer) String() string {
	return ""
}
