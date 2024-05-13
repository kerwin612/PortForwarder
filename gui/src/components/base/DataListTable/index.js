import React, { forwardRef } from 'react';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';

export default forwardRef( ( props, ref ) => {
    return (
        <DataTable ref={ref} {...props.table}>
            {props.columns.map( ( { key, ...col } ) => <Column key={key} {...col}></Column> )}
        </DataTable>
    );
} );
