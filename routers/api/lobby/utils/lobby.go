package utils

import "github.com/nowel-xyz/quiz/database/models"

type Settings struct {
	MemberLimit int `json:"member_limit"`
}

type LobbyInvite struct {
	Code string `json:"code"`
}


type Lobby struct {
	ID        string      `json:"lobby_id"`
	Invite    LobbyInvite `json:"lobby_invite"`
	HostID    string      `json:"host_id"`
	QuizID    string      `json:"quiz_id"`
	Members []models.LobbyUser `json:"members"`
	Settings  Settings    `json:"settings"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Started   bool        `json:"started"`

}