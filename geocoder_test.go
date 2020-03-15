package geocoder

import (
	"math"
	"math/rand"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"gotest.tools/assert"
)

func TestNewGeocoder(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")

	assert.DeepEqual(t, geo.data.Addresses[100000],
		&Address{
			ID:       100000,
			Lat:      49.963557842,
			Lng:      36.352988433,
			Number:   3,
			StreetID: 2009,
			Postcode: 61099,
		})

	if len(geo.data.Addresses) < 77034 {
		t.Errorf("len(Addresses) is too small %d", len(geo.data.Addresses))
	}
	if len(geo.data.StreetsAR) < 3720 {
		t.Errorf("len(StreetsAR) is too small %d", len(geo.data.StreetsAR))
	}
	if len(geo.data.Streets1562) < 2840 {
		t.Errorf("len(Streets1562) is too small %d", len(geo.data.Streets1562))
	}
	if len(geo.data.MappingArTo1562) < 2662 {
		t.Errorf("len(MappingArTo1562) is too small %d", len(geo.data.MappingArTo1562))
	}
	if len(geo.data.Mapping1562ToAr) < 2689 {
		t.Errorf("len(Mapping1562ToAr) is too small %d", len(geo.data.Mapping1562ToAr))
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

func TestGeocoder_buildSpatialIndex(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(100)
	if len(geo.index.Items) < 400 {
		t.Errorf("len(index.Items) is too small %d", len(geo.index.Items))
	}
}

func Test_getNearbyRectangles(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)
	rects := geo.getNearbyRectangles(49.944204004899994, 36.3421038691, 100)

	assert.DeepEqual(t, rects, []Rectangle{
		{27988, 43268},
		{27989, 43268},
	})
}

func TestReverseGeocode0(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)

	cupaloy.SnapshotT(t, geo.ReverseGeocode(49.977094, 36.219115, 300, 4))
}
func TestReverseGeocode1(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)

	cupaloy.SnapshotT(t, geo.ReverseGeocode(50.018105, 36.331791, 300, 10))
}

func BenchmarkReverseGeocode(b *testing.B) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)
	var latMin float64 = 49.929461
	var latMax float64 = 50.035745
	var lngMin float64 = 36.281144
	var lngMax float64 = 36.371397
	for i := 0; i < b.N; i++ {
		lat := rand.Float64()*(latMax-latMin) + latMin
		lng := rand.Float64()*(lngMax-lngMin) + lngMin
		geo.ReverseGeocode(lat, lng, 300, 10)
	}
}

func TestGeocoder_AddressByID_1(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)
	cupaloy.SnapshotT(t, geo.AddressByID(152699))
}

func TestGeocoder_AddressByID_2(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)
	cupaloy.SnapshotT(t, geo.AddressByID(150690))
}

func TestGeocoder_AddressByID_nil(t *testing.T) {
	geo := NewGeocoder("./data/gobs/geocoder-data.gob")
	geo.BuildSpatialIndex(200)
	cupaloy.SnapshotT(t, geo.AddressByID(180180180))
}
