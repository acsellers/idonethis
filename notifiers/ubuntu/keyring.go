package main

import (
	"encoding/gob"
	"os"
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

type saveData struct {
	Username string
	CheckMin int
}

func GetSimple() {
	f, e := os.Open(filepath.Join(basedir.ConfigHome, "idonethis"))
	if e != nil {
		return
	}
	defer f.Close()

	d := gob.NewDecoder(f)
	var sd saveData
	e = d.Decode(&sd)
	if e != nil {
		return
	}

	Username = sd.Username
	CheckMinutes = sd.CheckMin

}

func SaveSimple() {
	f, e := os.Create(filepath.Join(basedir.ConfigHome, "idonethis"))
	if e != nil {
		return
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	sd := saveData{
		Username: Username,
		CheckMin: CheckMinutes,
	}

	e = enc.Encode(sd)
	if e != nil {
		return
	}
}
