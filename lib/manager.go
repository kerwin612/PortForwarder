package lib

import (
    "net"
    "sync"
    "github.com/kerwin612/port-forwarder/lib/core"
)

type ForwardManager struct {
    Listeners sync.Map
}

func NewForwardManager() *ForwardManager {
    return &ForwardManager{
        Listeners: sync.Map{},
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

func (m *ForwardManager) Close(id string) error {
    if value, loaded := m.Listeners.Load(id); loaded {
        if ln, ok := value.(net.Listener); ok {
            if err := ln.Close(); err != nil {
                return err
            }
            m.Listeners.Delete(id)
        }
    }
    return nil
}

func (m *ForwardManager) Open(fwd *ForwardInfo) error {
    ln, err := core.Forward(core.Info{
        SourceAddr:  fwd.SourceAddr,
        SourcePort:  fwd.SourcePort,
        TargetAddr:  fwd.TargetAddr,
        TargetPort:  fwd.TargetPort,
        Protocol:    fwd.Protocol,
    })
    if err != nil {
        return err
    }
    m.Listeners.Store(fwd.Id, ln)
    return nil
}

func (m *ForwardManager) ReOpen(fwd *ForwardInfo) error {
    if err := m.Close(fwd.Id); err != nil {
        return err
    }
    return m.Open(fwd)
}
