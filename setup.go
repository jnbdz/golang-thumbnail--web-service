package main

import "os"

func addImgDirs() {

	var e Err
	var dn DirNav

	origPath := dn.getOrigDirPath()
	resizedPath := dn.getResizedDirPath()

	if stat, err := os.Stat(origPath); err == nil && stat.IsDir() {
		err = os.MkdirAll(origPath, 0777)
		if err != nil {
			e.setInternalServerError(err)
		}
	} else {
		e.setInternalServerError(err)
	}

	if stat, err := os.Stat(origPath); err == nil && stat.IsDir() {
		err = os.MkdirAll(resizedPath, 0777)
		if err != nil {
			e.setInternalServerError(err)
		}
	} else {
		e.setInternalServerError(err)
	}
}
