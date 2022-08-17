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
	"math/rand"
	"os"
	"time"
)

var logger = utils.GetLogger("game")

type Game struct {
	Players            []interfaces.IPlayer
	TotalPlayerCount   int
	Deck               interfaces.IDeck
	CurrentCard        *interfaces.CurrentCard
	WhoseTurn          int
	WhoseSendTurn      int
	WhoseFightTurn     int
	CurrentMessageCard interfaces.ICard
	MessageCardFaceUp  bool
	CardDirection      protos.Direction
	WhoIsLocked        []int
	CurrentPhase       protos.Phase
	Random             *rand.Rand
	cellnet.EventQueue
}

func (game *Game) Start(totalCount, robotCount int) {
	game.Random = rand.New(rand.NewSource(time.Now().Unix()))
	humanCount := totalCount - robotCount
	game.TotalPlayerCount = totalCount
	index := 0
	for ; index < robotCount; index++ {
		game.Players = append(game.Players, new(RobotPlayer))
	}
	logger.Info("已加入", robotCount, "个机器人，等待", humanCount, "人加入。。。")

	if !config.IsTcpDebugLogOpen() {
		msglog.SetCurrMsgLogMode(msglog.MsgLogMode_Mute)
	}
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	game.EventQueue = cellnet.NewEventQueue()

	// 创建一个tcp的侦听器，名称为server，所有连接将事件投递到queue队列,单线程的处理
	p := peer.NewGenericPeer("tcp.Acceptor", "server", config.GetListenAddress(), game.EventQueue)

	humanMap := make(map[int64]*HumanPlayer)
	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {
		switch pb := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			if index < totalCount {
				player := &HumanPlayer{Session: ev.Session()}
				game.Players = append(game.Players, player)
				humanMap[player.Session.ID()] = player
				index++
				logger.Info("server accepted", player.Session)
				if index == totalCount {
					game.Post(game.start)
				}
			} else {
				ev.Session().Close()
				logger.Info("房间人数已满")
			}
		case *cellnet.SessionClosed:
			logger.Info("session closed: ", ev.Session().ID())
			if _, ok := humanMap[ev.Session().ID()]; ok {
				logger.Info("目前不支持断线重连，程序将在3秒后关闭")
				time.Sleep(time.Second * 3)
				os.Exit(1)
			}
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
		}
	})
	p.Start()
	game.StartLoop()
	if humanCount == 0 {
		game.Post(game.start)
	}
	game.Wait()
}

func (game *Game) start() {
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
	game.Random.Shuffle(3, func(i, j int) {
		idTask[len(idTask)-3+i], idTask[len(idTask)-3+j] = idTask[len(idTask)-3+j], idTask[len(idTask)-3+i]
	})
	game.Random.Shuffle(len(game.Players), func(i, j int) {
		idTask[i], idTask[j] = idTask[j], idTask[i]
	})
	for location, player := range game.Players {
		player.Init(game, location, idTask[location].id, idTask[location].task)
	}
	game.Deck = NewDeck(game)
	game.WhoseTurn = game.Random.Intn(len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		game.Players[(game.WhoseTurn+i)%len(game.Players)].Draw(config.GetHandCardCountBegin())
	}
	game.DrawPhase()
}

func (game *Game) DrawPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		game.Post(game.NextTurn)
		return
	}
	logger.Info(player, "的回合开始了")
	game.CurrentPhase = protos.Phase_Draw_Phase
	for _, p := range game.Players {
		p.NotifyDrawPhase()
	}
	player.Draw(3)
	game.Post(game.MainPhase)
}

func (game *Game) MainPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		game.Post(game.NextTurn)
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
			game.GetDeck().Discard(player.DeleteAllCards()...)
			player.SetLose(true)
			player.SetAlive(false)
			for _, p := range game.GetPlayers() {
				p.NotifyDie(game.WhoseTurn, true)
			}
		}
	}
	if !player.IsAlive() {
		game.Post(game.NextTurn)
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
		p.NotifySendPhase(20)
	}
}

func (game *Game) MessageMoveNext() {
	if game.CardDirection == protos.Direction_Up {
		if game.Players[game.WhoseTurn].IsAlive() {
			game.WhoseSendTurn = game.WhoseTurn
		} else {
			game.Post(game.NextTurn)
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
				game.Post(game.NextTurn)
				return
			}
		}
	}
	logger.Info("情报到达", game.Players[game.WhoseSendTurn], "面前")
	for _, p := range game.Players {
		p.NotifySendPhase(20)
	}
}

func (game *Game) OnChooseReceiveCard() {
	logger.Info(game.Players[game.WhoseSendTurn], "选择接收情报")
	game.WhoseFightTurn = game.WhoseSendTurn
	game.CurrentPhase = protos.Phase_Fight_Phase
	for _, p := range game.Players {
		p.NotifyFightPhase(5)
	}
}

func (game *Game) FightPhaseNext() {
	for {
		game.WhoseFightTurn = (game.WhoseFightTurn + 1) % len(game.Players)
		if game.WhoseFightTurn == game.WhoseSendTurn {
			game.Post(game.ReceivePhase)
			return
		} else if game.Players[game.WhoseFightTurn].IsAlive() {
			break
		}
	}
	for _, p := range game.Players {
		p.NotifyFightPhase(5)
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
	}
	game.Post(game.NextTurn)
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
	game.Post(game.DrawPhase)
}

func (game *Game) GetPlayers() []interfaces.IPlayer {
	return game.Players
}

func (game *Game) GetDeck() interfaces.IDeck {
	return game.Deck
}

func (game *Game) GetWhoseTurn() int {
	return game.WhoseTurn
}

func (game *Game) GetWhoseSendTurn() int {
	return game.WhoseSendTurn
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

func (game *Game) GetRandom() *rand.Rand {
	return game.Random
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
