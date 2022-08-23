package game

import (
	"github.com/CuteReimu/FengSheng/config"
	_ "github.com/CuteReimu/FengSheng/core"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

var logger = utils.GetLogger("game")

var eventQueue cellnet.EventQueue

func Post(callback func()) {
	eventQueue.Post(callback)
}

type Game struct {
	AlreadyStart       bool
	Players            []interfaces.IPlayer
	Deck               interfaces.IDeck
	CurrentCard        *interfaces.CurrentCard
	WhoseTurn          int
	WhoseSendTurn      int
	WhoseFightTurn     int
	WhoDie             int
	DieState           interfaces.DieState
	CurrentMessageCard interfaces.ICard
	MessageCardFaceUp  bool
	CardDirection      protos.Direction
	WhoIsLocked        []int
	CurrentPhase       protos.Phase
	afterChengQing     func()
}

func Start(totalCount int) {
	game := &Game{Players: make([]interfaces.IPlayer, totalCount)}

	if !config.IsTcpDebugLogOpen() {
		msglog.SetCurrMsgLogMode(msglog.MsgLogMode_Mute)
	}
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	eventQueue = cellnet.NewEventQueue()

	// 创建一个tcp的侦听器，名称为server，所有连接将事件投递到queue队列,单线程的处理
	p := peer.NewGenericPeer("tcp.Acceptor", "server", config.GetListenAddress(), eventQueue)

	humanMap := make(map[int64]*HumanPlayer)

	onAdd := func(player interfaces.IPlayer) int {
		playerIndex, unready := 0, -1
		for index := range game.Players {
			if game.Players[index] == nil {
				unready++
				if unready == 0 {
					game.Players[index] = player
					playerIndex = index
				}
			}
		}
		msg := &protos.JoinRoomToc{Name: player.String(), Position: uint32(playerIndex)}
		for i := range game.Players {
			if player, ok := game.Players[i].(*HumanPlayer); ok {
				player.Send(msg)
			}
		}
		if unready == 0 {
			logger.Info(player, "加入了。已加入", totalCount, "个人，游戏开始。。。")
			Post(game.start)
			game = &Game{Players: make([]interfaces.IPlayer, totalCount)}
		} else {
			logger.Info(player, "加入了。已加入", totalCount-unready, "个人，等待", unready, "人加入。。。")
		}
		return playerIndex
	}

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {
		switch pb := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			player := &HumanPlayer{Session: ev.Session()}
			humanMap[player.Session.ID()] = player
			index := onAdd(player)
			msg := &protos.GetRoomInfoToc{Names: make([]string, len(game.Players)), MyPosition: uint32(index)}
			for i := range game.Players {
				if game.Players[i] != nil {
					msg.Names[i] = game.Players[i].String()
				}
			}
			player.Send(msg)
		case *cellnet.SessionClosed:
			logger.Info("session closed: ", ev.Session().ID())
			if player, ok := humanMap[ev.Session().ID()]; ok {
				game := player.GetGame().(*Game)
				if game.AlreadyStart {
					game.Players[player.Location()] = &RobotPlayer{BasePlayer: player.BasePlayer}
				} else {
					msg := &protos.LeaveRoomToc{}
					for i := range game.Players {
						if player == game.Players[i] {
							msg.Position = uint32(i)
							game.Players[i] = nil
							break
						}
					}
					for i := range game.Players {
						if player, ok := game.Players[i].(*HumanPlayer); ok {
							player.Send(msg)
						}
					}
				}
				delete(humanMap, ev.Session().ID())
			}
		case *protos.AddRobotTos:
			onAdd(&RobotPlayer{})
		case *protos.EndMainPhaseTos:
			humanMap[ev.Session().ID()].onEndMainPhase(pb)
		case *protos.UseShiTanTos:
			humanMap[ev.Session().ID()].onUseShiTan(pb)
		case *protos.ExecuteShiTanTos:
			humanMap[ev.Session().ID()].onExecuteShiTan(pb)
		case *protos.UseLiYouTos:
			humanMap[ev.Session().ID()].onUseLiYou(pb)
		case *protos.UsePingHengTos:
			humanMap[ev.Session().ID()].onUsePingHeng(pb)
		case *protos.UseWeiBiTos:
			humanMap[ev.Session().ID()].onUseWeiBi(pb)
		case *protos.WeiBiGiveCardTos:
			humanMap[ev.Session().ID()].onWeiBiGiveCard(pb)
		case *protos.UseChengQingTos:
			humanMap[ev.Session().ID()].onUseChengQing(pb)
		case *protos.SendMessageCardTos:
			humanMap[ev.Session().ID()].onSendMessageCard(pb)
		case *protos.ChooseWhetherReceiveTos:
			humanMap[ev.Session().ID()].onChooseWhetherReceive(pb)
		case *protos.EndFightPhaseTos:
			humanMap[ev.Session().ID()].onEndFightPhase(pb)
		case *protos.ChengQingSaveDieTos:
			humanMap[ev.Session().ID()].onChengQingSaveDie(pb)
		case *protos.DieGiveCardTos:
			humanMap[ev.Session().ID()].onDieGiveCard(pb)
		case *protos.UsePoYiTos:
			humanMap[ev.Session().ID()].onUsePoYi(pb)
		case *protos.PoYiShowTos:
			humanMap[ev.Session().ID()].onPoYiShow(pb)
		case *protos.UseJieHuoTos:
			humanMap[ev.Session().ID()].onUseJieHuo(pb)
		case *protos.UseDiaoBaoTos:
			humanMap[ev.Session().ID()].onUseDiaoBao(pb)
		case *protos.UseWuDaoTos:
			humanMap[ev.Session().ID()].onUseWuDao(pb)
		}
	})
	p.Start()
	eventQueue.StartLoop()
	eventQueue.Wait()
}

