package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/mpetavy/common"
	"os"
	"os/exec"
)

const (
	Ethernet     = "Ethernet"
	PiholeIp     = "192.168.1.1"
	CloudflareIp = "1.1.1.1"
)

var (
	LDFLAG_DEVELOPER = "mpetavy"                              // will be replaced with ldflag
	LDFLAG_HOMEPAGE  = "https://github.com/mpetavy/dnshobber" // will be replaced with ldflag
	LDFLAG_LICENSE   = common.APACHE                          // will be replaced with ldflag
	LDFLAG_VERSION   = "1.0.0"                                // will be replaced with ldflag
	LDFLAG_EXPIRE    = ""                                     // will be replaced with ldflag
	LDFLAG_GIT       = ""                                     // will be replaced with ldflag
	LDFLAG_COUNTER   = ""                                     // will be replaced with ldflag
)

var (
	icon           []byte
	menuPiHole     *systray.MenuItem
	menuCloudflare *systray.MenuItem
	menuQuit       *systray.MenuItem
)

func init() {
	common.Init(true, LDFLAG_VERSION, LDFLAG_GIT, "2019", "Simple DNS switcher", LDFLAG_DEVELOPER, LDFLAG_HOMEPAGE, LDFLAG_LICENSE, start, nil, nil, 0)
}

func start() error {
	var err error

	icon, err = Binpack_Icon_favicon32x32Ico.Unpack()
	if common.Error(err) {
		return err
	}

	go func() {
		systray.Run(onReady, onExit)
	}()

	return nil
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle(common.Title())
	systray.SetTooltip("DNS switcher")

	menuPiHole = systray.AddMenuItem("Pi-Hole", "Switch to Pi-Hole DNS")
	menuCloudflare = systray.AddMenuItem("Cloudflare", "Switch to Cloudflare DNS")
	systray.AddSeparator()
	menuQuit = systray.AddMenuItem("Quit", fmt.Sprintf("Quit %s", common.Title()))

	for {
		select {
		case <-menuPiHole.ClickedCh:
			menuCloudflare.Uncheck()
			menuPiHole.Check()

			dns(PiholeIp)
		case <-menuCloudflare.ClickedCh:
			menuPiHole.Uncheck()
			menuCloudflare.Check()

			dns(CloudflareIp)
		case <-menuQuit.ClickedCh:
			os.Exit(0)
		}
	}
}

func dns(ip string) {
	cmd := exec.Command("netsh", "interface", "ipv4", "set", "dns", fmt.Sprintf("name=%s", Ethernet), "static", ip)
	err := cmd.Run()
	if !common.Error(err) {
		common.Info("DNS server: %s", ip)
	}
}

func onExit() {
}

func main() {
	defer common.Done()

	common.Run(nil)
}
