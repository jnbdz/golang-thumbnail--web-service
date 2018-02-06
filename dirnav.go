package main

import "path/filepath"

type DirNav struct {
	currentDirPath string
	origDirPath    string
	resizedDirPath string
	origImgPath    string
	resizedImgPath string
	imgOrigName    string
	imgResizedName string
}

func (dn *DirNav) setCurrentDirPath() {
	var e Err
	currPath, err := filepath.Abs("./")
	if err != nil {
		e.setInternalServerError(err)
	}
	currPath += "/"

	dn.currentDirPath = currPath
}

func (dn *DirNav) getCurrentDirPath() string {
	if dn.currentDirPath == "" {
		dn.setCurrentDirPath()
	}
	return dn.currentDirPath
}

func (dn *DirNav) setOrigDirPath() {
	dn.origDirPath = dn.getCurrentDirPath() + "imgs/orig/"
}

func (dn *DirNav) getOrigDirPath() string {
	if dn.origDirPath == "" {
		dn.setOrigDirPath()
	}
	return dn.origDirPath
}

func (dn *DirNav) setResizedDirPath() {
	dn.resizedDirPath = dn.getCurrentDirPath() + "imgs/resized/"
}

func (dn *DirNav) getResizedDirPath() string {
	if dn.resizedDirPath == "" {
		dn.setResizedDirPath()
	}
	return dn.resizedDirPath
}

func (dn *DirNav) setOrigImgPath() {
	dn.origImgPath = dn.getOrigDirPath() + dn.getImgOrigName()
}

func (dn *DirNav) getOrigImgPath() string {
	if dn.origImgPath == "" {
		dn.setOrigImgPath()
	}
	return dn.origImgPath
}

func (dn *DirNav) setResizedImgPath() {
	dn.resizedImgPath = dn.getResizedDirPath() + dn.getImgResizedName()
}

func (dn *DirNav) getResizedImgPath() string {
	if dn.resizedImgPath == "" {
		dn.setResizedImgPath()
	}
	return dn.resizedImgPath
}

func (dn *DirNav) setImgOrigName(name string) {
	dn.imgOrigName = name
}

func (dn *DirNav) getImgOrigName() string {
	return dn.imgOrigName
}

func (dn *DirNav) setImgResizedName(name string) {
	dn.imgResizedName = name
}

func (dn *DirNav) getImgResizedName() string {
	return dn.imgResizedName
}
