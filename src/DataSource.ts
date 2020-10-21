import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { OwntracksDatasourceOptions, OwntracksQuery } from './types';

export class DataSource extends DataSourceWithBackend<OwntracksQuery, OwntracksDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<OwntracksDatasourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: OwntracksQuery): OwntracksQuery {
    const templateSrv = getTemplateSrv();

    return {
      ...query,
      user: query.user ? templateSrv.replace(query.user) : '',
      device: query.device ? templateSrv.replace(query.device) : '',
    };
  }
}
