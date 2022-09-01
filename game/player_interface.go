package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type IPlayer interface {
	Init(game *Game, location int, identity protos.Color, secretTask protos.SecretTask, roleSkillsData *RoleSkillsData)
	IncrSeq()
	GetGame() *Game
	Location() int
	GetAbstractLocation(location int) int
	GetAlternativeLocation(location int) uint32
	NotifyAddHandCard(location int, unknownCount int, cards ...ICard)
	Draw(count int)
	AddCards(cards ...ICard)
	GetCards() map[uint32]ICard
	FindCard(cardId uint32) ICard
	DeleteCard(card uint32)
	DeleteAllCards() []ICard
	AddMessageCards(cards ...ICard)
	GetMessageCards() map[uint32]ICard
	FindMessageCard(cardId uint32) ICard
	DeleteMessageCard(cardId uint32)
	DeleteAllMessageCards() []ICard
	CheckThreeSameMessageCard(colors ...protos.Color) bool
	NotifyDrawPhase(player IPlayer)
	NotifyMainPhase(player IPlayer, waitSecond uint32)
	NotifySendPhaseStart(player IPlayer, waitSecond uint32)
	NotifySendMessageCard(player, targetPlayer IPlayer, lockedPlayers []IPlayer, messageCard ICard, dir protos.Direction)
	NotifySendPhase(whoseTurn, inFrontOfWhom IPlayer, lockedPlayers []IPlayer, messageCard ICard, dir protos.Direction, isMessageFaceUp bool, waitSecond uint32)
	NotifyChooseReceiveCard(player IPlayer)
	NotifyFightPhase(whoseTurn, inFrontOfWhom, whoseFightTurn IPlayer, messageCard ICard, isMessageFaceUp bool, waitSecond uint32)
	NotifyReceivePhase(whoseTurn, inFrontOfWhom IPlayer, messageCard ICard)
	NotifyDie(location int, loseGame bool)
	NotifyWin(declareWinner []IPlayer, winner []IPlayer)
	NotifyAskForChengQing(whoDie IPlayer, askWhom IPlayer)
	WaitForDieGiveCard(whoDie IPlayer)
	SetAlive(alive bool)
	IsAlive() bool
	SetLose(lose bool)
	IsLose() bool
	SetHasNoIdentity(hasNoIdentity bool)
	HasNoIdentity() bool
	SetIdentity(identity protos.Color, secretTask protos.SecretTask)
	GetIdentity() (protos.Color, protos.SecretTask)
	SetSkills(skills []ISkill)
	GetSkills() []ISkill
	FindSkill(skillId SkillId) ISkill
	GetRole() protos.Role
	IsRoleFaceUp() bool
	AddSkillUseCount(skillId SkillId)
	GetSkillUseCount(skillId SkillId) int
	ResetSkillUseCount()
	String() string
}

type BasePlayer struct {
	game           *Game
	location       int
	cards          map[uint32]ICard
	messageCards   map[uint32]ICard
	identity       protos.Color
	secretTask     protos.SecretTask
	alive          bool
	lose           bool
	hasNoIdentity  bool
	roleSkillsData RoleSkillsData
	skillUseCount  map[SkillId]int
}

func (p *BasePlayer) Init(game *Game, location int, identity protos.Color, secretTask protos.SecretTask, roleSkillsData *RoleSkillsData) {
	logger.Info(location, "号的身份是", utils.IdentityColorToString(identity, secretTask))
	p.game = game
	p.location = location
	p.cards = make(map[uint32]ICard)
	p.messageCards = make(map[uint32]ICard)
	p.alive = true
	p.roleSkillsData = *roleSkillsData
	for _, skill := range p.roleSkillsData.Skills {
		skill.Init(game)
	}
	p.ResetSkillUseCount()
	p.SetIdentity(identity, secretTask)
}

func (p *BasePlayer) IncrSeq() {
}

func (p *BasePlayer) GetGame() *Game {
	return p.game
}

func (p *BasePlayer) Location() int {
	return p.location
}

func (p *BasePlayer) NotifyAddHandCard(int, int, ...ICard) {
}

func (p *BasePlayer) Draw(count int) {
	cards := p.game.GetDeck().Draw(count)
	p.AddCards(cards...)
	logger.Info(p.Location(), "号摸了", cards, "，现在有", len(p.cards), "张手牌")
	for _, player := range p.game.GetPlayers() {
		if player.Location() == p.Location() {
			player.NotifyAddHandCard(p.Location(), 0, cards...)
		} else {
			player.NotifyAddHandCard(p.Location(), len(cards))
		}
	}
}

