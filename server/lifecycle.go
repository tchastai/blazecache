package server

import (
	"math/rand"
	"strconv"
	"time"
)

func (s *Server) follower() {
	select {
	case <-s.isHeartbeatByLeader:
		s.logChan <- "Server number " + strconv.Itoa(s.me) + " have received a Heartbeat from the leader number : " + strconv.Itoa(s.leaderId)
	case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
		s.logChan <- "Follower number : " + strconv.Itoa(s.me) + " timeout"
		s.state = Candidate
	}
}

func (s *Server) candidate() {
	s.currentTerm++
	s.voteFor = s.me
	s.voteCount = 1

	go s.broadcastVote()

	select {
	case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
		s.state = Follower
	case <-s.isBecomeLeader:
		s.state = Leader
	}
}

func (s *Server) leader() {
	s.broadcastHearbeat()
	time.Sleep(100 * time.Millisecond)
}
