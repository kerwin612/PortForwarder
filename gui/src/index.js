import React from 'react';
import ReactDOM from 'react-dom/client';

import { PrimeReactProvider } from 'primereact/api';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import 'primereact/resources/themes/lara-light-cyan/theme.css';

import { Content } from 'components/business';

import './index.css';

const root = ReactDOM.createRoot( document.getElementById( 'root' ) );
root.render(
    <React.StrictMode>
        <PrimeReactProvider>
            <Content />
        </PrimeReactProvider>
    </React.StrictMode>
);
