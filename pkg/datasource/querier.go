package datasource

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jannickfahlbusch/owntracks-go/types"

	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/jannickfahlbusch/owntracks-go/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (owntracksDatasource *OwntracksDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData starting")

	// create response struct
	response := backend.NewQueryDataResponse()

	otClient, err := owntracksDatasource.Client(req.PluginContext)
	if err != nil {
		log.DefaultLogger.Error("Failed to obtain Datasource client", "error", err)
		return nil, err
	}

	// loop over queries and execute them individually.
	for _, query := range req.Queries {
		queryResponse := owntracksDatasource.query(ctx, otClient, query)
		response.Responses[query.RefID] = queryResponse
	}

	return response, nil
}

func (owntracksDatasource *OwntracksDatasource) query(ctx context.Context, client client.Client, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}
	request := &Query{}

	response.Error = json.Unmarshal(query.JSON, request)
	if response.Error != nil {
		log.DefaultLogger.Error("Failed to unmarshal query into request", "error", response.Error)
		return response
	}

	from := query.TimeRange.From
	to := query.TimeRange.To

	var locations *types.LocationList
	locations, response.Error = client.Locations(ctx, request.User, request.Device, from, to)
	if response.Error != nil {
		log.DefaultLogger.Error("Failed to gather location data", "error", response.Error)
		return response
	}

	log.DefaultLogger.Info("Got locations response", "amountLocations", locations.Count)

	// Table Frame
	/*var tableFrame *data.Frame
	tableFrame, response.Error = toTable(locations)
	if response.Error != nil {
		return response
	}

	response.Frames = append(response.Frames, tableFrame)*/

	// Time Series Frame
	var timeSeriesFrame *data.Frame
	timeSeriesFrame, response.Error = toTimeSeries(locations)
	response.Frames = append(response.Frames, timeSeriesFrame)

	return response
}

func toTimeSeries(locations *types.LocationList) (*data.Frame, error) {
	frame := data.NewFrame("location",
		data.NewField("time", nil, make([]*time.Time, locations.Count)),
		data.NewField("latitude", nil, make([]float64, locations.Count)),
		data.NewField("longitude", nil, make([]float64, locations.Count)),
		data.NewField("velocity", nil, make([]float64, locations.Count)),
		data.NewField("geohash", nil, make([]string, locations.Count)),
		data.NewField("altitude", nil, make([]float64, locations.Count)),
		data.NewField("accuracy", nil, make([]float64, locations.Count)),
		data.NewField("verticalAccuracy", nil, make([]float64, locations.Count)),
		data.NewField("address", nil, make([]string, locations.Count)),
		data.NewField("locality", nil, make([]string, locations.Count)),
		data.NewField("countryCode", nil, make([]string, locations.Count)),
	)

	for index, location := range locations.Data {
		frame.Set(0, index, location.Timestamp)
		frame.Set(1, index, location.Latitude)
		frame.Set(2, index, location.Longitude)
		frame.Set(3, index, location.Velocity)
		frame.Set(4, index, location.GeoHash)
		frame.Set(5, index, location.Altitude)
		frame.Set(6, index, location.Accuracy)
		frame.Set(7, index, location.VerticalAccuracy)
		frame.Set(8, index, location.Address)
		frame.Set(9, index, location.Locality)
		frame.Set(10, index, location.CountryCode)
	}

	return frame, nil
}

func toTable(locations *types.LocationList) (*data.Frame, error) {
	columns := []string{
		"time",
		"latitude",
		"longitude",
		"altitude",
		"velocity",
		"geohash",
		"accuracy",
		"radius",
		"verticalAccuracy",
		"barometricPressure",
	}

	frame := data.NewFrameOfFieldTypes("Response",
		0,
		data.FieldTypeTime,    // time
		data.FieldTypeFloat64, // latitude
		data.FieldTypeFloat64, // longitude
		data.FieldTypeFloat64, // altitude
		data.FieldTypeFloat64, // velocity
		data.FieldTypeString,  // geohash
		data.FieldTypeFloat64, // accuracy
	)

	err := frame.SetFieldNames(columns...)
	if err != nil {
		return nil, err
	}

	for _, location := range locations.Data {
		frame.AppendRow(location.Timestamp, location.Latitude, location.Longitude, location.Altitude, location.Velocity, location.GeoHash)
	}

	return frame, nil
}