func (p *BasePlayer) GetAbstractLocation(location int) int {
	return (location + p.Location()) % len(p.GetGame().GetPlayers())
}

func (p *BasePlayer) GetAlternativeLocation(location int) uint32 {
	location -= p.Location()
	totalPlayerCount := len(p.GetGame().GetPlayers())
	if location < 0 {
		location += totalPlayerCount
	}
	return uint32(location)
}

func (p *BasePlayer) SetAlive(alive bool) {
	p.alive = alive
}

func (p *BasePlayer) IsAlive() bool {
	return p.alive
}

func (p *BasePlayer) SetLose(lose bool) {
	p.lose = lose
}

func (p *BasePlayer) IsLose() bool {
	return p.lose
}

func (p *BasePlayer) SetHasNoIdentity(hasNoIdentity bool) {
	p.hasNoIdentity = hasNoIdentity
}

func (p *BasePlayer) HasNoIdentity() bool {
	return p.hasNoIdentity
}

func (p *BasePlayer) AddCards(cards ...ICard) {
	for _, card := range cards {
		p.cards[card.GetId()] = card
	}
}

func (p *BasePlayer) GetCards() map[uint32]ICard {
	return p.cards
}

func (p *BasePlayer) FindCard(cardId uint32) ICard {
	if card, ok := p.cards[cardId]; ok {
		return card
	}
	return nil
}

func (p *BasePlayer) DeleteCard(cardId uint32) {
	delete(p.cards, cardId)
}

func (p *BasePlayer) DeleteAllCards() []ICard {
	var cards []ICard
	for _, card := range p.cards {
		cards = append(cards, card)
	}
	p.cards = make(map[uint32]ICard)
	return cards
}

func (p *BasePlayer) AddMessageCards(cards ...ICard) {
	for _, card := range cards {
		p.messageCards[card.GetId()] = card
	}
}

func (p *BasePlayer) GetMessageCards() map[uint32]ICard {
	return p.messageCards
}

func (p *BasePlayer) FindMessageCard(cardId uint32) ICard {
	if card, ok := p.messageCards[cardId]; ok {
		return card
	}
	return nil
}

func (p *BasePlayer) DeleteMessageCard(cardId uint32) {
	delete(p.messageCards, cardId)
}

func (p *BasePlayer) DeleteAllMessageCards() []ICard {
	var cards []ICard
	for _, card := range p.messageCards {
		cards = append(cards, card)
	}
	p.messageCards = make(map[uint32]ICard)
	return cards
}

func (p *BasePlayer) CheckThreeSameMessageCard(colors ...protos.Color) bool {
	for _, color := range colors {
		var count int
		for _, card := range p.messageCards {
			if utils.IsColorIn(color, card.GetColors()) {
				count++
			}
		}
		if count >= 3 {
			return true
		}
	}
	return false
}

func (p *BasePlayer) SetIdentity(identity protos.Color, secretTask protos.SecretTask) {
	p.identity = identity
	p.secretTask = secretTask
}

func (p *BasePlayer) GetIdentity() (protos.Color, protos.SecretTask) {
	return p.identity, p.secretTask
}

func (p *BasePlayer) SetSkills(skills []ISkill) {
	p.roleSkillsData.Skills = skills
}

func (p *BasePlayer) GetSkills() []ISkill {
	return p.roleSkillsData.Skills
}

func (p *BasePlayer) FindSkill(skillId SkillId) ISkill {
	for _, skill := range p.roleSkillsData.Skills {
		if skill.GetSkillId() == skillId {
			return skill
		}
	}
	return nil
}

func (p *BasePlayer) GetRole() protos.Role {
	return p.roleSkillsData.Role
}

func (p *BasePlayer) IsRoleFaceUp() bool {
	return p.roleSkillsData.FaceUp
}

func (p *BasePlayer) AddSkillUseCount(skillId SkillId) {
	p.skillUseCount[skillId]++
}

func (p *BasePlayer) GetSkillUseCount(skillId SkillId) int {
	return p.skillUseCount[skillId]
}

func (p *BasePlayer) ResetSkillUseCount() {
	p.skillUseCount = make(map[SkillId]int)
}
