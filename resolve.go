package my1562geocoder

type FullAddress struct {
	address    *Address
	street1562 *Street1562
	streetAR   *StreetAR
}

func resolveAddress(address *Address) *FullAddress {
	fa := &FullAddress{}
	fa.address = address
	streetID := address.StreetID

	streetAR, streetARisResolved := StreetsAR[streetID]
	if streetARisResolved {
		fa.streetAR = &streetAR
	}

	street1562Id, street1562IdIsResolved := StreetsARto1562[streetID]
	if street1562IdIsResolved {
		street1562, street1562IsResolved := Streets1562[street1562Id]
		if street1562IsResolved {
			fa.street1562 = &street1562
		}
	}
	return fa
}
