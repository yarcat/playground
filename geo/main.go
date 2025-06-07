package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// ExampleStats demonstrates how to read a full file and gather some stats.
// This is similar to `osmconvert --out-statistics`
func main() {
	f, err := os.Open("/mnt/d/yarca/Downloads/latvia-latest.osm.pbf")
	if err != nil {
		fmt.Printf("could not open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	nodes, ways, relations := 0, 0, 0
	stats := newElementStats()

	minLat, maxLat := math.MaxFloat64, -math.MaxFloat64
	minLon, maxLon := math.MaxFloat64, -math.MaxFloat64

	minTS, maxTS := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC), time.Time{}

	var (
		maxNodeRefs   int
		maxNodeRefsID osm.WayID
	)

	var (
		maxRelRefs   int
		maxRelRefsID osm.RelationID
	)

	scanner := osmpbf.New(context.Background(), f, 3)
	defer scanner.Close()

	for scanner.Scan() {
		var ts time.Time

		switch e := scanner.Object().(type) {
		case *osm.Node:
			nodes++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if e.Lat > maxLat {
				maxLat = e.Lat
			}

			if e.Lat < minLat {
				minLat = e.Lat
			}

			if e.Lon > maxLon {
				maxLon = e.Lon
			}

			if e.Lon < minLon {
				minLon = e.Lon
			}
		case *osm.Way:
			data, _ := json.MarshalIndent(e, "", "  ")
			if bytes.Contains(bytes.ToLower(data), []byte("brīvības")) {
				fmt.Println(string(data))
			}
			ways++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if l := len(e.Nodes); l > maxNodeRefs {
				maxNodeRefs = l
				maxNodeRefsID = e.ID
			}
		case *osm.Relation:
			relations++
			ts = e.Timestamp
			stats.Add(e.ElementID(), e.Tags)

			if l := len(e.Members); l > maxRelRefs {
				maxRelRefs = l
				maxRelRefsID = e.ID
			}
		}

		if ts.After(maxTS) {
			maxTS = ts
		}

		if ts.Before(minTS) {
			minTS = ts
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner returned error: %v", err)
		os.Exit(1)
	}

	fmt.Println("timestamp min:", minTS.Format(time.RFC3339))
	fmt.Println("timestamp max:", maxTS.Format(time.RFC3339))
	fmt.Printf("lon min: %0.7f\n", minLon)
	fmt.Printf("lon max: %0.7f\n", maxLon)
	fmt.Printf("lat min: %0.7f\n", minLat)
	fmt.Printf("lat max: %0.7f\n", maxLat)
	fmt.Println("nodes:", nodes)
	fmt.Println("ways:", ways)
	fmt.Println("relations:", relations)
	fmt.Println("version max:", stats.MaxVersion)
	fmt.Println("node id min:", stats.Ranges[osm.TypeNode].Min)
	fmt.Println("node id max:", stats.Ranges[osm.TypeNode].Max)
	fmt.Println("way id min:", stats.Ranges[osm.TypeWay].Min)
	fmt.Println("way id max:", stats.Ranges[osm.TypeWay].Max)
	fmt.Println("relation id min:", stats.Ranges[osm.TypeRelation].Min)
	fmt.Println("relation id max:", stats.Ranges[osm.TypeRelation].Max)
	fmt.Println("keyval pairs max:", stats.MaxTags)
	fmt.Println("keyval pairs max object:", stats.MaxTagsID.Type(), stats.MaxTagsID.Ref())
	fmt.Println("noderefs max:", maxNodeRefs)
	fmt.Println("noderefs max object: way", maxNodeRefsID)
	fmt.Println("relrefs max:", maxRelRefs)
	fmt.Println("relrefs max object: relation", maxRelRefsID)

}

// Stats is a shared bit of code to accumulate stats from the element ids.
type elementStats struct {
	Ranges     map[osm.Type]*idRange
	MaxVersion int

	MaxTags   int
	MaxTagsID osm.ElementID
}

type idRange struct {
	Min, Max int64
}

func newElementStats() *elementStats {
	return &elementStats{
		Ranges: map[osm.Type]*idRange{
			osm.TypeNode:     {Min: math.MaxInt64},
			osm.TypeWay:      {Min: math.MaxInt64},
			osm.TypeRelation: {Min: math.MaxInt64},
		},
	}
}

func (s *elementStats) Add(id osm.ElementID, tags osm.Tags) {
	s.Ranges[id.Type()].Add(id.Ref())

	if v := id.Version(); v > s.MaxVersion {
		s.MaxVersion = v
	}

	if l := len(tags); l > s.MaxTags {
		s.MaxTags = l
		s.MaxTagsID = id
	}
}

func (r *idRange) Add(ref int64) {
	if ref > r.Max {
		r.Max = ref
	}

	if ref < r.Min {
		r.Min = ref
	}
}
