package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/acsellers/keyring"
	"github.com/rkoesters/xdg/basedir"
)

func GetPassword(username string) (string, error) {
	if pw, err := keyring.Get("idonethis_notify", username); err == nil {
		return pw, nil
	} else if err == keyring.ErrNotFound {
		return "", PasswordNotFound
	} else {
		return "", err
	}
}

func SetPassword(username, password string) error {
	return keyring.Set("idonethis_notify", username, password)
}

func GetUsername() string {
	fb, e := ioutil.ReadFile(filepath.Join(basedir.ConfigHome, "idonethis"))
	if e != nil {
		return ""
	}

	return string(fb)
}

func SetUsername(username string) {
	ioutil.WriteFile(filepath.Join(basedir.ConfigHome, "idonethis"), []byte(username), 0700)
}
