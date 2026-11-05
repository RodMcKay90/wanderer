package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Ride attributes for the bike-only fork: loop vs. point-to-point,
// whether the ride starts at home, and the bikerouter profile used.
func init() {
	m.Register(upAddRideAttributes1795000000, downAddRideAttributes1795000000)
}

var rideAttributeFields1795000000 = []string{
	`{
		"hidden": false,
		"id": "selectroutetype",
		"maxSelect": 1,
		"name": "route_type",
		"presentable": false,
		"required": false,
		"system": false,
		"type": "select",
		"values": ["loop", "traversal"]
	}`,
	`{
		"hidden": false,
		"id": "boolstartshome",
		"name": "starts_home",
		"presentable": false,
		"required": false,
		"system": false,
		"type": "bool"
	}`,
	`{
		"hidden": false,
		"id": "selectrideprof",
		"maxSelect": 1,
		"name": "ride_profile",
		"presentable": false,
		"required": false,
		"system": false,
		"type": "select",
		"values": ["road-low-traffic", "ffm-by-bicycle", "trekking", "quaelnix-gravel", "m11n-gravel", "other"]
	}`,
	`{
		"hidden": false,
		"id": "boolrideauto",
		"name": "ride_attributes_detected",
		"presentable": false,
		"required": false,
		"system": false,
		"type": "bool"
	}`,
}

func upAddRideAttributes1795000000(app core.App) error {
	trails, err := app.FindCollectionByNameOrId("trails")
	if err != nil {
		return err
	}
	names := []string{"route_type", "starts_home", "ride_profile", "ride_attributes_detected"}
	for i, field := range rideAttributeFields1795000000 {
		if trails.Fields.GetByName(names[i]) != nil {
			continue
		}
		if err := trails.Fields.AddMarshaledJSON([]byte(field)); err != nil {
			return err
		}
	}
	return app.Save(trails)
}

func downAddRideAttributes1795000000(app core.App) error {
	trails, err := app.FindCollectionByNameOrId("trails")
	if err != nil {
		return err
	}
	for _, name := range []string{"route_type", "starts_home", "ride_profile", "ride_attributes_detected"} {
		trails.Fields.RemoveByName(name)
	}
	return app.Save(trails)
}
