import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, OwntracksDatasourceOptions, OwntracksQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, OwntracksQuery, OwntracksDatasourceOptions>;

export class VariableQueryEditor extends PureComponent<Props> {
  onUserChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, user: event.target.value });
    // executes the query
    onRunQuery();
  };

  onDeviceChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, device: event.target.value });
    // executes the query

  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { user, device } = query;

    console.log('query');

    return (
      <div className="gf-form">
        <FormField width={4} value={user} onChange={this.onUserChange} label="User" type="text" />
        <FormField width={4} value={device} onChange={this.onDeviceChange} label="Device" type="text" />
      </div>
    );
  }
}
