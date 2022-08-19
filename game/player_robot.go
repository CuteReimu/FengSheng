package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"strconv"
	"time"
)

var AIMainPhase = make(map[protos.CardType]func(player interfaces.IPlayer, card interfaces.ICard) bool)

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
	for cardId := range cards {
		card := cards[cardId]
		ai := AIMainPhase[card.GetType()]
		if ai != nil && ai(r, card) {
			return
		}
	}
	time.AfterFunc(time.Second, func() {
		r.GetGame().Post(r.GetGame().SendPhaseStart)
	})
}

func (r *RobotPlayer) NotifySendPhaseStart(uint32) {
	if r.Location() != r.GetGame().GetWhoseTurn() {
		return
	}
	time.AfterFunc(time.Second, func() {
		autoSendMessageCard(r, true)
	})
}

func (r *RobotPlayer) NotifySendPhase(uint32, bool) {
	if r.Location() != r.GetGame().GetWhoseSendTurn() {
		return
	}
	time.AfterFunc(time.Second, func() {
		if func(e int, arr []int) bool {
			for _, a := range arr {
				if a == e {
					return true
				}
			}
			return false
		}(r.Location(), r.GetGame().GetLockPlayers()) || r.Location() == r.GetGame().GetWhoseTurn() || r.GetGame().GetRandom().Intn((len(r.GetGame().GetPlayers())-1)*2) == 0 {
			r.GetGame().Post(r.GetGame().OnChooseReceiveCard)
		} else {
			r.GetGame().Post(r.GetGame().MessageMoveNext)
		}
	})
}

func (r *RobotPlayer) NotifyFightPhase(uint32) {
	if r.Location() != r.GetGame().GetWhoseFightTurn() {
		return
	}
	time.AfterFunc(time.Second/2, func() {
		r.GetGame().Post(r.GetGame().FightPhaseNext)
	})
}

func (r *RobotPlayer) NotifyReceivePhase() {
}

func (r *RobotPlayer) NotifyDie(int, bool) {
}

func (r *RobotPlayer) NotifyWin(interfaces.IPlayer, []interfaces.IPlayer) {

}

func (r *RobotPlayer) NotifyAskForChengQing(_ interfaces.IPlayer, askWhom interfaces.IPlayer) {
	if askWhom.Location() != r.Location() {
		return
	}
	time.AfterFunc(time.Second/2, func() {
		r.GetGame().Post(func() {
			r.GetGame().Post(r.GetGame().AskNextForChengQing)
		})
	})
}

func (r *RobotPlayer) WaitForDieGiveCard(whoDie interfaces.IPlayer) {
	if whoDie.Location() != r.Location() {
		return
	}
	time.AfterFunc(time.Second, func() {
		r.GetGame().Post(func() {
			r.GetGame().Post(r.GetGame().AfterChengQing)
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
	if lock && card.CanLock() && r.GetGame().GetRandom().Intn(2) == 0 {
		lockLocation = append(lockLocation, availableLocations[r.GetGame().GetRandom().Intn(len(availableLocations))])
	}
	switch dir {
	case protos.Direction_Up:
		targetLocation = availableLocations[r.GetGame().GetRandom().Intn(len(availableLocations))]
	case protos.Direction_Left:
		targetLocation = (r.Location() + len(r.GetGame().GetPlayers()) - 1) % len(r.GetGame().GetPlayers())
	case protos.Direction_Right:
		targetLocation = (r.Location() + 1) % len(r.GetGame().GetPlayers())
	}
	r.GetGame().Post(func() { r.GetGame().OnSendCard(card, dir, targetLocation, lockLocation) })
}
