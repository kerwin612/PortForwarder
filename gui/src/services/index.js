import { httpclient } from 'utils';

export const IPS = ( [{name: '0.0.0.0(All local network addresses)', code: ''}, ...( ( ( await httpclient.get( '/v1/forward/ip/list' ) )?.data || [] ).map( i => ( {name: i, code: i} ) ) )] );

export const PROTOCOLS = ( ( await httpclient.get( '/v1/forward/protocol' ) )?.data || [] ).map( i => ( {name: i.toUpperCase(), code: i} ) );

export const exit = async () => {
    await httpclient.get( '/exit' );
};

export const explorerWS = async () => {
    return ( await httpclient.get( '/v1/ws/explorer' ) )?.data || [];
};

export const getSettings = async () => {
    return ( await httpclient.get( '/v1/settings' ) )?.data || [];
};

export const saveSettings = async ( data ) => {
    return ( await httpclient.post( '/v1/settings', data ) )?.data || {};
};

export const telnet = async ( {protocol, ip, port} ) => {
    return await httpclient.get( `/v1/telnet/${protocol}/${ip}/${port}` );
};

export const getForwardList = async ( params ) => {
    const result = await httpclient.get( '/v1/forward/list' );

    let data = result?.data ?? {};
    return Object.keys( data ).map( k => {
        let i = data[k];
        let o = i[0];
        o.id = k;
        o.status = i[1];
        return o;
    } );
};

export const stopForwards = async ( data ) => {
    return await httpclient.post( '/v1/forward/stop', data );
};

export const deleteForwards = async ( data ) => {
    return await httpclient.post( '/v1/forward/delete', data );
};

export const startForwards = async ( data ) => {
    return await httpclient.post( '/v1/forward/start', data );
};

export const restartForwards = async ( data ) => {
    return await httpclient.post( '/v1/forward/restart', data );
};

export const saveForward = async ( data ) => {
    return await httpclient.post( `/v1/forward/save/${data.id}`, data );
};

export const saveAndStartForward = async ( data ) => {
    return await httpclient.post( `/v1/forward/save_and_start/${data.id}`, data );
};
