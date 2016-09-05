package parser

import (
	"encoding/json"
	"fmt"
)

type jsonPeek struct {
	TypeName string `json:"type"`
	raw      []byte
}

// GeoJSON types are based on a subset of OpenGIS types, so
// formal definitions are split between the GEOJSON rfc:
// https://tools.ietf.org/html/rfc7946
// and the OpenGIS Simple Features Specification for SQL:
// https://portal.opengeospatial.org/files/?artifact_id=829
type Geometry interface {
	Type() string
	Raw() []byte
}

func (peek *jsonPeek) Type() string {
	return peek.TypeName
}

func (peek *jsonPeek) Raw() []byte {
	return peek.raw
}

func NewGeometry(geoJSON []byte) (Geometry, error) {
	peek := &jsonPeek{raw: geoJSON}
	err := json.Unmarshal(geoJSON, peek)
	if err != nil {
		return nil, err
	}
	return peek, nil
}

type Coordinate []float64

type Point struct {
	Type       string     `json:"type"`
	Coordinate Coordinate `json:"coordinates"`
}

func AsPoint(geo Geometry) (*Point, error) {
	if geo.Type() != "Point" {
		return nil, fmt.Errorf("Can't create a Point from type='%v'", geo.Type())
	}
	pt := new(Point)
	err := json.Unmarshal(geo.Raw(), pt)

	return pt, err
}

type LineString struct {
	Type        string       `json:"type"`
	Coordinates []Coordinate `json:"coordinates"`
}

func AsLineString(geo Geometry) (*LineString, error) {
	if geo.Type() != "LineString" {
		return nil, fmt.Errorf("Can't create a LineString from type='%v'", geo.Type())
	}
	ls := new(LineString)
	err := json.Unmarshal(geo.Raw(), ls)

	return ls, err
}

type Polygon struct {
	Type        string         `json:"type"`
	Coordinates [][]Coordinate `json:"coordinates"`
}

func AsPolygon(geo Geometry) (*Polygon, error) {
	if geo.Type() != "Polygon" {
		return nil, fmt.Errorf("Can't create a Polygon from type='%v'", geo.Type())
	}
	pg := new(Polygon)
	err := json.Unmarshal(geo.Raw(), pg)

	return pg, err
}

type MultiPoint struct {
	Type        string       `json:"type"`
	Coordinates []Coordinate `json:"coordinates"`
}

func AsMultiPoint(geo Geometry) (*MultiPoint, error) {
	if geo.Type() != "MultiPoint" {
		return nil, fmt.Errorf("Can't create a MultiPoint from type='%v'", geo.Type())
	}

	mp := new(MultiPoint)
	err := json.Unmarshal(geo.Raw(), mp)

	return mp, err
}

type MultiLineString struct {
	Type        string         `json:"type"`
	Coordinates [][]Coordinate `json:"coordinates"`
}

func AsMultiLineString(geo Geometry) (*MultiLineString, error) {
	if geo.Type() != "MultiLineString" {
		return nil, fmt.Errorf("Can't create a MultiLineSring from type='%v'", geo.Type())
	}

	mls := new(MultiLineString)
	err := json.Unmarshal(geo.Raw(), mls)

	return mls, err
}

type MultiPolygon struct {
	Type        string           `json:"type"`
	Coordinates [][][]Coordinate `json:"coordinates"`
}

func AsMultiPolygon(geo Geometry) (*MultiPolygon, error) {
	if geo.Type() != "MultiPolygon" {
		return nil, fmt.Errorf("Can't create a MultiPolygon from type='%v'", geo.Type())
	}

	mpg := new(MultiPolygon)
	err := json.Unmarshal(geo.Raw(), mpg)

	return mpg, err
}

type GeometryCollection struct {
	Type           string            `json:"type"`
	AnonymousGeoms []json.RawMessage `json:"geometries"`
	Geometries     []Geometry
}

func AsGeometryCollection(geo Geometry) (*GeometryCollection, error) {
	if geo.Type() != "GeometryCollection" {
		return nil, fmt.Errorf("Can't create a GeometryCollection from type='%v'", geo.Type())
	}

	gc := new(GeometryCollection)
	err := json.Unmarshal(geo.Raw(), gc)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error ='%v' for GeometryCollection=%v", err, string(geo.Raw()))
	}

	geoms := make([]Geometry, len(gc.AnonymousGeoms))
	for i, _ := range gc.AnonymousGeoms {
		g, err := NewGeometry(gc.AnonymousGeoms[i])
		if err != nil {
			return gc, err
		}
		geoms[i] = g
	}
	gc.Geometries = geoms

	return gc, err
}

func AsMap(geo Geometry) (map[string]interface{}, error) {
	geomap := make(map[string]interface{})
	err := json.Unmarshal(geo.Raw(), &geomap)
	return geomap, err
}
