package server

import (
	"blazecache/cache"
	"net/http"

	"github.com/labstack/echo/v4"
)

type setRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (s *Server) setHandler(c echo.Context) error {

	c.Logger().Info("Request handle by set method")
	s.logChan <- "Server on port " + s.port + " : handle set method"

	var setRequest setRequest
	err := c.Bind(&setRequest)
	if err != nil {
		return err
	}
	s.cache.Set(setRequest.Key, setRequest.Value, cache.DefaultExpiration)
	if s.state == Leader {
		s.BroadcastDatas(setRequest)
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) getHandler(c echo.Context) error {

	c.Logger().Info("Request handle by get method")
	s.logChan <- "Server on port " + s.port + " : handle get method"

	key := c.QueryParam("key")

	r, found := s.cache.Get(key)
	if !found {
		return c.JSON(http.StatusNotFound, "item does not exists")
	}

	return c.JSON(http.StatusOK, r)
}

type heartbeatResponse struct {
	Term    int  `json:"term"`
	Success bool `json:"success"`
}

func (s *Server) heartbeatHandler(c echo.Context) error {
	var heartbeatRequest heartbeatRequest
	var heartbeatResponse heartbeatResponse

	err := c.Bind(&heartbeatRequest)
	if err != nil {
		return err
	}

	if s.currentTerm > heartbeatRequest.Term {
		heartbeatResponse.Success = false
		heartbeatResponse.Term = s.currentTerm
		return c.JSON(http.StatusOK, heartbeatResponse)
	}

	s.isHeartbeatByLeader <- true
	s.leaderId = heartbeatRequest.Leader
	heartbeatResponse.Success = true
	heartbeatResponse.Term = s.currentTerm
	return c.JSON(http.StatusOK, heartbeatResponse)
}

type voteResponse struct {
	Term          int  `json:"term"`
	IsVoteGranted bool `json:"isVoteGranted"`
}

func (s *Server) voteHandler(c echo.Context) error {
	var voteRequest voteRequest
	var voteResponse voteResponse
	err := c.Bind(&voteRequest)
	if err != nil {
		return err
	}

	if voteRequest.Term < s.currentTerm {
		voteResponse.Term = s.currentTerm
		voteResponse.IsVoteGranted = false
	}

	if s.voteFor == -1 {
		s.currentTerm = voteRequest.Term
		s.voteFor = voteRequest.Candidate
		voteResponse.Term = s.currentTerm
		voteResponse.IsVoteGranted = true
	}

	return c.JSON(http.StatusOK, voteResponse)
}
