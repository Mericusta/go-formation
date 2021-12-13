package formation

type GameDataJsonObject struct {
	Format map[string]int  `json:"Format"`
	Data   [][]interface{} `json:"Data"`
}
