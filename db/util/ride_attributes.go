package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/tkrajina/gpxgo/gpx"
)

// RideAttributes are derived from a trail's GPX file when it is uploaded.
// Empty strings mean "could not be detected"; the user can always set them by hand.
type RideAttributes struct {
	RouteType   string // "loop" | "traversal" | ""
	StartsHome  bool
	RideProfile string // one of the ride_profile select values or ""
}

const (
	defaultHomeRadiusMeters     = 100.0
	defaultLoopMaxDistanceMeter = 10000.0
)

var bikerouterProfileRe = regexp.MustCompile(`[?&;]profile=([A-Za-z0-9_.\-]+)`)

// bikerouter profile ids (from the export link) -> our ride_profile values
var bikerouterProfiles = map[string]string{
	"quaelnix-gravel": "quaelnix-gravel",
	"trekking":        "trekking",
}

func mapBikerouterProfile(id string) string {
	id = strings.ToLower(id)
	if v, ok := bikerouterProfiles[id]; ok {
		return v
	}
	switch {
	case strings.HasPrefix(id, "m11n"):
		return "m11n-gravel"
	case strings.Contains(id, "quaelnix"):
		return "quaelnix-gravel"
	case strings.Contains(id, "trekking"):
		return "trekking"
	}
	return ""
}

func envFloat(name string, fallback float64) (float64, bool) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback, false
	}
	return v, true
}

// DetectRideAttributes inspects raw GPX content.
func DetectRideAttributes(content []byte) (RideAttributes, error) {
	attrs := RideAttributes{}

	if m := bikerouterProfileRe.FindSubmatch(content); m != nil {
		attrs.RideProfile = mapBikerouterProfile(string(m[1]))
	}
	if attrs.RideProfile == "" {
		// older bikerouter exports put the profile name into the track name
		lower := bytes.ToLower(content)
		if bytes.Contains(lower, []byte("rennrad low traffic")) {
			attrs.RideProfile = "road-low-traffic"
		}
	}

	data, err := gpx.ParseBytes(content)
	if err != nil {
		return attrs, fmt.Errorf("parse gpx: %w", err)
	}

	var first, last *gpx.GPXPoint
	for ti := range data.Tracks {
		for si := range data.Tracks[ti].Segments {
			pts := data.Tracks[ti].Segments[si].Points
			if len(pts) == 0 {
				continue
			}
			if first == nil {
				first = &pts[0]
			}
			last = &pts[len(pts)-1]
		}
	}
	if first == nil {
		for ri := range data.Routes {
			pts := data.Routes[ri].Points
			if len(pts) == 0 {
				continue
			}
			if first == nil {
				first = &pts[0]
			}
			last = &pts[len(pts)-1]
		}
	}
	if first == nil || last == nil {
		return attrs, nil
	}

	loopMax, _ := envFloat("RIDE_LOOP_MAX_DISTANCE_M", defaultLoopMaxDistanceMeter)
	if HaversineDistanceMeters(first.Latitude, first.Longitude, last.Latitude, last.Longitude) <= loopMax {
		attrs.RouteType = "loop"
	} else {
		attrs.RouteType = "traversal"
	}

	homeLat, okLat := envFloat("RIDE_HOME_LAT", 0)
	homeLon, okLon := envFloat("RIDE_HOME_LON", 0)
	if okLat && okLon {
		radius, _ := envFloat("RIDE_HOME_RADIUS_M", defaultHomeRadiusMeters)
		attrs.StartsHome = HaversineDistanceMeters(first.Latitude, first.Longitude, homeLat, homeLon) <= radius
	}

	return attrs, nil
}

func readTrailGPX(app core.App, r *core.Record) ([]byte, error) {
	gpxPath := r.GetString("gpx")
	if gpxPath == "" {
		return nil, nil
	}
	fsys, err := app.NewFilesystem()
	if err != nil {
		return nil, fmt.Errorf("open filesystem: %w", err)
	}
	defer fsys.Close()

	reader, err := fsys.GetReader(r.BaseFilesPath() + "/" + gpxPath)
	if err != nil {
		return nil, fmt.Errorf("open gpx file: %w", err)
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

// ApplyRideAttributes detects ride attributes once per trail and fills only
// fields the user has not set. It marks the trail as processed so manual
// changes are never overwritten later. It does not save the record.
func ApplyRideAttributes(app core.App, r *core.Record) (bool, error) {
	if r.GetBool("ride_attributes_detected") {
		return false, nil
	}
	content, err := readTrailGPX(app, r)
	if err != nil {
		return false, err
	}
	r.Set("ride_attributes_detected", true)
	if len(content) == 0 {
		return true, nil
	}
	attrs, err := DetectRideAttributes(content)
	if err != nil {
		return true, err
	}
	if r.GetString("route_type") == "" && attrs.RouteType != "" {
		r.Set("route_type", attrs.RouteType)
	}
	if !r.GetBool("starts_home") && attrs.StartsHome {
		r.Set("starts_home", true)
	}
	if r.GetString("ride_profile") == "" && attrs.RideProfile != "" {
		r.Set("ride_profile", attrs.RideProfile)
	}
	return true, nil
}