func (game *Game) start() {
	game.AlreadyStart = true
	idTask := make([]struct {
		id   protos.Color
		task protos.SecretTask
	}, (len(game.Players)+1)/2*2+1)
	for i := 0; i < (len(game.Players)-1)/2; i++ {
		idTask[i*2].id = protos.Color_Red
		idTask[i*2+1].id = protos.Color_Blue
	}
	for i := 0; i < 3; i++ {
		idTask[len(idTask)-3+i].task = protos.SecretTask(i)
	}
	utils.Random.Shuffle(3, func(i, j int) {
		idTask[len(idTask)-3+i], idTask[len(idTask)-3+j] = idTask[len(idTask)-3+j], idTask[len(idTask)-3+i]
	})
	utils.Random.Shuffle(len(game.Players), func(i, j int) {
		idTask[i], idTask[j] = idTask[j], idTask[i]
	})
	for location, player := range game.Players {
		player.Init(game, location, idTask[location].id, idTask[location].task)
	}
	game.Deck = NewDeck(game)
	game.WhoseTurn = utils.Random.Intn(len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		game.Players[(game.WhoseTurn+i)%len(game.Players)].Draw(config.GetHandCardCountBegin())
	}
	game.DrawPhase()
}

func (game *Game) DrawPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	logger.Info(player, "的回合开始了")
	game.CurrentPhase = protos.Phase_Draw_Phase
	for _, p := range game.Players {
		p.NotifyDrawPhase()
	}
	player.Draw(config.GetHandCardCountEachTurn())
	Post(game.MainPhase)
}

func (game *Game) MainPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	game.CurrentPhase = protos.Phase_Main_Phase
	for _, p := range game.Players {
		p.NotifyMainPhase(30)
	}
}

func (game *Game) SendPhaseStart() {
	player := game.Players[game.WhoseTurn]
	if player.IsAlive() {
		if len(player.GetCards()) == 0 {
			logger.Info(player, "没有情报可传，输掉了游戏")
			game.GetDeck().Discard(player.DeleteAllMessageCards()...)
			player.SetLose(true)
			player.SetAlive(false)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(game.WhoseTurn, true)
			}
		}
	}
	if !player.IsAlive() {
		Post(game.NextTurn)
		return
	}
	game.CurrentPhase = protos.Phase_Send_Start_Phase
	for _, p := range game.Players {
		p.NotifySendPhaseStart(20)
	}
}

func (game *Game) OnSendCard(card interfaces.ICard, dir protos.Direction, targetLocation int, lockLocations []int) {
	player := game.Players[game.WhoseTurn]
	logger.Info(player, "传出了", card)
	player.DeleteCard(card.GetId())
	game.CurrentMessageCard = card
	game.CardDirection = dir
	game.WhoseSendTurn = targetLocation
	game.WhoIsLocked = append(game.WhoIsLocked, lockLocations...)
	game.CurrentPhase = protos.Phase_Send_Phase
	logger.Info("情报到达", game.Players[game.WhoseSendTurn], "面前")
	for _, p := range game.Players {
		p.NotifySendPhase(20, true)
	}
}

