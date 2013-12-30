package main

import (
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/phrozen/gravatar"
	"github.com/rkoesters/xdg/basedir"
)

var (
	savedPics map[string]string
	cacheDir  string
)

func init() {
	cacheDir = filepath.Join(basedir.CacheHome, "idonethis")
	fi, e := os.Stat(cacheDir)
	if e != nil {
		e = os.Mkdir(cacheDir, 0755)
		if e != nil {
			panic(e)
		}
		fi, e = os.Stat(cacheDir)
	}
	if !fi.IsDir() {
		panic("cache folder is not folder")
	}

	savedPics = make(map[string]string)
}

func LocalAvatarPath(email string) (string, error) {
	// if we've already saved this person
	if p, ok := savedPics[email]; ok {
		return p, nil
	}

	g := gravatar.NewGravatar(email)
	i, e := g.Avatar(64)
	if e != nil {
		return "", e
	}

	fn := filepath.Join(cacheDir, g.Hash()+".jpeg")
	f, e := os.Create(fn)
	if e != nil {
		return "", e
	}

	e = jpeg.Encode(f, i, &jpeg.Options{80})
	if e != nil {
		return "", e
	}
	f.Close()

	savedPics[email] = fn
	return fn, nil
}
