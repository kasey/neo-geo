package parser

import (
	"fmt"
	"github.com/golang/geo/s2"
	"testing"
)

func AsS2Point(geo Geometry, t *testing.T) (s2.Point, error) {
	switch geo.Type() {
	case "Point":
		gjPoint, err := AsPoint(geo)
		if err != nil {
			return s2.Point{}, err
		}
		ll := s2.LatLngFromDegrees(gjPoint.Coordinate[0], gjPoint.Coordinate[1])
		return s2.PointFromLatLng(ll), nil
	default:
		return s2.Point{}, fmt.Errorf("Can't construct a s2.Point type from a geo-json %v type", geo.Type())
	}
}

func S2PointFromCoordinate(coord Coordinate) s2.Point {
	ll := s2.LatLngFromDegrees(coord[0], coord[1])
	return s2.PointFromLatLng(ll)
}

func AsS2PointSlice(geo Geometry) ([]s2.Point, error) {
	points := make([]s2.Point, 0)
	switch geo.Type() {
	case "LineString":
		ls, err := AsLineString(geo)
		if err != nil {
			return nil, err
		}

		for i, _ := range ls.Coordinates {
			points = append(points, S2PointFromCoordinate(ls.Coordinates[i]))
		}
		return points, nil
	case "MultiPoint":
		mp, err := AsMultiPoint(geo)
		if err != nil {
			return nil, err
		}

		for i, _ := range mp.Coordinates {
			points = append(points, S2PointFromCoordinate(mp.Coordinates[i]))
		}
		return points, nil
	default:
		return nil, fmt.Errorf("Can't make an []s2.Point slice from type=%v", geo.Type())
	}
}

// AsS2Polygon is built on top of S2PointFromCoordinate, LoopFromPoints and PolygonFromLoops.
// BUG: The documentation for PolygonFromLoops comes with the following warning:
// "NOTE: this function is NOT YET IMPLEMENTED for more than one loop and will panic if given a slice of length > 1."
// link: https://godoc.org/github.com/golang/geo/s2#PolygonFromLoops
func AsS2Polygon(geo Geometry) (*s2.Polygon, error) {
	switch geo.Type() {
	case "Polygon":
		pg, err := AsPolygon(geo)
		if err != nil {
			return nil, err
		}

		loops := make([]*s2.Loop, len(pg.Coordinates))
		for i, _ := range pg.Coordinates {
			coordLoop := pg.Coordinates[i]
			points := make([]s2.Point, len(coordLoop))
			for pi, _ := range coordLoop {
				points[pi] = S2PointFromCoordinate(coordLoop[pi])
			}
			loops[i] = s2.LoopFromPoints(points)
		}
		return s2.PolygonFromLoops(loops), nil
	default:
		return nil, fmt.Errorf("Can't construct an s2.Polygon from type=%v", geo.Type())
	}
}

func MultiLineStringToS2PointSliceSlice(geo Geometry) ([][]s2.Point, error) {
	return make([][]s2.Point, 0), nil
}

func MultiPolygonToS2PolygonSlice(geo Geometry) ([]*s2.Polygon, error) {
	return make([]*s2.Polygon, 0), nil
}
