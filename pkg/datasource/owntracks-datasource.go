package datasource

import (
	"github.com/jannickfahlbusch/owntracks-go/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

const (
	APISuffix = "/api/0/"
)

func NewDatasource() datasource.ServeOpts {
	instanceManager := datasource.NewInstanceManager(newDatasourceSettings)

	ds := &OwntracksDatasource{
		instanceManager: instanceManager,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

func newDatasourceSettings(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &datasourceSettings{
		owntracksClient: client.New(settings.URL + APISuffix),
		settings:        settings,
	}, nil
}
