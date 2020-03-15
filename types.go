package geocoder

type GeoIndex map[string][]uint32
type GeoIndexWithResolution struct {
	Resolution float64  `json:"resolution"`
	Items      GeoIndex `json:"index"`
}

// Address represents an address from Kharkiv Architectural Registry database
type Address struct {
	ID     uint32  `json:"id,omitempty"`
	Lat    float64 `json:"lat,omitempty"`
	Lng    float64 `json:"lng,omitempty"`
	Number uint16  `json:"number,omitempty"`
	// Suffix - a letter after the building number. Example - "Г"
	// Suffix can also mean a secondary building number. Example - "/26"
	Suffix string `json:"suffix,omitempty"`
	// Block - block (корпус)
	Block    string `json:"block,omitempty"`
	StreetID uint32 `json:"streetID,omitempty"`
	// Detail - example "ділянка №"
	Detail string `json:"detail,omitempty"`
	// DetailNumber - example "5" (works together with detail)
	DetailNumber string `json:"detailNumber,omitempty"`
	Postcode     uint32 `json:"postcode,omitempty"`
}

type AddressMap map[uint32]*Address

type StreetAR struct {
	ID     uint32 `json:"id,omitempty"`
	NameUk string `json:"name_ukr,omitempty"`
	NameRu string `json:"name_ru,omitempty"`
	TypeUk string `json:"typeUKR,omitempty"`
	TypeRu string `json:"typeRU,omitempty"`
}

type StreetsARMap map[uint32]*StreetAR

type Street1562 struct {
	ID   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Streets1562Map map[uint32]*Street1562

type IDToIDMap map[uint32]uint32

type GeocoderData struct {
	Addresses       AddressMap
	StreetsAR       StreetsARMap
	Streets1562     Streets1562Map
	MappingArTo1562 IDToIDMap
	Mapping1562ToAr IDToIDMap
}

type FullAddress struct {
	Address    *Address
	Street1562 *Street1562
	StreetAR   *StreetAR
}

type ReverseGeocodingResult struct {
	FullAddress *FullAddress
	Distance    float64
}
