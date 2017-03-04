package main

import (
	"github.com/fatih/color"
	"github.com/freignat91/blockchain/api"
	"os"
)

type bchainCLI struct {
	printColor [6]*color.Color
	server     string
	verbose    bool
	silence    bool
	debug      bool
	fullColor  bool
	userName   string
	keyPath    string
}

var currentColorTheme = "default"
var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
	colDebug   = 5
)

func (m *bchainCLI) init() error {
	m.setColors()
	//
	return nil
}

func (m *bchainCLI) printf(col int, format string, args ...interface{}) {
	if m.silence {
		return
	}
	colorp := m.printColor[0]
	if col > 0 && col < len(m.printColor) {
		colorp = m.printColor[col]
	}
	if !m.verbose && col == colInfo && !m.fullColor {
		return
	}
	if !m.debug && col == colDebug {
		return
	}
	colorp.Printf(format, args...)
}

func (m *bchainCLI) fatal(format string, args ...interface{}) {
	m.printf(colError, format, args...)
	os.Exit(1)
}

func (m *bchainCLI) pError(format string, args ...interface{}) {
	m.printf(colError, format, args...)
}

func (m *bchainCLI) pWarn(format string, args ...interface{}) {
	m.printf(colWarn, format, args...)
}

func (m *bchainCLI) pInfo(format string, args ...interface{}) {
	m.printf(colInfo, format, args...)
}

func (m *bchainCLI) pSuccess(format string, args ...interface{}) {
	m.printf(colSuccess, format, args...)
}

func (m *bchainCLI) pRegular(format string, args ...interface{}) {
	m.printf(colRegular, format, args...)
}

func (m *bchainCLI) pDebug(format string, args ...interface{}) {
	m.printf(colDebug, format, args...)
}

func (m *bchainCLI) setColors() {
	theme := config.colorTheme
	if theme == "dark" {
		m.printColor[0] = color.New(color.FgHiWhite)
		m.printColor[1] = color.New(color.FgHiBlack)
		m.printColor[2] = color.New(color.FgYellow)
		m.printColor[3] = color.New(color.FgRed)
		m.printColor[4] = color.New(color.FgGreen)
		m.printColor[5] = color.New(color.FgHiBlack)
	} else {
		m.printColor[0] = color.New(color.FgMagenta)
		m.printColor[1] = color.New(color.FgHiBlack)
		m.printColor[2] = color.New(color.FgYellow)
		m.printColor[3] = color.New(color.FgRed)
		m.printColor[4] = color.New(color.FgGreen)
		m.printColor[5] = color.New(color.FgHiBlack)
	}
	//add theme as you want.
}

func (m *bchainCLI) setAPI(api *api.BchainAPI) error {
	if m.silence {
		api.SetLogLevel("error")
	} else if m.verbose {
		api.SetLogLevel("info")
	} else if m.debug {
		api.SetLogLevel("debug")
	} else {
		api.SetLogLevel("warn")
	}
	if err := api.SetUser(m.userName, m.keyPath); err != nil {
		return err
	}
	return nil
}
