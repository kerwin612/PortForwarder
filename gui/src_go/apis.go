package src_go

import (
    "fmt"
    "net"
    "time"
    "errors"
    "os/exec"
    "strconv"
    "runtime"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "github.com/rakyll/statik/fs"
    "github.com/getlantern/systray"
    l "github.com/kerwin612/hybrid-launcher"
    "github.com/kerwin612/port-forwarder/lib"
    "github.com/kerwin612/port-forwarder/lib/core"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

var fp string
var rootdir string
var sm *SettingsManager
var launcher *l.Launcher
var fc *lib.Configuration
var fm *lib.ForwardManager

func Startup(homedir string, autoopen bool) {

    rootdir = homedir

    var err error
    sm, err = NewSettingsManager(filepath.Join(rootdir, "cfg"))
    if err != nil {
        logger.Fatalln(err)
    }

    si, err := sm.Load()
    if err != nil {
        logger.Println("[Error] cfg invalid:", err)
    }

    statikFS, err := fs.New()
    if err != nil {
        logger.Fatalln(err)
    }

    fp = filepath.Join(rootdir, "fwd")

    c, e := lib.NewConfigurationFromFile(fp)
    if e != nil {
        logger.Println("[Error] load fwd faild:", e)
    } else {
        logger.Println("forwards loaded.")
    }

    fc = c
    fm = lib.NewForwardManager()

    routes := map[string]map[string]HandlerFunc{
        "/exit": map[string]HandlerFunc{
            "GET": func(w http.ResponseWriter, r *http.Request) {
                _writeString(w, "success")
                logger.Printf("quit by gui.\n")
                go func(){
                    time.Sleep(time.Millisecond * 1500)
                    launcher.Exit()
                }()
            },
        },
        "/v1/ws/explorer": map[string]HandlerFunc{
            "GET": func(w http.ResponseWriter, r *http.Request) {
                err := _openFolderInDefaultBrowser(rootdir)
                if err != nil {
                    _writeError(w, err)
                } else {
                    _writeString(w, "success")
                }
            },
        },
        "/v1/settings": map[string]HandlerFunc{
            "GET": _settings_get,
            "POST": _settings_save,
        },
        "/v1/telnet/{protocol}/{ip}/{port}": map[string]HandlerFunc{
            "GET": _telnet,
        },
        "/v1/forward/protocol": map[string]HandlerFunc{
            "GET": _forward_protocol,
        },
        "/v1/forward/ip/list": map[string]HandlerFunc{
            "GET": _forward_ip_list,
        },
        "/v1/forward/list": map[string]HandlerFunc{
            "GET": _forward_list,
        },
        "/v1/forward/stop": map[string]HandlerFunc{
            "POST": _forward_stop,
        },
        "/v1/forward/delete": map[string]HandlerFunc{
            "POST": _forward_del,
        },
        "/v1/forward/start": map[string]HandlerFunc{
            "POST": _forward_start,
        },
        "/v1/forward/restart": map[string]HandlerFunc{
            "POST": _forward_restart,
        },
        "/v1/forward/save/{id}": map[string]HandlerFunc{
            "POST": func(w http.ResponseWriter, r *http.Request) {
                _forward_save_and_start(w, r, false)
            },
        },
        "/v1/forward/save_and_start/{id}": map[string]HandlerFunc{
            "POST": func(w http.ResponseWriter, r *http.Request) {
                _forward_save_and_start(w, r, true)
            },
        },
    }

    for path, handlerMap := range routes {
        for method, handler := range handlerMap {
            http.Handle(fmt.Sprintf("%s %s", method, path), _commonHandlerMiddleware(&handler))
        }
        http.Handle(fmt.Sprintf("%s %s", http.MethodOptions, path), _commonHandlerMiddleware(nil))
    }


    cfg, err := l.DefaultConfig()
    if err != nil {
        logger.Fatalln(err)
    }

    cfg.Pid = filepath.Join(rootdir, "pid")
    cfg.Ip = si.Ip
    cfg.Port = si.Port
    cfg.Icon = IconData
    cfg.Title = "PortForwarder"
    cfg.Tooltip = "PortForwarder"
    cfg.RootHandler = http.StripPrefix("/", http.FileServer(statikFS))
    cfg.TrayOnReady = func() {

        mStart := systray.AddMenuItem("Start", "Start All Forwards")
        mStop := systray.AddMenuItem("Stop", "Stop All Forwards")
        mRestart := systray.AddMenuItem("Restart", "Restart All Forwards")

        systray.AddSeparator()

        mShow := systray.AddMenuItem("Open", "Open PortForwarder")
        mQuit := systray.AddMenuItem("Quit", "Quit PortForwarder")

        go func() {
            for {
                select {
                    case <-mRestart.ClickedCh:
                        _stopForward([]string{})
                        _startForward([]string{})
                    case <-mStart.ClickedCh:
                        _startForward([]string{})
                    case <-mStop.ClickedCh:
                        _stopForward([]string{})
                    case <-mShow.ClickedCh:
                        go launcher.Open()
                    case <-mQuit.ClickedCh:
                        logger.Printf("quit by systray.\n")
                        launcher.Exit()
                        return
                }
            }
        }()

    }

    var existingLauncher *l.ExistingLauncher
    if launcher, existingLauncher, err = l.NewWithConfig(cfg); err != nil {
        if existingLauncher != nil {
            logger.Printf("listener: [%s].\n", existingLauncher.Addr)
            existingLauncher.Open()
            logger.Fatalln(err)
        }
        cfg.Ip = ""
        cfg.Port = 0
        if launcher, existingLauncher, err = l.NewWithConfig(cfg); err != nil {
            if existingLauncher != nil {
                logger.Printf("listener: [%s].\n", existingLauncher.Addr)
                existingLauncher.Open()
            }
            logger.Fatalln(err)
        }
    }
    logger.Printf("listener: [%s].\n", launcher.Addr())

    var se error
    if autoopen {
        se = launcher.StartAndOpen()
    } else {
        se = launcher.Start()
    }
    if se != nil {
        logger.Fatalln(se)
    }

}

func _settings_get(w http.ResponseWriter, r *http.Request) {

    si, err := sm.Load()
    if err != nil {
        _writeJson(w, &[]interface{}{rootdir, fmt.Sprintf("%v", err)})
    } else {
        _writeJson(w, &[]interface{}{rootdir, &si})
    }

}

func _settings_save(w http.ResponseWriter, r *http.Request) {
    var settingsInfo *SettingsInfo
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&settingsInfo); err != nil {
        _writeError(w, err)
    } else {
        si, _ := sm.Load()
        si.Ip = settingsInfo.Ip
        si.Port = settingsInfo.Port
        if err := sm.Save(si); err != nil {
            _writeError(w, err)
        } else {
            logger.Printf("settings saved.")
            _writeString(w, "success")
        }
    }
}

