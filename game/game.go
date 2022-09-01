package game

import (
	"github.com/CuteReimu/FengSheng/config"
	_ "github.com/CuteReimu/FengSheng/core"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"google.golang.org/protobuf/proto"
	"time"
)

var logger = utils.GetLogger("game")

var eventQueue cellnet.EventQueue

func Post(callback func()) {
	eventQueue.Post(callback)
}

type Fsm interface {
	Resolve() (next Fsm, continueResolve bool)
}

type WaitingFsm interface {
	Fsm
	ResolveProtocol(player IPlayer, pb proto.Message) (next Fsm, continueResolve bool)
}

type Game struct {
	Id      uint32
	Players []IPlayer
	Deck    *Deck
	fsm     Fsm
}

func (game *Game) end() {
	delete(Cache, game.Id)
}

var Cache = make(map[uint32]*Game)

func Start(totalCount int) {
	gameId := uint32(1)
	game := &Game{Id: gameId, Players: make([]IPlayer, totalCount)}

	if !config.IsTcpDebugLogOpen() {
		msglog.SetCurrMsgLogMode(msglog.MsgLogMode_Mute)
	}
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	eventQueue = cellnet.NewEventQueue()

	// 创建一个tcp的侦听器，名称为server，所有连接将事件投递到queue队列,单线程的处理
	p := peer.NewGenericPeer("tcp.Acceptor", "server", config.GetListenAddress(), eventQueue)

	humanMap := make(map[int64]*HumanPlayer)

	onAdd := func(player IPlayer) int {
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
			if player, ok := game.Players[i].(*HumanPlayer); ok && i != playerIndex {
				player.Send(msg)
			}
		}
		if unready == 0 {
			logger.Info(player, "加入了。已加入", totalCount, "个人，游戏开始。。。")
			Post(game.start)
			gameId++
			game = &Game{Id: gameId, Players: make([]IPlayer, totalCount)}
		} else {
			logger.Info(player, "加入了。已加入", totalCount-unready, "个人，等待", unready, "人加入。。。")
		}
		return playerIndex
	}

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {
		switch pb := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			logger.Info("session connected: ", ev.Session().ID())
		case *protos.JoinRoomTos:
			if pb.Version >= config.GetClientVersion() {
				if len(Cache) <= config.GetMaxRoomCount() {
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
				} else {
					ev.Session().Send(&protos.ErrorCodeToc{
						Code: protos.ErrorCode_no_more_room,
					})
				}
			} else {
				ev.Session().Send(&protos.ErrorCodeToc{
					Code:      protos.ErrorCode_client_version_not_match,
					IntParams: []int64{int64(config.GetClientVersion())},
				})
			}
		case *cellnet.SessionClosed:
			logger.Info("session closed: ", ev.Session().ID())
			if player, ok := humanMap[ev.Session().ID()]; ok {
				if player.GetGame() != nil {
					game := player.GetGame()
					if func(players []IPlayer) bool {
						for i := range players {
							if _, ok := players[i].(*HumanPlayer); ok {
								return true
							}
						}
						return false
					}(game.GetPlayers()) {
						game.Players[player.Location()] = &RobotPlayer{BasePlayer: player.BasePlayer}
					} else {
						for i := range game.GetPlayers() {
							switch p := game.Players[i].(type) {
							case *HumanPlayer:
								game.Players[i] = &IdlePlayer{BasePlayer: p.BasePlayer}
							case *RobotPlayer:
								game.Players[i] = &IdlePlayer{BasePlayer: p.BasePlayer}
							}
						}
						game.end()
					}
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
	utils.Random.Shuffle(len(RoleCache), func(i, j int) {
		RoleCache[i], RoleCache[j] = RoleCache[j], RoleCache[i]
	})
	for location, player := range game.Players {
		if location < len(RoleCache) {
			player.Init(game, location, idTask[location].id, idTask[location].task, RoleCache[location])
		} else {
			player.Init(game, location, idTask[location].id, idTask[location].task, nil)
		}
	}
	Cache[game.Id] = game
	game.Deck = NewDeck(game)
	whoseTurn := utils.Random.Intn(len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		game.Players[(whoseTurn+i)%len(game.Players)].Draw(config.GetHandCardCountBegin())
	}
	time.AfterFunc(time.Second, func() {
		game.Resolve(&DrawPhase{Player: game.Players[whoseTurn]})
	})
}

func (game *Game) GetPlayers() []IPlayer {
	return game.Players
}

func (game *Game) GetDeck() *Deck {
	return game.Deck
}

func (game *Game) PlayerDiscardCard(player IPlayer, cards ...ICard) {
	if len(cards) == 0 {
		return
	}
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

func (game *Game) GetFsm() Fsm {
	return game.fsm
}

func (game *Game) ContinueResolve() {
	Post(func() {
		var continueResolve bool
		game.fsm, continueResolve = game.fsm.Resolve()
		if continueResolve {
			game.ContinueResolve()
		}
	})
}

func (game *Game) TryContinueResolveProtocol(player *HumanPlayer, pb proto.Message) {
	fsm, ok := game.GetFsm().(WaitingFsm)
	if !ok {
		logger.Error("时机错误", fsm)
		return
	}
	Post(func() {
		game.fsm = fsm
		var continueResolve bool
		game.fsm, continueResolve = fsm.ResolveProtocol(player, pb)
		if continueResolve {
			game.ContinueResolve()
		}
	})
}

func (game *Game) Resolve(fsm Fsm) {
	game.fsm = fsm
	game.ContinueResolve()
}
