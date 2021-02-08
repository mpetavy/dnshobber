package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/mpetavy/common"
	"os"
	"os/exec"
	"strings"
)

const (
	Ethernet = "Ethernet"
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
	icon       []byte
	dnsServers common.MultiValueFlag
	menus      []*systray.MenuItem
)

func init() {
	common.Init(false, LDFLAG_VERSION, LDFLAG_GIT, "2021", "Simple DNS switcher", LDFLAG_DEVELOPER, LDFLAG_HOMEPAGE, LDFLAG_LICENSE, start, nil, nil, 0)

	flag.Var(&dnsServers, "s", "DNS servers")
	dnsServers = []string{"192.168.1.1", "192.168.1.7", "1.1.1.1", "8.8.4.4"}
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

	menus = make([]*systray.MenuItem, 0)
	clickCh := make(chan *systray.MenuItem)

	for _, dnsServer := range dnsServers {
		menu := systray.AddMenuItem(dnsServer, fmt.Sprintf("Switch to %s DNS", dnsServer))
		menus = append(menus, menu)

		go func(menu *systray.MenuItem) {
			for {
				<-menu.ClickedCh
				clickCh <- menu
			}
		}(menu)
	}

	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", fmt.Sprintf("Quit %s", common.Title()))
	go func() {
		for {
			<-menuQuit.ClickedCh
			os.Exit(0)
		}
	}()

	for {
		select {
		case clickMenu := <-clickCh:
			for _, menu := range menus {
				menu.Uncheck()
			}

			clickMenu.Check()

			t := clickMenu.String()
			t = t[:len(t)-2]

			p := strings.LastIndex(t, "\"")
			t = t[p+1:]

			dns(t)
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
