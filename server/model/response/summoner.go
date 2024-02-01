package response

type SummonerDTO struct {
	Name          string      `json:"name"`
	Loc           string      `json:"loc"`
	ProfileIconID int         `json:"profileIconID"`
	SummonerLevel int         `json:"summonerLevel"`
	SoloEntry     *EntryDTO   `json:"soloEntry"`
	FlexEntry     *EntryDTO   `json:"flexEntry"`
	Matches       []*MatchDTO `json:"matches"`
}
