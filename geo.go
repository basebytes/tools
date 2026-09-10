package tools

import "math"

func GeoDistance(lng1, lat1, lng2, lat2 float64) float64 {
	const radius = 6371.0
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLng := (lng2 - lng1) * (math.Pi / 180)
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Cos(lat1*(math.Pi/180))*
			math.Cos(lat2*(math.Pi/180))*
			math.Pow(math.Sin(dLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return radius * c * 1000
}
