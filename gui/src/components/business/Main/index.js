
import React, {
    useRef,
    useState,
    useEffect
} from 'react';

import {
    Badge,
    Toast,
    Dialog,
    Button,
    Tooltip,
    Dropdown,
    InputText,
    InputIcon,
    IconField,
    InputNumber,
    ContextMenu,
    InputTextarea,
    DataListTable,
} from 'components/base';

import * as services from 'services';

import './index.css';
import Menu from '../Menu';
import Brand from '../Brand';

export default function Main() {

    const emptyForward = {
        id: null,
        source_addr: '',
        source_port: 0,
        target_addr: '',
        target_port: 0,
        protocol: '',
        description: ''
    };

    const cm = useRef( null );
    const toast = useRef( null );
    const [errors, setErrors] = useState( null );
    const [forward, setForward] = useState( {} );
    const [submitted, setSubmitted] = useState( false );
    const [selectedIp, setSelectedIp] = useState( null );
    const [forwardList, setForwardList] = useState( [] );
    const [globalFilter, setGlobalFilter] = useState( null );
    const [totalForwards, setTotalForwards] = useState( [] );
    const [forwardDialog, setForwardDialog] = useState( false );
    const [selectedForward, setSelectedForward] = useState( null );
    const [selectedForwards, setSelectedForwards] = useState( [] );
    const [selectedProtocol, setSelectedProtocol] = useState( null );
    const [deleteForwardDialog, setDeleteForwardDialog] = useState( false );
    const [deleteForwardsDialog, setDeleteForwardsDialog] = useState( false );

    useEffect( () => {
        loadForwardList();
    }, [] );

    useEffect( () => {
        setForward( prevForward => ( { ...prevForward, source_addr: selectedIp?.code ?? ( selectedIp || '' ) } ) );
    }, [selectedIp] );

    useEffect( () => {
        setForward( prevForward => ( { ...prevForward, protocol: selectedProtocol?.code ?? '' } ) );
    }, [selectedProtocol] );

    useEffect( () => {
        setForwardList( ( !globalFilter ) ? totalForwards : totalForwards.filter( i => {
            return (
                `${i.source_addr ?? ''}`.indexOf( globalFilter ) !== -1
                || `${i.source_port ?? ''}`.indexOf( globalFilter ) !== -1
                || `${i.target_addr ?? ''}`.indexOf( globalFilter ) !== -1
                || `${i.target_port ?? ''}`.indexOf( globalFilter ) !== -1
                || `${i.description ?? ''}`.indexOf( globalFilter ) !== -1
            );
        } ) );
    }, [totalForwards, globalFilter] );

    const empty = ( array ) => {
        return ( array == null || array.length === 0 ) ? null : array;
    };

    const generateId = () => `${( String.fromCharCode( 97 + Math.floor( Math.random() * 26 ) ) + Math.random().toString( 36 ).slice( 2, 6 ) + Date.now().toString( 36 ) + Math.random().toString( 36 ).slice( 2, 8 ) + Date.now().toString( 36 ) + Math.random().toString( 36 ).slice( 2, 6 ) ) + String.fromCharCode( 97 + Math.floor( Math.random() * 26 ) )}`;

    const telnetTarget = async ( {target_addr, target_port, protocol} ) => {
        let rst = await services.telnet( {protocol, ip: target_addr, port: target_port} );
        if ( rst?.data === 'success' ) {
            showToast( { detail: `telnet [${protocol}/${target_addr}/${target_port}] success` } );
        } else {
            showToast( { severity: 'error', summary: 'Failure', detail: `telnet [${protocol}/${target_addr}/${target_port}] => ${rst?.data}`, life: 3000 } );
        }
    };

    const loadForwardList = async () => {
        setTotalForwards( await services.getForwardList() );
        // setTimeout(() => {
        //     forwardList.forEach(async (item, index) => {
        //         refreshRowTargetStatus(index, await services.telnet({protocol: item.protocol, ip: item.target_addr, port: item.target_port}));
        //     });
        // }, 1500);
    };

    const resetForwardFrom = () => {
        setErrors( null );
        setSelectedIp( services.IPS[0] );
        setSelectedProtocol( services.PROTOCOLS[0] );
        setForward( { ...emptyForward, id: generateId(), source_addr: services.IPS[0].code, protocol: services.PROTOCOLS[0].code } );
    };

    const openNewForwardForm = () => {
        resetForwardFrom();
        setSubmitted( false );
        setForwardDialog( true );
    };

    const openEditForwardForm = ( forward ) => {
        setForward( { ...forward } );
        let ips = services.IPS.filter( i => i.code === ( forward.source_addr ?? '' ) );
        setSelectedIp( ips.length > 0 ? ips[0] : forward.source_addr );
        setSelectedProtocol( {code: forward.protocol, name: forward.protocol.toUpperCase()} );
        setForwardDialog( true );
    };

    // const toggleRowHover = ( index, status ) => {
    //     if ( !forwardList || !forwardList.length || forwardList.length < index )  return;
    //     let newForwardList = [...forwardList];
    //     newForwardList[index].onMouseEnter = status;
    //     setForwardList( newForwardList );
    // };

    // const refreshRowTargetStatus = (index, status) => {
    //     if (!forwardList || !forwardList.length || forwardList.length < index)  return;
    //     let newForwardList = [...forwardList];
    //     newForwardList[index].targetStatus = status;
    //     setForwardList(newForwardList);
    // };

    const confirmDeleteSelected = () => {
        setDeleteForwardsDialog( true );
    };

    const confirmDeleteForward = ( forward ) => {
        setForward( forward );
        setDeleteForwardDialog( true );
    };

    const showToast = ( info ) => {
        toast.current.show( { severity: 'success', summary: 'Successful', life: 1000, ...info } );
    };

    const startForwards = async ( ids ) => {
        let data = ( await services.startForwards( empty( ids ) || ( selectedForwards ?? [] ).map( i => i.id ) ) )?.data;
        loadForwardList();
        if ( empty( ids ) && data[ids[0]] !== 'success' ) {
            showToast( { severity: 'error', summary: 'Failure', detail: `start faild => ${data[ids[0]]}`, life: 3000 } );
        } else {
            showToast( { detail: `forward${empty( ids ) ? '' : 's'} started` } );
        }
    };

    const stopForwards = async ( ids ) => {
        await services.stopForwards( empty( ids ) || ( selectedForwards ?? [] ).map( i => i.id ) );
        loadForwardList();
        showToast( { detail: `forward${empty( ids ) ? '' : 's'} stoped` } );
    };

    const restartForwards = async ( ids ) => {
        let data = ( await services.restartForwards( empty( ids ) || ( selectedForwards ?? [] ).map( i => i.id ) ) )?.data;
        loadForwardList();
        if ( empty( ids ) && data[ids[0]] !== 'success' ) {
            showToast( { severity: 'error', summary: 'Failure', detail: `restart faild => ${data[ids[0]]}`, life: 3000 } );
        } else {
            showToast( { detail: `forward${empty( ids ) ? '' : 's'} restarted` } );
        }
    };

    const saveForward = async ( isStart ) => {
        if ( submitted )  return;
        setSubmitted( true );
        setErrors( null );
        let {id, source_addr, source_port, target_addr, target_port, protocol, description} = forward;
        source_addr = ( source_addr ?? '' ).trim();
        if ( source_port < 1 ) {
            setErrors( {source_port: 'Local Port is required.'} );
        } else if ( !( target_addr = target_addr.trim() ) ) {
            setErrors( {target_addr: 'Target Address is required.'} );
        } else if ( target_port < 1 ) {
            setErrors( {target_port: 'Target Port is required.'} );
        } else if ( !( protocol = protocol.trim() ) ) {
            setErrors( {protocol: 'Target Protocol is required.'} );
        } else {
            let rst = await ( services[( isStart ? 'saveAndStartForward' : 'saveForward' )] )( {
                id, source_addr, source_port, target_addr, target_port, protocol, description
            } );
            if ( rst.status === 200 ) {
                showToast( { detail: 'forward saved' } );
                setForwardDialog( false );
                loadForwardList();
            } else {
                showToast( { severity: 'error', summary: 'Failure', detail: rst.data } );
            }
        }
        setSubmitted( false );
    };

    const hideForwardDialog = () => {
        setSubmitted( false );
        setForwardDialog( false );
    };

    const deleteForward = async () => {
        await services.deleteForwards( [forward.id] );
        loadForwardList();
        showToast( { detail: 'forward deleted' } );
        setDeleteForwardDialog( false );
    };

    const hideDeleteForwardDialog = () => {
        setDeleteForwardDialog( false );
    };

    const deleteForwards = async () => {
        await services.deleteForwards( ( selectedForwards ?? [] ).map( i => i.id ) );
        loadForwardList();
        showToast( { detail: 'forwards deleted' } );
        setDeleteForwardsDialog( false );
        setSelectedForwards( [] );
    };

    const hideDeleteForwardsDialog = () => {
        setDeleteForwardsDialog( false );
    };

    const onInputChange = ( e, name ) => {
        const val = ( e.target && e.target.value ) || '';
        let _forward = { ...forward };

        _forward[`${name}`] = val;

        setForward( _forward );
    };

    const onInputNumberChange = ( e, name ) => {
        const val = e.value || 0;
        let _forward = { ...forward };

        _forward[`${name}`] = val;

        setForward( _forward );
    };

    const listTableHeader = () => {
        let tooltipOptions = { showDelay: 800, showOnDisabled: true, position: 'bottom' };
        return (
            <div className="flex flex-wrap gap-2 align-items-center justify-content-between">
                <Brand />
                <div className="flex gap-2" style={{padding: 10, maxHeight: 60}}>
                    <div className="flex flex-wrap gap-2">
                        <Button label="New" severity="success" onClick={openNewForwardForm} tooltip="add a record" tooltipOptions={tooltipOptions} size="small" text />
                        <Button label="Delete" severity="danger" onClick={confirmDeleteSelected} disabled={!selectedForwards || !selectedForwards.length} tooltip="delete the selected records" tooltipOptions={tooltipOptions} size="small" text />
                        <Button label="Restart" severity="warning" onClick={() => restartForwards( [] )} disabled={!forwardList || !forwardList.length} tooltip="restart selected or all records" tooltipOptions={tooltipOptions} size="small" text />
                        <Button label="Start" severity="info" onClick={() => startForwards( [] )} disabled={!forwardList || !forwardList.length} tooltip="start selected or all records" tooltipOptions={tooltipOptions} size="small" text />
                        <Button label="Stop" severity="secondary" onClick={() => stopForwards( [] )} disabled={!forwardList || !forwardList.length} tooltip="stop selected or all records" tooltipOptions={tooltipOptions} size="small" text />
                    </div>
                    <Menu />
                </div>
            </div>
        );
    };

    const actionsHeader = () => {
        return (
            <span className="p-input-icon-left">
                <IconField iconPosition="left">
                    <InputIcon className="pi pi-search" />
                    <InputText onInput={( e ) => setGlobalFilter( e.target.value )} placeholder="Search..." />
                </IconField>
            </span>
        );
    };

    const dataIndexBodyTemplate = ( rowData, options ) => {
        return (
            <Badge value={options.rowIndex + 1} severity={rowData.status === 2 ? 'success' : 'secondary'} />
        );
    };

    const dataSourceAddrBodyTemplate = ( rowData ) => {
        let value = rowData.source_addr;
        let ips = services.IPS.filter( i => i.code === ( value ?? '' ) );
        return (
            <>{( ips.length > 0 ? ips[0] : {code: value} ).code || '0.0.0.0'}</>
        );
    };

    const dataDescriptionBodyTemplate = ( rowData ) => {
        let id = generateId();
        let cn = `${id}_description`;
        return (
            <>
                <Tooltip target={`.${cn}`} position="mouse" showDelay={300} autoHide={false} style={{maxWidth: '30vw'}}/>
                <div className={cn} data-pr-tooltip={rowData.description} style={{maxWidth: '20vw', overflow: 'hidden', whiteSpace: 'nowrap', textOverflow: 'ellipsis'}}>
                    {rowData.description}
                </div>
            </>
        );
    };

    const cmModel = [
        { label: 'Restart the record', command: () => restartForwards( [selectedForward.id] ) },
        { label: 'Start the record', command: () => startForwards( [selectedForward.id] ) },
        { label: 'Stop the record', command: () => stopForwards( [selectedForward.id] ) },
        { label: 'Telnet [Target Port] of the record', command: () => telnetTarget( selectedForward ) },
    ];

    const dataActionBodyTemplate = ( rowData ) => {
        let tooltipOptions = { showDelay: 800, showOnDisabled: true, position: 'bottom' };
        return (
            <div className='flex flex-wrap justify-content-end gap-2'>
                <Button icon="pi pi-pencil" className="action-button" tooltip="edit the record" tooltipOptions={tooltipOptions} rounded outlined onClick={() => openEditForwardForm( rowData )} />
                <Button icon="pi pi-trash" className="action-button" tooltip="delete the record" tooltipOptions={tooltipOptions} rounded outlined severity="danger" onClick={() => confirmDeleteForward( rowData )} />
                <Button icon="pi pi-bars" className="action-button" tooltip="more actions" tooltipOptions={tooltipOptions} rounded outlined severity="warning" onClick={( e ) => {
                    setSelectedForward( rowData );
                    cm.current.show( e );
                }} />
            </div>
        );
    };

    const formErrorTemplate = ( field ) => {
        let error = ( errors || {} )[field];
        return (
            error && <small className="p-error">{error}</small>
        );
    };

    const forwardDialogFooter = () => {
        return (
            <React.Fragment>
                <Button label="Cancel" icon="pi pi-times" outlined onClick={hideForwardDialog} />
                <Button label="Save" icon="pi pi-check" onClick={() => saveForward( false )} disabled={submitted} />
                <Button label="Save and Start" icon="pi pi-check" onClick={() => saveForward( true )} disabled={submitted} />
            </React.Fragment>
        );
    };

    const deleteForwardDialogFooter = () => {
        return (
            <React.Fragment>
                <Button label="No" icon="pi pi-times" outlined onClick={hideDeleteForwardDialog} />
                <Button label="Yes" icon="pi pi-check" severity="danger" onClick={deleteForward} />
            </React.Fragment>
        );
    };

    const deleteForwardsDialogFooter = () => {
        return (
            <React.Fragment>
                <Button label="No" icon="pi pi-times" outlined onClick={hideDeleteForwardsDialog} />
                <Button label="Yes" icon="pi pi-check" severity="danger" onClick={deleteForwards} />
            </React.Fragment>
        );
    };

    return (
        <div className="main">
            <Toast ref={toast} />
            <div className="card">
                <ContextMenu model={cmModel} ref={cm} onHide={() => setSelectedForward( null )} />
                <DataListTable
                    table={{
                        dataKey: 'id',
                        value: forwardList,
                        selection: selectedForwards,
                        selectionMode: 'checkbox',
                        onSelectionChange: ( e ) => setSelectedForwards( e.value ),
                        header: listTableHeader,
                        scrollable: true,
                        scrollHeight: 'flex',
                        removableSort: true,
                        sortMode: 'multiple',
                        multiSortMeta: [{field: 'source_port', order: 1}],
                        className: 'list',
                        emptyMessage: 'No results found',
                        paginator: true,
                        rows: 15,
                        alwaysShowPaginator: false,
                        rowsPerPageOptions: [5, 15, 25, 50],
                        onContextMenu: ( e ) => cm.current.show( e.originalEvent ),
                        contextMenuSelection: selectedForward,
                        onContextMenuSelectionChange: ( e ) => setSelectedForward( e.value )
                        // onRowMouseEnter: ( {index} ) => {
                        //     toggleRowHover( index, true );
                        // },
                        // onRowMouseLeave: ( {index} ) => {
                        //     toggleRowHover( index, false );
                        // },
                    }}
                    columns={[
                        { key: 'checkbox', field: 'id', selectionMode: 'multiple', exportable: false },
                        { key: 'index', headerStyle: { width: '3rem' }, header: '#', body: dataIndexBodyTemplate, style: { width: '4%' } },
                        { key: 'source_addr', field: 'source_addr', sortable: true, header: 'Local Address', body: dataSourceAddrBodyTemplate, style: { width: '14%' } },
                        { key: 'source_port', field: 'source_port', sortable: true, header: 'Local Port', style: { width: '9%' } },
                        { key: 'target_addr', field: 'target_addr', sortable: true, header: 'Target Address', style: { width: '14%' } },
                        { key: 'target_port', field: 'target_port', sortable: true, header: 'Target Port', style: { width: '9%' } },
                        { key: 'protocol', field: 'protocol', header: 'Target Protocol', style: { width: '10%' } },
                        { key: 'description', field: 'description', header: 'Description', body: dataDescriptionBodyTemplate, style: { width: '20%' } },
                        { key: 'actions', header: actionsHeader, headerClassName: 'actions-header', body: dataActionBodyTemplate, style: { width: '20%' } },
                    ]}
                />
            </div>

            <Dialog visible={forwardDialog} style={{ width: '38rem' }} breakpoints={{ '1140px': '75vw', '641px': '90vw' }} header="Forward Details" modal className="p-fluid" footer={forwardDialogFooter} onHide={hideForwardDialog}>

                <div className="formgrid grid">
                    <div className="field col-8">
                        <label htmlFor="source_addr" className="font-bold">
                            Local Address
                        </label>
                        <Dropdown value={selectedIp} onChange={( e ) => setSelectedIp( e.value )} options={services.IPS} optionLabel="name"
                            editable placeholder="Select or type a valid IP address" />
                        {formErrorTemplate( 'source_addr' )}
                    </div>
                    <div className="field col-4">
                        <label htmlFor="source_port" className="font-bold">
                            Local Port
                        </label>
                        <InputNumber id="source_port" min={0} useGrouping={false} value={forward.source_port} onValueChange={( e ) => onInputNumberChange( e, 'source_port' )} placeholder="Type a valid Port" />
                        {formErrorTemplate( 'source_port' )}
                    </div>
                </div>

                <div className="formgrid grid">
                    <div className="field col-8">
                        <label htmlFor="target_addr" className="font-bold">
                            Target Address
                        </label>
                        <InputText id="target_addr" value={forward.target_addr} onChange={( e ) => onInputChange( e, 'target_addr' )} placeholder="Type a valid IP address" />
                        {formErrorTemplate( 'target_addr' )}
                    </div>
                    <div className="field col-4">
                        <label htmlFor="target_port" className="font-bold">
                            Target Port
                        </label>
                        <InputNumber id="target_port" min={0} useGrouping={false} value={forward.target_port} onValueChange={( e ) => onInputNumberChange( e, 'target_port' )} placeholder="Type a valid Port" />
                        {formErrorTemplate( 'target_port' )}
                    </div>
                </div>

                <div className="field">
                    <label htmlFor="protocol" className="font-bold">
                        Target Protocol
                    </label>
                    <Dropdown id="protocol" value={selectedProtocol} onChange={( e ) => setSelectedProtocol( e.value )} options={services.PROTOCOLS} optionLabel="name"
                        placeholder="Select a Protocol" />
                    {formErrorTemplate( 'protocol' )}
                </div>

                <div className="field">
                    <label htmlFor="description" className="font-bold">
                        Description
                    </label>
                    <InputTextarea id="description" autoResize value={forward.description} onChange={( e ) => onInputChange( e, 'description' )} rows={3} cols={20} placeholder="Type a Description" />
                    {formErrorTemplate( 'description' )}
                </div>

            </Dialog>

            <Dialog visible={deleteForwardDialog} style={{ width: '32rem' }} breakpoints={{ '960px': '75vw', '641px': '90vw' }} header="Confirm" modal footer={deleteForwardDialogFooter} onHide={hideDeleteForwardDialog}>
                <div className="flex align-items-center justify-content-center confirmation-content">
                    <i className="pi pi-exclamation-triangle mr-3" style={{ fontSize: '2rem' }} />
                    {forward && (
                        <span>
                            Are you sure you want to delete [<b>{forward.id}</b>] ?
                        </span>
                    )}
                </div>
            </Dialog>

            <Dialog visible={deleteForwardsDialog} style={{ width: '32rem' }} breakpoints={{ '960px': '75vw', '641px': '90vw' }} header="Confirm" modal footer={deleteForwardsDialogFooter} onHide={hideDeleteForwardsDialog}>
                <div className="flex align-items-center justify-content-center confirmation-content">
                    <i className="pi pi-exclamation-triangle mr-3" style={{ fontSize: '2rem' }} />
                    {forward && <span>Are you sure you want to delete the selected forwards?</span>}
                </div>
            </Dialog>
        </div>
    );
}
