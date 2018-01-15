package standings

type Stats struct {
	Wins   int
	Losses int
	Draws  int
}

type Standing struct {
	Place            int
	Name             string
	Points           int
	GamesPlayed      int
	GoalsFor         int
	GoalsAgainst     int
	GoalDifferential int
	OverallStats     Stats
	HomeStats        Stats
	AwayStats        Stats
	Clinched         bool
	ConferenceWinner bool
	ShieldWinner     bool
}

type Standings []Standing
