# MLS Data Gatherer

## Routes

- `/`: JSON of current season's fixtures
- `/standings/shield`: JSON of current season's Supporters Shield Standings
- `/standings/conference/(east|west)`: JSON of current season's Conference Standings
- `/reddit/(club abbreviation)/automod`: AutoMod-formatted list of current season's fixtures to schedule Reddit Automod posts. Accepts `offset` query string param for hours before kickoff: `?offset=1` for one hour before kickoff.
- `/reddit/(club abbreviation)/schedule`: Reddit markdown-formatted table of previous and upcoming matches. Takes `prevCount`, `nextCount`, and `showForm` params, all of which have numerical values
- `/reddit/(club abbreviation)/standings`: Reddit markdown-formatted standings for specified team's conference
- More to come
