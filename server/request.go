package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type voteRequest struct {
	Term      int `json:"term"`
	Candidate int `json:"candidate"`
}

func (s *Server) broadcastVote() error {
	v := voteRequest{
		Term:      s.currentTerm,
		Candidate: s.voteFor,
	}

	vJson, err := json.Marshal(v)
	if err != nil {
		return err
	}

	for _, v := range s.agentAdress {

		s.logChan <- "Server number " + strconv.Itoa(s.me) + ": Send vote request to : " + v
		resp, err := http.Post("http://127.0.0.1"+v+"/consensus/vote", "application/json", bytes.NewBuffer(vJson))
		if err != nil {
			s.logChan <- "Server number " + strconv.Itoa(s.me) + "have an error while trying to send vote request to server with the adress localhost" + v + " with err = " + err.Error()
		}

		s.logChan <- "Server number " + strconv.Itoa(s.me) + ": deserialize response from : " + v
		var voteResponse voteResponse
		if err := json.NewDecoder(resp.Body).Decode(&voteResponse); err != nil {
			s.logChan <- "Cannot deserialize vote response for vote request"
		}

		if voteResponse.Term > s.currentTerm {
			s.currentTerm = voteResponse.Term
			s.state = Follower
			s.voteFor = -1
			return nil
		}

		if voteResponse.IsVoteGranted {
			s.voteCount++
		}

		if s.voteCount >= len(s.agentAdress)-1 {
			s.logChan <- v + " have vote for " + strconv.Itoa(s.me) + " he becomes a leader"

			s.isBecomeLeader <- true
		}

		resp.Body.Close()

	}

	return nil
}

type heartbeatRequest struct {
	Term   int
	Leader int
}

func (s *Server) broadcastHearbeat() error {
	for _, v := range s.agentAdress {

		g := heartbeatRequest{
			Term:   s.currentTerm,
			Leader: s.me,
		}

		gJson, err := json.Marshal(g)
		if err != nil {
			return err
		}

		resp, err := http.Post("http://127.0.0.1"+v+"/consensus/heartbeat", "application/json", bytes.NewBuffer(gJson))
		if err != nil {
			s.logChan <- "Server number " + strconv.Itoa(s.me) + " have an error while trying to send heartbeat request to server with the adress localhost" + v + " with err = " + err.Error()
		}

		var heartbeatResponse heartbeatResponse
		if err := json.NewDecoder(resp.Body).Decode(&heartbeatResponse); err != nil {
			s.logChan <- "Cannot deserialize vote response for heartbeat"
		}

		if !heartbeatResponse.Success && heartbeatResponse.Term > s.currentTerm {
			s.currentTerm = heartbeatResponse.Term
			s.state = Follower
			s.voteFor = -1
			return nil
		}

		resp.Body.Close()
	}

	return nil
}

func (s *Server) BroadcastDatas(setRequest setRequest) error {
	for _, v := range s.agentAdress {

		setRequestJson, err := json.Marshal(setRequest)
		if err != nil {
			return err
		}

		resp, err := http.Post("http://127.0.0.1"+v+"/cache", "application/json", bytes.NewBuffer(setRequestJson))
		if err != nil {
			s.logChan <- "Server number " + strconv.Itoa(s.me) + " have an error while trying to send broadcast to server with the adress localhost" + v + " with err = " + err.Error()
		}

		if resp.StatusCode != 200 {
			s.logChan <- "Server number " + strconv.Itoa(s.me) + " have a bad status while broadcasting datas status = " + strconv.Itoa(resp.StatusCode)
		}

		resp.Body.Close()
	}

	return nil
}
