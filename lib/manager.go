package lib

import (
    "net"
    "sync"
    "time"
    "errors"
    "log"
    "github.com/kerwin612/port-forwarder/lib/core"
)

type ForwardStatus struct {
    IsRunning      bool
    StartTime      time.Time
    LastError      error
    ConnectionCount int
}

type ForwardManager struct {
    Listeners sync.Map  // id -> net.Listener
    Status    sync.Map  // id -> *ForwardStatus
}

func NewForwardManager() *ForwardManager {
    return &ForwardManager{
        Listeners: sync.Map{},
        Status:    sync.Map{},
    }
}

func (m *ForwardManager) Exist(id string) *net.Listener {
    if value, loaded := m.Listeners.Load(id); loaded {
        if ln, ok := value.(net.Listener); ok {
            return &ln
        }
    }
    return nil
}

func (m *ForwardManager) GetStatus(id string) *ForwardStatus {
    if value, loaded := m.Status.Load(id); loaded {
        if status, ok := value.(*ForwardStatus); ok {
            return status
        }
    }
    return nil
}

func (m *ForwardManager) Close(id string) error {
    if value, loaded := m.Listeners.Load(id); loaded {
        if ln, ok := value.(net.Listener); ok {
            log.Printf("Closing forward with ID: %s", id)
            err := ln.Close()
            m.Listeners.Delete(id)
            
            if status := m.GetStatus(id); status != nil {
                status.IsRunning = false
                if err != nil {
                    status.LastError = err
                }
            }
            
            if err != nil {
                log.Printf("Error closing forward %s: %v", id, err)
                return err
            }
            log.Printf("Successfully closed forward with ID: %s", id)
        }
    } else {
        log.Printf("Attempted to close non-existent forward with ID: %s", id)
        return errors.New("forward not found")
    }
    return nil
}

func (m *ForwardManager) Open(fwd *ForwardInfo) error {
    if fwd == nil {
        return errors.New("invalid forward configuration: nil value")
    }
    
    if fwd.SourcePort <= 0 || fwd.TargetPort <= 0 {
        return errors.New("invalid port configuration")
    }
    
    log.Printf("Opening forward %s: %s:%d -> %s:%d [%s]", 
        fwd.Id, fwd.SourceAddr, fwd.SourcePort, fwd.TargetAddr, fwd.TargetPort, fwd.Protocol)
    
    status := &ForwardStatus{
        IsRunning:      false,
        StartTime:      time.Now(),
        ConnectionCount: 0,
    }
    
    ln, err := core.Forward(core.Info{
        SourceAddr:  fwd.SourceAddr,
        SourcePort:  fwd.SourcePort,
        TargetAddr:  fwd.TargetAddr,
        TargetPort:  fwd.TargetPort,
        Protocol:    fwd.Protocol,
    })
    
    if err != nil {
        log.Printf("Failed to open forward %s: %v", fwd.Id, err)
        status.LastError = err
        m.Status.Store(fwd.Id, status)
        return err
    }
    
    // Verify that the listener is actually working
    if ln == nil {
        err := errors.New("listener is nil after successful creation")
        log.Printf("Failed to open forward %s: %v", fwd.Id, err)
        status.LastError = err
        m.Status.Store(fwd.Id, status)
        return err
    }
    
    // Test if target connection is reachable
    canConnect, err := core.Telnet(fwd.Protocol, fwd.TargetAddr, fwd.TargetPort, 3000)
    if !canConnect {
        if err != nil {
            log.Printf("Warning: Target %s:%d is not reachable: %v", 
                       fwd.TargetAddr, fwd.TargetPort, err)
            // Only log warning without returning error, as target might be temporarily unavailable
            status.LastError = err
        } else {
            log.Printf("Warning: Target %s:%d is not reachable for unknown reason", 
                      fwd.TargetAddr, fwd.TargetPort)
        }
    }
    
    m.Listeners.Store(fwd.Id, ln)
    status.IsRunning = true
    m.Status.Store(fwd.Id, status)
    
    log.Printf("Successfully opened forward with ID: %s", fwd.Id)
    return nil
}

func (m *ForwardManager) ReOpen(fwd *ForwardInfo) error {
    if err := m.Close(fwd.Id); err != nil && !errors.Is(err, errors.New("forward not found")) {
        return err
    }
    return m.Open(fwd)
}
