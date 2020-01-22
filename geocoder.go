package my1562geocoder

import (
	"fmt"
	"math"
	"sort"
	"strconv"
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

func getDistance(lat0, lng0, lat1, lng1 float64) float64 {
	dLatMeters := (lat0 - lat1) * meridianDegreeLength
	dLngMeters := (lng0 - lng1) * parallelDegreeLength(lat0)
	return math.Sqrt(dLatMeters*dLatMeters + dLngMeters*dLngMeters)
}

type Rectangle struct {
	iLat int32
	iLng int32
}

func (rect *Rectangle) ToString() string {
	return strconv.Itoa(int(rect.iLat)) + "." + strconv.Itoa(int(rect.iLng))
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

func ReverseGeocode(lat float64, lng float64, accuracyMeters float64, limit int) {
	rectangles := getNearbyRectangles(lat, lng, GeoIndexResolution, accuracyMeters)
	addressIDs := make([]uint32, 0)
	for _, rect := range rectangles {
		key := rect.ToString()
		addressIDsInRectangle, ok := GeoIndexData[key]
		if !ok {
			continue
		}
		addressIDs = append(addressIDs, addressIDsInRectangle...)
	}

	type addressWithDistance struct {
		address  *FullAddress
		distance float64
	}

	addresses := make([]*addressWithDistance, 0)
	for _, id := range addressIDs {
		addr, ok := Addresses[id]
		if !ok {
			continue
		}
		distance := getDistance(lat, lng, addr.Lat, addr.Lng)
		if distance <= accuracyMeters {
			fullAddress := resolveAddress(&addr)
			if fullAddress.street1562 != nil {
				addresses = append(addresses, &addressWithDistance{
					fullAddress,
					distance,
				})
			}
		}
	}
	sort.Slice(addresses, func(i, j int) bool {
		return addresses[i].distance < addresses[j].distance
	})
	safeLimit := limit
	if len(addresses) < limit {
		safeLimit = len(addresses)
	}
	addresses = addresses[:safeLimit]
	fmt.Println(addresses)
	// return addresses
}