func _telnet(w http.ResponseWriter, r *http.Request) {
    timeout, err := strconv.Atoi(r.URL.Query().Get("timeout"))
    if err != nil || timeout <= 0 {
        timeout = 3000 // Default 3 second timeout
    }
    
    protocol := r.PathValue("protocol")
    if protocol == "" {
        protocol = "tcp" // Default TCP protocol
    }
    
    ip := r.PathValue("ip")
    port, _ := strconv.Atoi(r.PathValue("port"))
    
    // Record test information
    logger.Printf("Telnet test request: %s://%s:%d (timeout: %dms)", protocol, ip, port, timeout)
    
    // Add multiple retries to ensure it's not a temporary network issue
    maxRetries := 3
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        success, err := core.Telnet(protocol, ip, port, timeout)
        if success {
            _writeString(w, "true")
            return
        }
        lastErr = err
        logger.Printf("Telnet attempt %d failed: %v", i+1, err)
        time.Sleep(time.Millisecond * 200) // Wait 200ms between retries
    }
    
    // Log more detailed network information after all attempts fail
    logger.Printf("All telnet attempts failed to %s:%d. Last error: %v", ip, port, lastErr)
    
    // Return error information
    if lastErr != nil {
        _writeString(w, fmt.Sprintf("false|%v", lastErr))
    } else {
        _writeString(w, "false")
    }
}

func _forward_protocol(w http.ResponseWriter, r *http.Request) {
    _writeJson(w, &[]string{"tcp", "udp"})
}

func _forward_ip_list(w http.ResponseWriter, r *http.Request) {

    ips, err := _getListenIPs()
    if err != nil {
        _writeError(w, err)
    } else {
        _writeJson(w, &ips)
        // _writeJson(w, append([]string{"0.0.0.0"}, ips...))
    }

}

func _forward_del(w http.ResponseWriter, r *http.Request) {

    ids := _readIdArray(w, r)
    if ids == nil {
        return
    }
    _stopForward(*ids)
    for _, i := range *ids {
        fc.DeleteForward(i)
        logger.Printf("forward [%s] deleted.", i)
    }
    if err := _saveForwardConfig(); err != nil {
        _writeError(w, err)
    } else {
        _writeString(w, "success")
    }

}

func _forward_list(w http.ResponseWriter, r *http.Request) {

    forwards := fc.GetForwards()
    regularMap := make(map[string]interface{})
    forwards.Range(func(key, value any) bool {
        var status int
        forward := value.(*lib.ForwardInfo)
        if ln := fm.Exist(forward.Id); ln == nil {
            status = 1
        } else {
            status = 2
        }
        regularMap[key.(string)] = []interface{}{forward, status}
        return true
    })
    _writeJson(w, regularMap)

}

func _forward_stop(w http.ResponseWriter, r *http.Request) {

    ids := _readIdArray(w, r)
    if ids == nil {
        return
    }

    _stopForward(*ids)

    _writeString(w, "success")

}

func _forward_start(w http.ResponseWriter, r *http.Request) {

    ids := _readIdArray(w, r)
    if ids == nil {
        return
    }

    _writeJson(w, _startForward(*ids))

}

