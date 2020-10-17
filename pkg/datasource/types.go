package datasource

import (
	"errors"

	"github.com/jannickfahlbusch/owntracks-go/client"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

// OwntracksDatasource implements a datasource with Owntracks' Recorder as its backend
type OwntracksDatasource struct {
	instanceManager instancemgmt.InstanceManager
}

func (ds *OwntracksDatasource) Settings(pluginctx backend.PluginContext) (*datasourceSettings, error) {
	instance, err := ds.instanceManager.Get(pluginctx)
	if err != nil {
		return nil, err
	}

	settings, ok := instance.(*datasourceSettings)
	if !ok {
		return nil, errors.New("stored settings do not have correct type")
	}

	return settings, nil

}

func (ds *OwntracksDatasource) Client(pluginctx backend.PluginContext) (client.Client, error) {
	settings, err := ds.Settings(pluginctx)
	if err != nil {
		return nil, err
	}

	return settings.Client(), nil
}

type datasourceSettings struct {
	settings backend.DataSourceInstanceSettings

	owntracksClient client.Client
}

func (settings *datasourceSettings) Settings() backend.DataSourceInstanceSettings {
	return settings.settings
}

func (settings *datasourceSettings) Client() client.Client {
	return settings.owntracksClient
}

type Query struct {
	User   string `json:"user"`
	Device string `json:"device"`
}
