
import React, { forwardRef } from 'react';
import { Toast } from 'primereact/toast';

export default forwardRef( ( props, ref ) => {
    return (
        <Toast ref={ref} {...props} />
    );
} );
