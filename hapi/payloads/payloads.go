package payloads

type MatchListPayload struct {
	Gamertag string `json:"gamertag"`
	Count    int    `json:"count"`
	Offset   int    `json:"offset"`
	Mode     string `json:"mode"`
}

type MatchList struct {
	Additional struct {
		Gamertag string `json:"gamertag"`
		Mode     string `json:"mode"`
	} `json:"additional"`
	Count int64 `json:"count"`
	Data  []struct {
		Details struct {
			Category struct {
				Asset struct {
					Id           string `json:"id"`
					ThumbnailUrl string `json:"thumbnail_url"`
					Version      string `json:"version"`
				} `json:"asset"`
				Name string `json:"name"`
			} `json:"category"`
			Map struct {
				Asset struct {
					Id           string `json:"id"`
					ThumbnailUrl string `json:"thumbnail_url"`
					Version      string `json:"version"`
				} `json:"asset"`
				Name string `json:"name"`
			} `json:"map"`
			Playlist struct {
				Asset struct {
					Id           string `json:"id"`
					ThumbnailUrl string `json:"thumbnail_url"`
					Version      string `json:"version"`
				} `json:"asset"`
				Name       string `json:"name"`
				Properties struct {
					Input  string `json:"input"`
					Queue  string `json:"queue"`
					Ranked bool   `json:"ranked"`
				} `json:"properties"`
			} `json:"playlist"`
		} `json:"details"`
		Duration struct {
			Human   string `json:"human"`
			Seconds int64  `json:"seconds"`
		} `json:"duration"`
		Experience string `json:"experience"`
		Id         string `json:"id"`
		PlayedAt   string `json:"played_at"`
		Player     struct {
			Outcome string `json:"outcome"`
			Rank    int64  `json:"rank"`
			Stats   struct {
				Core struct {
					Breakdowns struct {
						Assists struct {
							Callouts int64 `json:"callouts"`
							Driver   int64 `json:"driver"`
							Emp      int64 `json:"emp"`
						} `json:"assists"`
						Kills struct {
							Grenades     int64 `json:"grenades"`
							Headshots    int64 `json:"headshots"`
							Melee        int64 `json:"melee"`
							PowerWeapons int64 `json:"power_weapons"`
						} `json:"kills"`
					} `json:"breakdowns"`
					Damage struct {
						Dealt int64 `json:"dealt"`
						Taken int64 `json:"taken"`
					} `json:"damage"`
					Kda    int64   `json:"kda"`
					Kdr    float64 `json:"kdr"`
					Rounds struct {
						Lost int64 `json:"lost"`
						Tied int64 `json:"tied"`
						Won  int64 `json:"won"`
					} `json:"rounds"`
					Score int64 `json:"score"`
					Shots struct {
						Accuracy float64 `json:"accuracy"`
						Fired    int64   `json:"fired"`
						Landed   int64   `json:"landed"`
						Missed   int64   `json:"missed"`
					} `json:"shots"`
					Summary struct {
						Assists   int64 `json:"assists"`
						Betrayals int64 `json:"betrayals"`
						Deaths    int64 `json:"deaths"`
						Kills     int64 `json:"kills"`
						Medals    int64 `json:"medals"`
						Suicides  int64 `json:"suicides"`
						Vehicles  struct {
							Destroys int64 `json:"destroys"`
							Hijacks  int64 `json:"hijacks"`
						} `json:"vehicles"`
					} `json:"summary"`
				} `json:"core"`
				Mode interface{} `json:"mode"`
			} `json:"stats"`
			Team struct {
				EmblemUrl string `json:"emblem_url"`
				Id        int64  `json:"id"`
				Name      string `json:"name"`
			} `json:"team"`
		} `json:"player"`
		Teams struct {
			Enabled bool `json:"enabled"`
			Scoring bool `json:"scoring"`
		} `json:"teams"`
	} `json:"data"`
	Paging struct {
		Count  int64 `json:"count"`
		Offset int64 `json:"offset"`
	} `json:"paging"`
}