func (game *Game) MessageMoveNext() {
	if game.CardDirection == protos.Direction_Up {
		if game.Players[game.WhoseTurn].IsAlive() {
			game.WhoseSendTurn = game.WhoseTurn
		} else {
			Post(game.NextTurn)
			return
		}
	} else {
		for {
			if game.CardDirection == protos.Direction_Left {
				game.WhoseSendTurn = (game.WhoseSendTurn + len(game.Players) - 1) % len(game.Players)
			} else {
				game.WhoseSendTurn = (game.WhoseSendTurn + 1) % len(game.Players)
			}
			if game.Players[game.WhoseSendTurn].IsAlive() {
				break
			} else if game.WhoseTurn == game.WhoseSendTurn {
				Post(game.NextTurn)
				return
			}
		}
	}
	logger.Info("情报到达", game.Players[game.WhoseSendTurn], "面前")
	for _, p := range game.Players {
		p.NotifySendPhase(20, false)
	}
}

func (game *Game) OnChooseReceiveCard() {
	logger.Info(game.Players[game.WhoseSendTurn], "选择接收情报")
	game.WhoseFightTurn = game.WhoseSendTurn
	game.CurrentPhase = protos.Phase_Fight_Phase
	for _, p := range game.Players {
		p.NotifyChooseReceiveCard()
		p.NotifyFightPhase(20)
	}
}

func (game *Game) FightPhaseNext() {
	for {
		game.WhoseFightTurn = (game.WhoseFightTurn + 1) % len(game.Players)
		if game.WhoseFightTurn == game.WhoseSendTurn {
			Post(game.ReceivePhase)
			return
		} else if game.Players[game.WhoseFightTurn].IsAlive() {
			break
		}
	}
	for _, p := range game.Players {
		p.NotifyFightPhase(20)
	}
}

func (game *Game) ReceivePhase() {
	player := game.Players[game.WhoseSendTurn]
	if player.IsAlive() {
		game.CurrentPhase = protos.Phase_Receive_Phase
		player.AddMessageCards(game.CurrentMessageCard)
		logger.Info(player, "成功接收情报")
		for _, p := range game.Players {
			p.NotifyReceivePhase()
		}
		if game.checkWinOrDie() {
			return
		}
	}
	Post(game.NextTurn)
}

func (game *Game) NextTurn() {
	game.CurrentMessageCard = nil
	game.WhoIsLocked = nil
	for {
		game.WhoseTurn = (game.WhoseTurn + 1) % len(game.Players)
		if game.Players[game.WhoseTurn].IsAlive() {
			break
		}
	}
	Post(game.DrawPhase)
}

func (game *Game) GetPlayers() []interfaces.IPlayer {
	return game.Players
}

func (game *Game) GetDeck() interfaces.IDeck {
	return game.Deck
}

func (game *Game) GetWhoDie() int {
	return game.WhoDie
}

func (game *Game) GetWhoseTurn() int {
	return game.WhoseTurn
}

func (game *Game) GetWhoseSendTurn() int {
	return game.WhoseSendTurn
}

func (game *Game) SetWhoseSendTurn(whoseSendTurn int) {
	game.WhoseSendTurn = whoseSendTurn
}

func (game *Game) GetWhoseFightTurn() int {
	return game.WhoseFightTurn
}

func (game *Game) GetMessageCardDirection() protos.Direction {
	return game.CardDirection
}

func (game *Game) GetCurrentCard() *interfaces.CurrentCard {
	return game.CurrentCard
}

func (game *Game) SetCurrentCard(card *interfaces.CurrentCard) {
	game.CurrentCard = card
}

func (game *Game) GetCurrentMessageCard() interfaces.ICard {
	return game.CurrentMessageCard
}

func (game *Game) SetCurrentMessageCard(currentMessageCard interfaces.ICard) {
	game.CurrentMessageCard = currentMessageCard
}

func (game *Game) IsMessageCardFaceUp() bool {
	return game.MessageCardFaceUp
}

func (game *Game) SetMessageCardFaceUp(messageCardFaceUp bool) {
	game.MessageCardFaceUp = messageCardFaceUp
}

func (game *Game) GetLockPlayers() []int {
	return game.WhoIsLocked
}

func (game *Game) IsIdleTimePoint() bool {
	return game.CurrentCard == nil
}

func (game *Game) GetCurrentPhase() protos.Phase {
	return game.CurrentPhase
}

func (game *Game) GetDieState() interfaces.DieState {
	return game.DieState
}

func (game *Game) PlayerDiscardCard(player interfaces.IPlayer, cards ...interfaces.ICard) {
	for _, card := range cards {
		player.DeleteCard(card.GetId())
	}
	logger.Info(player, "弃掉了", cards, "，剩余手牌", len(player.GetCards()), "张")
	game.GetDeck().Discard(cards...)
	for _, p := range game.Players {
		if h, ok := p.(*HumanPlayer); ok {
			msg := &protos.DiscardCardToc{PlayerId: h.GetAlternativeLocation(player.Location())}
			for _, card := range cards {
				msg.Cards = append(msg.Cards, card.ToPbCard())
			}
			h.Send(msg)
		}
	}
}

