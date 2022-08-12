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
	Players          []interfaces.IPlayer
	TotalPlayerCount int
	Deck             interfaces.IDeck
	CurrentCard      interfaces.ICard
	WhoseTurn        int
	CurrentPhase     protos.Phase
	Random           *rand.Rand
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
	logger.Infof("已加入%d个机器人，等待%d人加入。。。", robotCount, humanCount)

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
		case *protos.UseShiTanTos:
			humanMap[ev.Session().ID()].onUseShiTan(pb)
		case *protos.ExecuteShiTanTos:
			humanMap[ev.Session().ID()].onExecuteShiTan(pb)
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

func (game *Game) SendPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		game.Post(game.NextTurn)
		return
	}
	game.Post(game.FightPhase)
}

func (game *Game) FightPhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		game.Post(game.NextTurn)
		return
	}
	game.Post(game.ReceivePhase)
}

func (game *Game) ReceivePhase() {
	player := game.Players[game.WhoseTurn]
	if !player.IsAlive() {
		game.Post(game.NextTurn)
		return
	}
	game.Post(game.NextTurn)
}

func (game *Game) NextTurn() {
	for {
		game.WhoseTurn = (game.WhoseTurn + 1) % len(game.Players)
		player := game.Players[game.WhoseTurn]
		if player.IsAlive() {
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

func (game *Game) GetCurrentCard() interfaces.ICard {
	return game.CurrentCard
}

func (game *Game) SetCurrentCard(card interfaces.ICard) {
	game.CurrentCard = card
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
	game.Deck.Discard(cards...)
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
