module go-halo.com/fetchers

go 1.17

replace go-halo.com/fetchers/max => ./max

replace go-halo.com/hapi/payloads => ../hapi/payloads

replace go-halo.com/hapi/stats/matches => ../hapi/stats/matches

require go-halo.com/fetchers/max v0.0.0-00010101000000-000000000000

require (
	go-halo.com/hapi/payloads v0.0.0-00010101000000-000000000000 // indirect
	go-halo.com/hapi/stats/matches v0.0.0-00010101000000-000000000000 // indirect
)
