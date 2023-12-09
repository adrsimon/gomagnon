package typing

import (
	"slices"
	"time"
)

var Needs = map[string]map[ResourceType]int{
	"hut": map[ResourceType]int{
		WOOD: 3,
		ROCK: 3,
	},
}

type Hut struct {
	Position  *Hexagone
	Inventory []ResourceType
	Owner     *Human
	Ballot    Ballot
}

type Ballot struct {
	reason         string
	VoteInProgress bool
	EndTimeVote    time.Time
	VotersID       []string
	Profile        []bool
}

func (hu *Hut) removeVoter(agID string) {
	for i, v := range hu.Ballot.VotersID {
		if v == agID {
			hu.Ballot.VotersID = append(hu.Ballot.VotersID[:i], hu.Ballot.VotersID[i+1:]...)
		}
	}
}
func (hu *Hut) StartNewVote(agent *Human, reason string) bool {
	if hu.Ballot.VoteInProgress == false && agent == agent.Clan.chief {
		switch reason {
		case "VoteNewPerson":
			hu.Ballot.VoteInProgress = true
			hu.Ballot.EndTimeVote = time.Now().Add(15 * time.Second)
			for _, v := range agent.Clan.members {
				hu.Ballot.VotersID = append(hu.Ballot.VotersID, v.ID)
			}
			hu.Ballot.VotersID = append(hu.Ballot.VotersID, agent.ID)
		}
		return true
	}
	return false
}

func (hu *Hut) Vote(agent *Human, choice string) bool {
	if !hu.Ballot.EndTimeVote.Before(time.Now()) {
		if slices.Contains(hu.Ballot.VotersID, agent.ID) {
			hu.removeVoter(agent.ID)
			if choice == "VoteYes" {
				hu.Ballot.Profile = append(hu.Ballot.Profile, true)
			} else {
				hu.Ballot.Profile = append(hu.Ballot.Profile, false)
			}
			return true

		}
	}
	return false
}

func (hu *Hut) CountVotes() bool {
	var yesCount, noCount int
	for _, vote := range hu.Ballot.Profile {
		if vote {
			yesCount++
		} else {
			noCount++
		}
	}
	return yesCount >= noCount
}

func (hu *Hut) GetResult(agent *Human) bool {
	if hu.Ballot.EndTimeVote.Before(time.Now()) && agent == agent.Clan.chief {
		hu.Ballot.reason = ""
		hu.Ballot.VotersID = make([]string, 0)
		hu.Ballot.EndTimeVote = time.Time{}
		hu.Ballot.VoteInProgress = false
		return hu.CountVotes()
	}
	return false
}
