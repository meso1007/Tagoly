package main

import (
	"flag"
	"fmt"
	"tagoly/internal/hook"
)

func processInstallHookCommand(args []string) error {
	fs := flag.NewFlagSet("install-hook", flag.ContinueOnError)
	force := fs.Bool("force", false, "Overwrite an existing commit-msg hook")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("invalid install-hook arguments: %v", err)
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("install-hook does not accept positional arguments")
	}

	hookPath, err := hook.InstallCommitMsgHook(".", hook.InstallOptions{Force: *force})
	if err != nil {
		return err
	}
	fmt.Printf("Installed commit-msg hook: %s\n", hookPath)
	return nil
}