func _forward_restart(w http.ResponseWriter, r *http.Request) {

    ids := _readIdArray(w, r)
    if ids == nil {
        return
    }

    _stopForward(*ids)

    _writeJson(w, _startForward(*ids))

}

func _forward_save_and_start(w http.ResponseWriter, r *http.Request, start bool) {

    var forward *lib.ForwardInfo
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&forward); err != nil {
        _writeError(w, err)
    } else {
        forward.Id = r.PathValue("id")
        if fwd := _getForwardBySourcePort(forward.SourcePort); fwd != nil && fwd.Id != forward.Id {
            _writeError(w, errors.New("port already exists."))
        } else {
            fc.AddOrUpdateForward(forward)
            if err := _saveForwardConfig(); err != nil {
                _writeError(w, err)
            } else {
                logger.Printf("forward [%s] saved.", forward.Id)
                if start {
                    _writeJson(w, _startForward([]string{forward.Id}))
                } else {
                    _writeString(w, "success")
                }
            }
        }
    }

}

func _saveForwardConfig() error {
    return fc.SaveToFile(fp)
}

func _stopForward(ids []string) {
    forwards := fc.GetForwards()
    forwards.Range(func(key, _ any) bool {
        if id, ok := key.(string); ok {
            if ln := fm.Exist(id); ((len(ids) == 0 || _contains(ids, id)) && ln != nil) {
                if err := fm.Close(id); err != nil {
                    logger.Printf("[Error] forward [%s] stop faild: %v", id, err)
                } else {
                    logger.Printf("forward [%s] stopped.", id)
                }
            }
        }
        return true
    })
}

func _startForward(ids []string) map[string]string {
    rstMap := make(map[string]string)
    forwards := fc.GetForwards()
    forwards.Range(func(key, value any) bool {
        if id, ok := key.(string); ok {
            if fwd, ok := value.(*lib.ForwardInfo); ok {
                if len(ids) == 0 || _contains(ids, id){
                    if ln := fm.Exist(fwd.Id); ln == nil {
                        if err := fm.Open(fwd); err != nil {
                            rstMap[id] = fmt.Sprintf("%v", err)
                            logger.Printf("[Error] forward [%s] start faild: %v", id, err)
                        } else {
                            rstMap[id] = "success"
                            logger.Printf("forward [%s] started.", id)
                        }
                    } else {
                        rstMap[id] = "success"
                        logger.Printf("forward [%s] has already started.", id)
                    }
                }
            }
        }
        return true
    })
    return rstMap
}

func _getForwardBySourcePort(port int) *lib.ForwardInfo {
    var forward *lib.ForwardInfo
    forwards := fc.GetForwards()
    forwards.Range(func(key, value any) bool {
        if fwd, ok := value.(*lib.ForwardInfo); ok {
            if fwd.SourcePort == port {
                forward = fwd
                return false
            }
        }
        return true
    })
    return forward
}

func _contains(array []string, item string) bool {
    for _, i := range array {
        if i == item {
            return true
        }
    }
    return false
}

func _writeError(w http.ResponseWriter, err error) {
    http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
}

func _writeJson(w http.ResponseWriter, object interface{}) {
    w.Header().Set("Content-Type", "application/json")
    data, err := json.Marshal(&object)
    if err != nil {
        _writeError(w, err)
    } else {
        fmt.Fprintf(w, string(data))
    }
}

func _writeString(w http.ResponseWriter, str string) {
    fmt.Fprintf(w, str)
}

func _readIdArray(w http.ResponseWriter, r *http.Request) *[]string {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        _writeError(w, errors.New("Error reading request body"))
        return nil
    }
    defer r.Body.Close()

    var strings []string
    if err := json.Unmarshal(body, &strings); err != nil {
        if len(body) == 0 {
            strings = []string{}
        } else {
            _writeError(w, errors.New("Error parsing JSON body"))
            return nil
        }
    }
    return &strings
}

func _getListenIPs() ([]string, error) {
    var ips []string

    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    for _, iface := range interfaces {
        if iface.Name == "lo" || iface.Flags&net.FlagUp == 0 {
            continue
        }

        addrs, err := iface.Addrs()
        if err != nil {
            return nil, err
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
            }

            if ip == nil || ip.IsLoopback() || ip.To4() == nil {
                continue
            }

            ips = append(ips, ip.String())
        }
    }

    if len(ips) == 0 {
        return nil, errors.New("no listenable IP addresses found")
    }

    return ips, nil
}

func _commonHandlerMiddleware(next *HandlerFunc) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                _writeString(w, "success")
            }).ServeHTTP(w, r)
        } else {
            http.HandlerFunc(*next).ServeHTTP(w, r)
        }
    })
}

func _openFolderInDefaultBrowser(folderPath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
        case "windows":
            cmd = exec.Command("explorer", folderPath)
        case "darwin":
            cmd = exec.Command("open", folderPath)
        case "linux":
            cmd = exec.Command("xdg-open", folderPath)
        default:
            return fmt.Errorf("unsupported platform")
	}

	if err := cmd.Start(); err != nil {
		return err
	}

    return nil
}
