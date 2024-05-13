package lib

import (
    "os"
    "testing"
)

func Test_NewConfiguration(t *testing.T) {

    c1 := NewConfiguration()

    c1.AddOrUpdateForward(&ForwardInfo{
        Id:          "test",
        SourcePort:  8081,
        TargetPort:  8089,
    })

    tmpfile, _ := os.CreateTemp("", "temp")
    defer tmpfile.Close()

    tempFilePath := tmpfile.Name()
    t.Logf("%+v\n", tempFilePath)

    c1.SaveToFile(tempFilePath)

    c2 := NewConfiguration()
    c2.LoadFromFile(tempFilePath)

    entry := c2.GetForward("test")

    if entry.SourcePort != 8081 || entry.TargetPort != 8089 {
        t.Errorf("Config entries do not match: %+v", entry)
    }
}

