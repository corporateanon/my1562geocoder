package geocoder

import "math"

const meridianDegreeLength = 111000.0

func parallelDegreeLength(lat float64) float64 {
	return math.Cos(lat*math.Pi/180.0) * meridianDegreeLength
}

func transpose(lat float64, lng float64, dLatMeters float64, dLngMeters float64) (float64, float64) {
	dLatDegrees := dLatMeters / meridianDegreeLength
	dLngDegrees := dLngMeters / parallelDegreeLength(lat)
	return lat + dLatDegrees, lng + dLngDegrees
}

func getDistance(lat0, lng0, lat1, lng1 float64) float64 {
	dLatMeters := (lat0 - lat1) * meridianDegreeLength
	dLngMeters := (lng0 - lng1) * parallelDegreeLength(lat0)
	return math.Sqrt(dLatMeters*dLatMeters + dLngMeters*dLngMeters)
}
