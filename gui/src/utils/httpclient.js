import axios from 'axios';

const TIMEOUT = 5000;
const baseURL = process.env.REACT_APP_PORTFORWARDER_API_BASE_URL || '';

export const instance = axios.create( {
    baseURL,
    timeout: TIMEOUT,
} );

instance.interceptors.response.use(
    ( r ) => r,
    ( e ) => {
        if ( e?.response?.status === 500 ) {
            return Promise.resolve( e.response );
        } else {
            return Promise.reject( e );
        }
    },
);

export default instance;
