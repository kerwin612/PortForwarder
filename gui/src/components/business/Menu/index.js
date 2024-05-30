import React, { useRef, useState } from 'react';

import { Tag, Menu, Button, Dialog, Tooltip, Message, Dropdown, InputText, InputNumber } from 'components/base';
import { IPS, exit, getSettings, saveSettings, explorerWS } from 'services';
import './index.css';

export default function NavBar() {
    const menuRight = useRef( null );
    const [hostInvalid, setHostInvalid] = useState( false );
    const [portInvalid, setPortInvalid] = useState( false );
    const [settings, setSettings] = useState( {} );
    const [workspace, setWorkspace] = useState( null );
    const [byeDialog, setByeDialog] = useState( false );
    const [selectedIp, setSelectedIp] = useState( null );
    const [aboutDialog, setAboutDialog] = useState( false );
    const [showSuccess, setShowSuccess] = useState( false );
    const [settingsDialog, setSettingsDialog] = useState( false );
    const items = [
        {
            label: 'Settings',
            command: async () => {
                let rst = await getSettings();
                let ws = rst[0];
                let si = ( typeof rst[1] ) !== 'object' ? {} : rst[1];
                setSettings( si );
                setWorkspace( ws );
                setSelectedIp( si.ip ? si.ip : IPS[0] );
                let hostWithPort = window.location.host;
                let pos = hostWithPort.lastIndexOf( ':' );
                let lh = pos > 0 ? hostWithPort.substring( 0, pos ) : hostWithPort;
                let lp = pos > 0 ? hostWithPort.substring( pos + 1 ) : '80';
                if ( si.ip && si.ip !== lh ) {
                    setHostInvalid( true );
                }
                if ( si.port !== 0 && si.port !== parseInt( lp ) ) {
                    setPortInvalid( true );
                }
                setSettingsDialog( true );
            }
        },
        {
            label: 'About',
            command: () => {
                setAboutDialog( true );
            }
        },
        {
            label: 'Quit',
            command: () => {
                exit();
                setByeDialog( true );
            }
        }
    ];

    let timeout = null;
    const closeSuccess = () => {
        if ( timeout != null ) {
            clearTimeout( timeout );
            timeout = null;
        }
        timeout = setTimeout( () => {
            setShowSuccess( false );
        }, 3 * 1000 );
    };
    const onSelectionChange = ( { value } ) => {
        setSelectedIp( value );
        setHostInvalid( false );
        onInputChange( { value: value?.code ?? ( value || '' ) }, 'ip' );
    };
    const onInputChange = async ( e, name ) => {
        let value = ( e.value || '' );
        if ( name === 'port' ) {
            setPortInvalid( false );
        }
        let newSettings = {...settings, ...{[name]: name === 'port' ? parseInt( value || '0' ) : value}};
        await saveSettings( newSettings );
        setSettings( newSettings );
        setShowSuccess( true );
        closeSuccess();
    };
    return (
        <div className="nav flex justify-content-center">
            <Menu model={items} popup ref={menuRight} id="popup_menu_right" popupAlignment="right" />
            <Button className="btn" icon="pi pi-ellipsis-h" size="small" severity="secondary" rounded  onClick={( event ) => menuRight.current.toggle( event )} aria-controls="popup_menu_right" aria-haspopup />
            <Dialog visible={aboutDialog} resizable={false} draggable={false} style={{ width: '32rem' }} breakpoints={{ '960px': '75vw', '641px': '90vw' }} header="About" modal className="p-fluid" onHide={() => { setAboutDialog( false ); }}>
                <div>PortForwarder is a lightweight, single-purpose utility for forwarding local ports to destination ports without installation.</div>
                <div className="flex justify-content-center align-items-center gap-2" style={{marginTop: '32px'}}>
                    <a href="https://github.com/kerwin612/PortForwarder" target="_blank" rel="noreferrer">Feedback</a>
                    <span style={{border: '0.05rem solid lightgrey', height: '0.8rem'}}/>
                    <a href="https://github.com/kerwin612" target="_blank" rel="noreferrer">©kerwin612</a>
                </div>
            </Dialog>
            <Dialog visible={settingsDialog} resizable={false} draggable={false} style={{ width: '38rem' }} breakpoints={{ '1140px': '75vw', '641px': '90vw' }} header="Settings" modal className="p-fluid" onHide={() => { setSettingsDialog( false ); }}>
                {
                    showSuccess
                    ? (
                        <Message severity="success" text="Changes have been saved" />
                    )
                    : (
                        <Message severity="warn" text="Restart required to apply settings changes" />
                    )
                }
                <div className="formgrid grid">
                    <div className="field col-9">
                        <label htmlFor="workspace" className="font-bold">
                            Workspace
                        </label>
                        <span className="tooltip-workspace"><InputText id="workspace" value={workspace} disabled /></span>
                        <Tooltip target=".tooltip-workspace" position="mouse" autoHide={false}>
                            <small>Changes only via startup args:</small><Tag value="PortForwarder -ws=new_workspace" />
                        </Tooltip>
                    </div>
                    <div className="field col-3">
                        <label htmlFor="explorer" className="font-bold">
                            &nbsp;
                        </label>
                        <Button id="explorer" label="Explorer" severity="info" onClick={explorerWS} />
                    </div>
                    <div className="field col-9">
                        <label htmlFor="ip" className="font-bold">
                            IP
                        </label>
                        <Dropdown id="ip" value={selectedIp} onChange={onSelectionChange} options={IPS} optionLabel="name"
                            editable placeholder="Select or type a valid IP address" />
                        {hostInvalid && <small className="p-error">The IP address is invalid.</small>}
                    </div>
                    <div className="field col-3">
                        <label htmlFor="port" className="font-bold">
                            Port
                        </label>
                        <InputNumber id="port" value={settings.port} min={0} useGrouping={false} onValueChange={( e ) => onInputChange( e, 'port' )} />
                        {portInvalid && <small className="p-error">The port is invalid.</small>}
                    </div>
                </div>
            </Dialog>
            <Dialog
                resizable={false}
                draggable={false}
                visible={byeDialog}
                content={() => (
                    <Message
                        severity="success"
                        style={{
                            background: 'rgba(228, 248, 240)',
                            color: '#1ea97c'
                        }}
                        content={() => (
                            <>Bye-bye, buddy!</>
                        )}
                    />
                )}
            />
        </div>
    );
}
