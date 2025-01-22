package server

import (
	"blazecache/cache"
	"blazecache/util"
	"time"

	"github.com/labstack/echo/v4"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Server struct {
	me                  int
	cache               *cache.Cache
	port                string
	logChan             chan string
	state               State
	currentTerm         int
	agentAdress         []string
	voteFor             int
	voteCount           int
	isHeartbeatByLeader chan bool
	isBecomeLeader      chan bool
	leaderId            int
}

func New(me int, port string, logChan chan string, agentAdress []string) *Server {

	c := cache.New(5*time.Minute, 10*time.Minute)

	return &Server{
		me:                  me,
		cache:               c,
		port:                port,
		logChan:             logChan,
		voteFor:             -1,
		state:               Follower,
		currentTerm:         0,
		agentAdress:         util.RemoveStringFromList(agentAdress, port),
		isHeartbeatByLeader: make(chan bool),
		isBecomeLeader:      make(chan bool),
	}
}

func (s *Server) Start() {

	e := echo.New()

	cacheGroup := e.Group("/cache")
	{
		cacheGroup.GET("", s.getHandler)
		cacheGroup.POST("", s.setHandler)
	}

	consensusGroup := e.Group("/consensus")
	{
		consensusGroup.POST("/heartbeat", s.heartbeatHandler)
		consensusGroup.POST("/vote", s.voteHandler)
	}

	go s.startConsensus()

	s.logChan <- "Starting server on port: " + s.port
	e.Logger.Fatal(e.Start(s.port))
}

func (s *Server) startConsensus() {
	for {
		switch s.state {
		case Follower:
			s.follower()
		case Candidate:
			s.candidate()
		case Leader:
			s.leader()
		}
	}
}