func (game *Game) checkWinOrDie() bool {
	var killer, stealer interfaces.IPlayer
	var redPlayers, bluePlayers []interfaces.IPlayer
	for _, p := range game.GetPlayers() {
		if p.IsAlive() {
			identity, secretTask := p.GetIdentity()
			switch identity {
			case protos.Color_Black:
				switch secretTask {
				case protos.SecretTask_Killer:
					killer = p
				case protos.SecretTask_Stealer:
					stealer = p
				}
			case protos.Color_Red:
				redPlayers = append(redPlayers, p)
			case protos.Color_Blue:
				bluePlayers = append(bluePlayers, p)
			}
		}
	}
	player := game.GetPlayers()[game.WhoseSendTurn]
	var red, blue, black int
	for _, card := range player.GetMessageCards() {
		for _, color := range card.GetColor() {
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
	var declareWinner interfaces.IPlayer
	var winner []interfaces.IPlayer
	identity, secretTask := player.GetIdentity()
	switch identity {
	case protos.Color_Black:
		if secretTask == protos.SecretTask_Collector && (red >= 3 || blue >= 3) {
			declareWinner = player
			winner = append(winner, player)
		}
	case protos.Color_Red:
		if red >= 3 {
			declareWinner = player
			winner = redPlayers
		}
	case protos.Color_Blue:
		if blue >= 3 {
			declareWinner = player
			winner = bluePlayers
		}
	}
	if declareWinner != nil && stealer != nil && game.WhoseTurn == stealer.Location() {
		declareWinner = stealer
		winner = []interfaces.IPlayer{stealer}
	}
	if declareWinner != nil {
		logger.Info(declareWinner, "宣告胜利，胜利者有", winner)
		for _, p := range game.GetPlayers() {
			p.NotifyWin(declareWinner, winner)
		}
		return true
	}
	if black >= 3 {
		if red+blue < 2 {
			killer = nil
		}
		logger.Info(game.Players[game.WhoDie], "濒死")
		game.WhoDie = game.WhoseSendTurn
		game.DieState = interfaces.DieStateWaitForChengQing
		game.WhoseFightTurn = game.WhoseTurn
		game.AskForChengQing()
		game.afterChengQing = func() {
			if !game.Players[game.WhoDie].IsAlive() && killer != nil && game.WhoseTurn == killer.Location() {
				logger.Info(declareWinner, "宣告胜利，胜利者有", winner)
				for _, p := range game.GetPlayers() {
					p.NotifyWin(killer, []interfaces.IPlayer{killer})
				}
				return
			}
			Post(game.NextTurn)
		}
		return true
	}
	return false
}

func (game *Game) AskForChengQing() {
	if !game.Players[game.WhoseFightTurn].IsAlive() {
		game.AskNextForChengQing()
	}
	logger.Info("正在询问", game.Players[game.WhoseFightTurn], "是否使用澄清")
	for _, p := range game.GetPlayers() {
		p.NotifyAskForChengQing(game.Players[game.WhoDie], game.Players[game.WhoseFightTurn])
	}
}

func (game *Game) AskNextForChengQing() {
	for {
		game.WhoseFightTurn = (game.WhoseFightTurn + 1) % len(game.Players)
		if game.WhoseFightTurn == game.WhoseTurn {
			player := game.Players[game.WhoDie]
			player.SetAlive(false)
			game.GetDeck().Discard(player.DeleteAllCards()...)
			game.GetDeck().Discard(player.DeleteAllMessageCards()...)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(game.WhoDie, false)
			}
			logger.Info("无人拯救，", player, "已死亡")
			game.DieState = interfaces.DieStateDying
			break
		}
		if game.Players[game.WhoseFightTurn].IsAlive() {
			break
		}
	}
	for _, p := range game.GetPlayers() {
		p.NotifyAskForChengQing(game.Players[game.WhoseSendTurn], game.Players[game.WhoseFightTurn])
	}
}

func (game *Game) AfterChengQing() {
	var black int
	for _, card := range game.Players[game.WhoseSendTurn].GetMessageCards() {
		if utils.IsColorIn(protos.Color_Black, card.GetColor()) {
			black++
		}
	}
	if black >= 3 {
		for _, p := range game.GetPlayers() {
			p.NotifyAskForChengQing(game.Players[game.WhoseSendTurn], game.Players[game.WhoseFightTurn])
		}
	} else {
		game.DieState = interfaces.DieStateNone
		if game.afterChengQing != nil {
			Post(game.afterChengQing)
			game.afterChengQing = nil
		}
	}
}
