package datasource

import (
	"context"
	"encoding/json"
	"time"
	"unsafe"

	"github.com/jannickfahlbusch/owntracks-go/types"

	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/jannickfahlbusch/owntracks-go/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func (owntracksDatasource *OwntracksDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	otClient, err := owntracksDatasource.Client(req.PluginContext)
	if err != nil {
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
		return response
	}

	from := query.TimeRange.From
	to := query.TimeRange.To

	labels := data.Labels{
		"User":   request.User,
		"Device": request.Device,
	}

	var locations *types.LocationList
	locations, response.Error = client.Locations(ctx, request.User, request.Device, from, to)
	if response.Error != nil {
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
	timeSeriesFrame, response.Error = toTimeSeries(locations, labels)
	response.Frames = append(response.Frames, timeSeriesFrame)

	log.DefaultLogger.Info("Returning frames", "length", len(response.Frames), "size", unsafe.Sizeof(response))

	return response
}

func toTimeSeries(locations *types.LocationList, labels Labels) (*data.Frame, error) {
	frame := data.NewFrame("location",
		data.NewField("time", labels, make([]time.Time, locations.Count)),
		data.NewField("latitude", labels, make([]float64, locations.Count)),
		data.NewField("longitude", labels, make([]float64, locations.Count)),
		data.NewField("geohash", labels, make([]string, locations.Count)),
		data.NewField("velocity", labels, make([]int32, locations.Count)),
		data.NewField("altitude", labels, make([]float64, locations.Count)),
	)

	for index, location := range locations.Data {
		timestamp := time.Unix(location.Timestamp, 0)
		frame.Set(0, index, timestamp)
		frame.Set(1, index, location.Latitude)
		frame.Set(2, index, location.Longitude)
		frame.Set(3, index, location.GeoHash)
		frame.Set(4, index, int32(location.Velocity))
		frame.Set(5, index, location.Altitude)
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
		data.FieldTypeInt32,   // velocity
		data.FieldTypeString,  // geohash
		data.FieldTypeFloat64, // accuracy
	)

	err := frame.SetFieldNames(columns...)
	if err != nil {
		return nil, err
	}

	for _, location := range locations.Data {
		timestamp := time.Unix(location.Timestamp, 0)

		frame.AppendRow(timestamp, location.Latitude, location.Longitude, location.Altitude, int32(location.Velocity), location.GeoHash)
		frame.Row
	}

	return frame, nil
}
