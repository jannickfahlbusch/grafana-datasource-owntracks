package datasource

import (
	"context"
	"errors"
	"reflect"

	"github.com/jannickfahlbusch/owntracks-go/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (owntracksDatasource *OwntracksDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	healthCheckResult := &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working as expected",
	}

	h, err := owntracksDatasource.instanceManager.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	log.DefaultLogger.Info("Healthcheck", "type", reflect.TypeOf(h))

	settings, ok := h.(*datasourceSettings)
	if !ok {
		healthCheckResult.Status = backend.HealthStatusUnknown
		healthCheckResult.Message = "Could not check connection"
		return healthCheckResult, errors.New("did not received settings properly")
	}

	// Create a new client
	owntracksClient := client.New(settings.Settings().URL + APISuffix)
	_, err = owntracksClient.Users(ctx)
	if err != nil {
		healthCheckResult.Status = backend.HealthStatusError
		healthCheckResult.Message = err.Error()
		return healthCheckResult, err
	}

	return healthCheckResult, nil
}
