import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface OwntracksQuery extends DataQuery {
  user: string;
  device: string;
}

export const defaultQuery: Partial<OwntracksQuery> = {
  user: 'tracking',
  device: 'iphone',
};

/**
 * These are options configured for each DataSource instance
 */
export interface OwntracksDatasourceOptions extends DataSourceJsonData {
  recorderURL?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface OwntracksSecureJSONData {
  apiKey?: string;
}
