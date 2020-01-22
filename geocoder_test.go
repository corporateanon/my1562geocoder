package my1562geocoder

import (
	"math"
	"reflect"
	"testing"
)

func TestGetAddressById(t *testing.T) {
	type args struct {
		id uint32
	}
	tests := []struct {
		name string
		args args
		want *Address
	}{
		{
			"Existing address",
			args{181901},
			&Address{
				ID:           181901,
				Lat:          49.9604596373,
				Lng:          36.3265315275,
				Number:       4,
				Suffix:       "",
				Block:        "",
				StreetID:     4453,
				Detail:       "",
				DetailNumber: 0,
				Postcode:     61082,
			},
		},
		{
			"Non-existing address",
			args{666666},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAddressByID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRectangleByLatLng(t *testing.T) {
	type args struct {
		lat float64
		lng float64
		res int32
	}
	tests := []struct {
		name string
		args args
		want Rectangle
	}{
		{
			"Resolve coordinates to indexing rectangle",
			args{49.944204004899994, 36.3421038691, 200},
			Rectangle{27988, 43268},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRectangleByLatLng(tt.args.lat, tt.args.lng, tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRectangleByLatLng() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNearbyRectangles(t *testing.T) {
	type args struct {
		lat            float64
		lng            float64
		res            int32
		accuracyMeters float64
	}
	tests := []struct {
		name string
		args args
		want []Rectangle
	}{
		{
			"Get slice of nearby rectangles for a point",
			args{49.944204004899994, 36.3421038691, 200, 100},
			[]Rectangle{
				{27988, 43268},
				{27989, 43268},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNearbyRectangles(tt.args.lat, tt.args.lng, tt.args.res, tt.args.accuracyMeters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNearbyRectangles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDistance(t *testing.T) {
	type args struct {
		lat0 float64
		lng0 float64
		lat1 float64
		lng1 float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"Measures distance",
			args{
				50.00577382744244, 36.22907459735871,
				49.92835056065926, 36.308602094650276,
			},
			10300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := math.Round(getDistance(tt.args.lat0, tt.args.lng0, tt.args.lat1, tt.args.lng1)/100) * 100; got != tt.want {
				t.Errorf("getDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseGeocode(t *testing.T) {
	type args struct {
		lat            float64
		lng            float64
		accuracyMeters float64
		limit          int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Converts coordinates to addresses",
			args{
				49.977094, 36.219115, 100, 4, //Моечная 11/5
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReverseGeocode(tt.args.lat, tt.args.lng, tt.args.accuracyMeters, tt.args.limit)
		})
	}
}
