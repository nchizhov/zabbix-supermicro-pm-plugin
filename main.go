package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"zbplugin"

	"golang.zabbix.com/sdk/errs"
	sdkplugin "golang.zabbix.com/sdk/plugin"
	"golang.zabbix.com/sdk/plugin/flag"
)

const copyrightMessage = "Copyright 2026-%d Inok Ltd."

var (
	PLUGIN_VERSION_MAJOR = 1
	PLUGIN_VERSION_MINOR = 0
	PLUGIN_VERSION_PATCH = 0
	PLUGIN_VERSION_RC    = "a1"
	PLUGIN_LICENSE_YEAR  = 2026
)

func main() {
	args, err := flag.HandleFlags()
	if err != nil {
		exitWithError(errs.Wrap(err, "failed to handle flags: "))
	}

	pluginInfo := &sdkplugin.Info{
		Name:             zbplugin.Name,
		BinName:          os.Args[0],
		CopyrightMessage: fmt.Sprintf(copyrightMessage, PLUGIN_LICENSE_YEAR),
		MajorVersion:     PLUGIN_VERSION_MAJOR,
		MinorVersion:     PLUGIN_VERSION_MINOR,
		PatchVersion:     PLUGIN_VERSION_PATCH,
		Alphatag:         PLUGIN_VERSION_RC,
	}

	ctx := context.Background()

	p, err := zbplugin.New(ctx)
	if err != nil {
		exitWithError(errs.Wrap(err, "failed to initialize plugin: "))
	}

	err = flag.DecideActionFromFlags(args, p, pluginInfo, nil)
	if err != nil {
		if errors.Is(err, errs.ErrExitGracefully) {
			exitGracefully()
		}
		exitWithError(errs.Wrap(err, "failed to execute plugin functions: "))
	}

	err = p.Run()
	if err != nil {
		exitWithError(errs.Wrap(err, "failed to run plugin: "))
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func exitGracefully() {
	os.Exit(0)
}
