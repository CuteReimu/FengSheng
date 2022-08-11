package game

import (
	"github.com/CuteReimu/FengSheng/config"
	_ "github.com/CuteReimu/FengSheng/core"
	"github.com/CuteReimu/FengSheng/game/interfaces"
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
	Random           *rand.Rand
	cellnet.EventQueue
}

func (game *Game) Start(totalCount, robotCount int) {
	game.Random = rand.New(rand.NewSource(time.Now().Unix()))
	humanCount := totalCount - robotCount
	game.TotalPlayerCount = totalCount
	index := 0
	for ; index < robotCount; index++ {
		game.Players = append(game.Players, new(interfaces.BasePlayer))
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
		switch ev.Message().(type) {
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
	for location, player := range game.Players {
		player.Init(game, location)
	}
	game.Deck = NewDeck(game)
	game.WhoseTurn = game.Random.Intn(len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		game.Players[(game.WhoseTurn+i)%len(game.Players)].Draw(config.GetHandCardCountBegin())
	}
	game.WhoseTurn = len(game.Players) - 1
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
