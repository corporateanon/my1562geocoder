package my1562geocoder

import (
	"encoding/gob"
	"os"
	"sort"
)

type Geocoder struct {
	data  *GeocoderData
	index *GeoIndexWithResolution
}

func NewGeocoder(gobFile string) *Geocoder {
	file, err := os.Open(gobFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gob.RegisterName("*GeocoderData", &GeocoderData{})
	dec := gob.NewDecoder(file)

	data := GeocoderData{}
	err = dec.Decode(&data)
	if err != nil {
		panic(err)
	}

	return &Geocoder{
		data: &data,
	}
}

func (geo *Geocoder) BuildSpatialIndex(resolution float64) {
	index := &GeoIndexWithResolution{Resolution: resolution, Items: GeoIndex{}}
	for _, address := range geo.data.Addresses {
		rect := NewRectangle(address.Lat, address.Lng, resolution)
		key := rect.ToString()
		if index.Items[key] == nil {
			index.Items[key] = []uint32{}
		}
		index.Items[key] = append(index.Items[key], address.ID)
	}
	geo.index = index
}

func (geo *Geocoder) getNearbyRectangles(
	lat float64,
	lng float64,
	accuracyMeters float64,
) []Rectangle {
	minLat, minLng := transpose(lat, lng, -accuracyMeters, -accuracyMeters)
	maxLat, maxLng := transpose(lat, lng, accuracyMeters, accuracyMeters)

	minRect := NewRectangle(minLat, minLng, geo.index.Resolution)
	maxRect := NewRectangle(maxLat, maxLng, geo.index.Resolution)

	length := (maxRect.ILat - minRect.ILat + 1) * (maxRect.ILng - minRect.ILng + 1)
	slice := make([]Rectangle, length)
	index := 0
	for iLat := minRect.ILat; iLat <= maxRect.ILat; iLat++ {
		for iLng := minRect.ILng; iLng <= maxRect.ILng; iLng++ {
			slice[index] = Rectangle{iLat, iLng}
			index++
		}
	}

	return slice
}

func (geo *Geocoder) resolveAddress(address *Address) *FullAddress {
	fa := &FullAddress{}
	fa.Address = address
	streetID := address.StreetID

	streetAR, streetARisResolved := geo.data.StreetsAR[streetID]
	if streetARisResolved {
		fa.StreetAR = streetAR
	}

	street1562Id, street1562IdIsResolved := geo.data.MappingArTo1562[streetID]
	if street1562IdIsResolved {
		street1562, street1562IsResolved := geo.data.Streets1562[street1562Id]
		if street1562IsResolved {
			fa.Street1562 = street1562
		}
	}
	return fa
}

func (geo *Geocoder) ReverseGeocode(
	lat float64,
	lng float64,
	accuracyMeters float64,
	limit int,
) []*ReverseGeocodingResult {
	rectangles := geo.getNearbyRectangles(lat, lng, accuracyMeters)
	addressIDs := make([]uint32, 0)
	for _, rect := range rectangles {
		key := rect.ToString()
		addressIDsInRectangle, ok := geo.index.Items[key]
		if !ok {
			continue
		}
		addressIDs = append(addressIDs, addressIDsInRectangle...)
	}

	results := make([]*ReverseGeocodingResult, 0)
	for _, id := range addressIDs {
		addr, ok := geo.data.Addresses[id]
		if !ok {
			continue
		}
		distance := getDistance(lat, lng, addr.Lat, addr.Lng)
		if distance <= accuracyMeters {
			fullAddress := geo.resolveAddress(addr)
			if fullAddress.Street1562 != nil {
				results = append(results, &ReverseGeocodingResult{
					fullAddress,
					distance,
				})
			}
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
	safeLimit := limit
	if len(results) < limit {
		safeLimit = len(results)
	}
	results = results[:safeLimit]
	return results
}
