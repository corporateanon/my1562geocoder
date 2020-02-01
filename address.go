package my1562geocoder

import "strconv"

func (addr *Address) GetBuildingAsString() string {
	building := ""
	if addr.Number != 0 {
		building = building + strconv.FormatInt(int64(addr.Number), 10)
	}
	if addr.Suffix != "" {
		building = building + addr.Suffix
	}
	if addr.Block != "" {
		building = building + " к. №" + addr.Block
	}
	if addr.Detail != "" {
		building = building + " " + addr.Detail
		if addr.DetailNumber != "" {
			building = building + addr.DetailNumber
		}
	}
	return building
}
