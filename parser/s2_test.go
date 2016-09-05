package parser

import (
	"testing"

	"github.com/golang/geo/s2"
)

func TestAsS2Point(t *testing.T) {
	ptf := pointFixture()
	geom, err := NewGeometry(ptf)
	if err != nil {
		t.Fatalf("Receieved error from NewGeometry=%v", err)
	}

	s2Point, err := AsS2Point(geom, t)
	if err != nil {
		t.Fatalf("Got error=%v from AsS2Point using input =%v", err, string(ptf))
	}

	pt, err := AsPoint(geom)
	if err != nil {
		t.Fatalf("Got error=%v from AsPoint using input =%v", err, string(ptf))
	}

	ll := s2.LatLngFromDegrees(pt.Coordinate[0], pt.Coordinate[1])

	ptCellID := s2.CellIDFromLatLng(s2.LatLngFromPoint(s2Point))
	llCellID := s2.CellIDFromLatLng(ll)
	if ptCellID != llCellID {
		t.Fatalf("S2/Geo Point mismatch.  CellFromPoint=%v, CellFromLatLng=%v", ptCellID.String(), llCellID.String())
	}
}

func TestLatLngRoundtrip(t *testing.T) {
	llFloats := []float64{102.0, 0.5}
	ll := s2.LatLngFromDegrees(llFloats[0], llFloats[1])
	pt := s2.PointFromLatLng(ll)
	llcell := s2.CellFromLatLng(ll)
	ptcell := s2.CellFromPoint(pt)
	if ptcell != llcell {
		t.Fatalf("Something is busted in the floats -> LatLng -> Point round trip; CellFromLatLng=%v, CellFromPoint=%v", llcell, ptcell)
	}
}
