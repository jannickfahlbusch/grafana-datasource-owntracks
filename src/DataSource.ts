import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { OwntracksDatasourceOptions, OwntracksQuery } from './types';

export class DataSource extends DataSourceWithBackend<OwntracksQuery, OwntracksDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<OwntracksDatasourceOptions>) {
    super(instanceSettings);
  }
}
