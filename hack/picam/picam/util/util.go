/*
 * Copyright (C) 2018 Josh A. Beam
 * All rights reserved.
 *
 * See https://github.com/joshb/pi-camera-go
 */

package util

import (
	"os"
	"os/user"
	"path"
)

func ConfigDir(components ...string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	configDir := path.Join(append([]string{u.HomeDir, ".pi-camera-go"}, components...)...)
	if err := os.MkdirAll(configDir, os.ModeDir|os.ModePerm); err != nil {
		return "", err
	}

	return configDir, nil
}
