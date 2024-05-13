package lib

import (
    "sync"
    "io/ioutil"
    "encoding/json"
)

type ForwardInfo struct {
    Id          string   `json:"-"`
    SourceAddr  string   `json:"source_addr,omitempty"`
    SourcePort  int      `json:"source_port"`
    TargetAddr  string   `json:"target_addr,omitempty"`
    TargetPort  int      `json:"target_port"`
    Protocol    string   `json:"protocol"`
    Description string   `json:"description,omitempty"`
}

type Json struct {
    Forwards map[string]*ForwardInfo `json:"forwards"`
}

type Configuration struct {
    Forwards sync.Map
}

func NewConfiguration() *Configuration {
    return &Configuration{}
}

func NewConfigurationFromFile(filePath string) (*Configuration, error) {
    c := NewConfiguration()
    return c, c.LoadFromFile(filePath)
}

func (c *Configuration) LoadFromFile(filePath string) error {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }

    j := Json{
        Forwards: make(map[string]*ForwardInfo),
    }
    if err := json.Unmarshal(data, &j); err != nil {
        return err
    }

    for id, fwd := range j.Forwards {
        fwd.Id = id
        c.AddOrUpdateForward(fwd);
    }

    return nil
}

func (c *Configuration) SaveToFile(filePath string) error {
    j := Json{
        Forwards: make(map[string]*ForwardInfo),
    }

    c.Forwards.Range(func(key, value any) bool {
        if id, ok := key.(string); ok {
            if fwd, ok := value.(*ForwardInfo); ok {
                j.Forwards[id] = fwd
            }
        }
        return true
    })

    data, err := json.MarshalIndent(&j, "", "\t")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(filePath, data, 0644)

}

func (c *Configuration) DeleteForward(id string) {
    c.Forwards.Delete(id)
}

func (c *Configuration) AddOrUpdateForward(fwd *ForwardInfo) {
    c.DeleteForward(fwd.Id)
    c.Forwards.Store(fwd.Id, fwd)
}

func (c *Configuration) GetForward(id string) *ForwardInfo {
    if value, loaded := c.Forwards.Load(id); loaded {
        if fwd, ok := value.(*ForwardInfo); ok {
            return fwd
        }
        return nil
    }
    return nil
}

func (c *Configuration) GetForwards() sync.Map {
    return c.Forwards
}
