package myutils

import (
	DS "pudding/datastructs"
)

type GeoRect struct {
	minLon float64
	maxLon float64
	minLat float64
	maxLat float64
}

func (r *GeoRect) inBound(lon, lat float64) bool {
	return lon >= r.minLon &&
		lon <= r.maxLon &&
		lat >= r.minLat &&
		lat <= r.maxLat
}

func (r *GeoRect) isValid() bool {
	return !(r.minLon == 0 ||
		r.maxLon == 0 ||
		r.minLat == 0 ||
		r.maxLat == 0)
}

func WithinRect(update *DS.Update, query *DS.Query) bool {
	br := &GeoRect{query.Lower_left.Lon,
		query.Top_right.Lon,
		query.Lower_left.Lat,
		query.Top_right.Lat}

	if br.isValid() {
		return br.inBound(update.Start_coord.Lon, update.Start_coord.Lat) ||
			br.inBound(update.End_coord.Lon, update.End_coord.Lat)
	} else {
		return true
	}

}
