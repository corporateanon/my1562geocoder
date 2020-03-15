package geocoder

import "strconv"

type Rectangle struct {
	ILat int32
	ILng int32
}

func (rect *Rectangle) ToString() string {
	return strconv.Itoa(int(rect.ILat)) + "." + strconv.Itoa(int(rect.ILng))
}

func NewRectangle(lat float64, lng float64, res float64) *Rectangle {
	var iLat, iLng int32
	iLat = int32((lat + 90.0) * res)
	iLng = int32((lng + 180.0) * res)
	return &Rectangle{iLat, iLng}
}
