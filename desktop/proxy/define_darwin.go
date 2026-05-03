package proxy

import (
	"desktop/internal"
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

// netsh interface ipv4 set address name="PrismTUN" source=static addr=172.19.0.1 mask=255.255.255.255
const TUN_NAME = "utun"
const TUN_MTU = 1500
const TUN_OFFSET = 4

var TUN_IP netip.Prefix
var Routes []netip.Prefix

// 自定义网卡GUID 方便wintun复用
func init() {
	TUN_IP = netip.MustParsePrefix("172.19.0.1/32")
	Routes = []netip.Prefix{
		netip.MustParsePrefix("1.0.0.0/8"),
		netip.MustParsePrefix("2.0.0.0/7"),
		netip.MustParsePrefix("4.0.0.0/6"),
		netip.MustParsePrefix("8.0.0.0/5"),
		netip.MustParsePrefix("16.0.0.0/4"),
		netip.MustParsePrefix("32.0.0.0/3"),
		netip.MustParsePrefix("64.0.0.0/2"),
		netip.MustParsePrefix("128.0.0.0/1"),
		netip.MustParsePrefix("198.18.0.0/15"),
	}
}

func tunPermission() bool {
	// 已经是 root
	if os.Geteuid() == 0 {
		return true
	}

	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
		exePath = resolved
	}

	u, err := user.Current()
	if err != nil {
		return false
	}

	sudoFile := "/etc/sudoers.d/prism-desktop"

	// 如果 sudoers 已存在 → 直接 sudo 启动
	if _, err := os.Stat(sudoFile); err == nil {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("sudo %q &", exePath))
		_ = cmd.Start()

		time.Sleep(100 * time.Millisecond)
		internal.App.Event.Emit(internal.AppConfig.EventNameServiceAppQuit, true)
		return false
	}

	// sudoers 内容（⚠️ 不要转义路径）
	line := fmt.Sprintf("%s ALL=(ALL) NOPASSWD: %s\n", u.Username, exePath)

	// 单行脚本（无换行）
	script := fmt.Sprintf(
		`echo '%s' > %s && chmod 440 %s`,
		line, sudoFile, sudoFile,
	)

	apple := fmt.Sprintf(`do shell script %q with administrator privileges`, script)
	cmd := exec.Command("osascript", "-e", apple)

	go func() {
		// 第一次尝试
		err := cmd.Run()

		if err != nil {
			// 修复（再写一次）
			fixScript := fmt.Sprintf(
				`echo '%s' > %s && chmod 440 %s`,
				line, sudoFile, sudoFile,
			)
			fixApple := fmt.Sprintf(`do shell script %q with administrator privileges`, fixScript)
			_ = exec.Command("osascript", "-e", fixApple).Run()
		}

		// 用 sudo 启动（NOPASSWD）
		exec.Command("sh", "-c", fmt.Sprintf("sudo %q &", exePath)).Start()

		time.Sleep(150 * time.Millisecond)
		internal.App.Event.Emit(internal.AppConfig.EventNameServiceAppQuit, true)
	}()

	return false
}
