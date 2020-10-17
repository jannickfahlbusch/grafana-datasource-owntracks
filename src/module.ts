import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { OwntracksQuery, OwntracksDatasourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, OwntracksQuery, OwntracksDatasourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
