package src_go

import (
    "os"
    "io/ioutil"
    "encoding/json"
)

type SettingsInfo struct {
    Ip      string      `json:"ip"`
    Port    int         `json:"port"`
}

type SettingsManager struct {
    path    string
}

var dftSettingsInfo = SettingsInfo{
    Ip: "",
    Port: 0,
}

func NewSettingsManager(filePath string) (*SettingsManager, error) {
    sm := &SettingsManager{
        path: filePath,
    }
    _, err := os.Open(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            if err := sm.Save(&dftSettingsInfo); err != nil {
                return nil, err
            }
        } else {
            return nil, err
        }
    }
    return sm, nil
}

func (sm *SettingsManager) Load() (*SettingsInfo, error) {
    data, err := ioutil.ReadFile(sm.path)
    if err != nil {
        return &dftSettingsInfo, err
    }
    var s SettingsInfo
    if err := json.Unmarshal(data, &s); err != nil {
        return &dftSettingsInfo, err
    }
    return &s, nil
}

func (sm *SettingsManager) Save(settings *SettingsInfo) error {
    data, err := json.MarshalIndent(&settings, "", "\t")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(sm.path, data, 0644)
}
