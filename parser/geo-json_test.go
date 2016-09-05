package parser

import (
	"bytes"
	"testing"
)

func pointFixture() []byte {
	return []byte(`{
		"type": "Point",
		"coordinates": [102.0, 0.5]
	}`)
}

func linestringFixture() []byte {
	return []byte(`{
	"type": "LineString",
	"coordinates": [
		[102.0, 0.0],
		[103.0, 1.0],
		[104.0, 0.0],
		[105.0, 1.0]
		]
	}`)
}

func polygonFixture() []byte {
	return []byte(`{
	"type": "Polygon",
	"coordinates": [
		[
			[100.0, 0.0],
			[101.0, 0.0],
			[101.0, 1.0],
			[100.0, 1.0],
			[100.0, 0.0]
		]
	]}`)
}

func multipointFixture() []byte {
	return []byte(`{
	"type": "MultiPoint",
	"coordinates": [
		[100.0, 0.0],
		[101.0, 1.0]
		]
	}`)
}

func multiLineStringFixture() []byte {
	return []byte(`{
	"type": "MultiLineString",
	"coordinates": [
		[
			[100.0, 0.0],
			[101.0, 1.0]
		],
		[
			[102.0, 2.0],
			[103.0, 3.0]
		]
	]}`)
}

func multiPolygonFixture() []byte {
	return []byte(`{
	"type": "MultiPolygon",
	"coordinates": [
		[
			[
				[180.0, 40.0], [180.0, 50.0], [170.0, 50.0],
				[170.0, 40.0], [180.0, 40.0]
			]
		],
		[
			[
				[-170.0, 40.0], [-170.0, 50.0], [-180.0, 50.0],
				[-180.0, 40.0], [-170.0, 40.0]
			]
		]
	]}`)
}

func geometryCollectionFixture() []byte {
	return []byte(`{
	"type": "GeometryCollection",
	"geometries": [{
		"type": "Point",
		"coordinates": [100.0, 0.0]
		}, {
		"type": "LineString",
		"coordinates": [
			[101.0, 0.0],
			[102.0, 1.0]
		]
		}
	]}`)
}

func TestNewGeometry(t *testing.T) {
	ptf := pointFixture()

	geom, err := NewGeometry(ptf)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}
	if !bytes.Equal(geom.Raw(), ptf) {
		t.Fatalf("Expected geometry.Raw()(%v) to be equal to input (%v)", geom.Raw(), ptf)
	}

	ptype := geom.Type()
	if ptype != "Point" {
		t.Fatalf("GeoJSON type mismatch, wanted='Point', got='%v'", ptype)
	}

	gm, err := AsMap(geom)
	if err != nil {
		t.Fatalf("Receieved error from geo.GeomMap=%v", err)
	}

	expectedCoords := []float64{102.0, 0.5}
	coordinates := gm["coordinates"].([]interface{})
	for i, v := range coordinates {
		if v.(float64) != expectedCoords[i] {
			t.Fatal("Unexpected coordinates from point parse=v%v", coordinates)
		}
	}
}

func TestAsPoint(t *testing.T) {
	ptf := pointFixture()

	geom, err := NewGeometry(ptf)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	point, err := AsPoint(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsPoint=%v %v", err, point)
	}
}

func TestAsLinestring(t *testing.T) {
	lf := linestringFixture()

	geom, err := NewGeometry(lf)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	lineString, err := AsLineString(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsLineString=%v %v", err, lineString)
	}
}

func TestAsPolygon(t *testing.T) {
	pg := polygonFixture()

	geom, err := NewGeometry(pg)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	polygon, err := AsPolygon(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsPolygon=%v %v", err, polygon)
	}
}

func TestAsMultiPoint(t *testing.T) {
	mp := multipointFixture()

	geom, err := NewGeometry(mp)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	multipoint, err := AsMultiPoint(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsMultiPoint=%v %v", err, multipoint)
	}
}

func TestAsMultiLineString(t *testing.T) {
	mp := multiLineStringFixture()

	geom, err := NewGeometry(mp)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	multiLineString, err := AsMultiLineString(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsMultiLineString=%v %v", err, multiLineString)
	}
}

func TestAsMultiPolygon(t *testing.T) {
	mp := multiPolygonFixture()

	geom, err := NewGeometry(mp)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	multiPolygon, err := AsMultiPolygon(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsMultiPolygon=%v %v", err, multiPolygon)
	}
}

func TestAsGeometryCollection(t *testing.T) {
	gc := geometryCollectionFixture()
	geom, err := NewGeometry(gc)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	geometryCollection, err := AsGeometryCollection(geom)
	if err != nil {
		t.Fatalf("Receieved error from AsGeometryCollection=%v %v", err, geometryCollection)
	}
	rawPoint := geometryCollection.Geometries[0]
	for _, g := range geometryCollection.Geometries {
		t.Log(g.Type())
	}
	point, err := AsPoint(rawPoint)
	if err != nil {
		t.Fatalf("Error creating point using AsPoint and first element of .Geometries in GeometryCollection=%v", err)
	}
	expectedCoords := []float64{100.0, 0.0}
	for i, v := range point.Coordinate {
		if v != expectedCoords[i] {
			t.Fatal("Unexpected coordinates from point via GeometryCollection parse=v%v", point.Coordinate)
		}
	}
}
