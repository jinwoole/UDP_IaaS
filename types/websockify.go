// websockify.go
package types

import (
	"fmt"
	"os/exec"
)

func NewWebsockifyManager() *WebsockifyManager {
    return &WebsockifyManager{
        processes: make(map[int]*exec.Cmd),
    }
}

func (wm *WebsockifyManager) Start(vncPort int) (int, error) {
    wm.mu.Lock()
    defer wm.mu.Unlock()

    // Websocket 포트는 VNC 포트 + 1000
    wsPort := vncPort + 1000

    // 이미 실행 중인 프로세스가 있다면 종료
    if cmd, exists := wm.processes[wsPort]; exists {
        if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
            return wsPort, nil // 이미 실행 중
        }
    }

    // websockify 프로세스 시작
    cmd := exec.Command("websockify",
        "--web", "/usr/share/novnc/",
        fmt.Sprintf("%d", wsPort),
        fmt.Sprintf("localhost:%d", vncPort),
    )

    if err := cmd.Start(); err != nil {
        return 0, fmt.Errorf("failed to start websockify: %w", err)
    }

    wm.processes[wsPort] = cmd

    // 프로세스 모니터링
    go func() {
        cmd.Wait()
        wm.mu.Lock()
        delete(wm.processes, wsPort)
        wm.mu.Unlock()
    }()

    return wsPort, nil
}