package my1562geocoder

import (
	"math"
)

func GetAddressByID(id uint32) *Address {
	addr, ok := Addresses[id]
	if ok {
		return &addr
	}
	return nil
}

const meridianDegreeLength = 111000.0

func parallelDegreeLength(lat float64) float64 {
	return math.Cos(lat*math.Pi/180.0) * meridianDegreeLength
}

func transpose(lat float64, lng float64, dLatMeters float64, dLngMeters float64) (float64, float64) {
	dLatDegrees := dLatMeters / meridianDegreeLength
	dLngDegrees := dLngMeters / parallelDegreeLength(lat)
	return lat + dLatDegrees, lng + dLngDegrees
}

func distance(lat0, lng0, lat1, lng1 float64) float64 {
	dLatMeters := (lat0 - lat1) * meridianDegreeLength
	dLngMeters := (lng0 - lng1) * parallelDegreeLength(lat0)
	return math.Sqrt(dLatMeters*dLatMeters + dLngMeters*dLngMeters)
}

type Rectangle struct {
	iLat int32
	iLng int32
}

func getRectangleByLatLng(lat float64, lng float64, res int32) Rectangle {
	var iLat, iLng int32
	iLat = int32((lat + 90.0) * GeoIndexResolution)
	iLng = int32((lng + 180.0) * GeoIndexResolution)
	return Rectangle{iLat, iLng}
}

func getNearbyRectangles(lat float64, lng float64, res int32, accuracyMeters float64) []Rectangle {
	minLat, minLng := transpose(lat, lng, -accuracyMeters, -accuracyMeters)
	maxLat, maxLng := transpose(lat, lng, accuracyMeters, accuracyMeters)

	minRect := getRectangleByLatLng(minLat, minLng, res)
	maxRect := getRectangleByLatLng(maxLat, maxLng, res)

	length := (maxRect.iLat - minRect.iLat + 1) * (maxRect.iLng - minRect.iLng + 1)
	slice := make([]Rectangle, length)
	index := 0
	for iLat := minRect.iLat; iLat <= maxRect.iLat; iLat++ {
		for iLng := minRect.iLng; iLng <= maxRect.iLng; iLng++ {
			slice[index] = Rectangle{iLat, iLng}
			index++
		}
	}

	return slice
}
