// +build windows

package process

/*
   Copyright 2016 Simon Schmidt

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import "syscall"
import "unsafe"
import "errors"
import "path/filepath"

const NULL = uintptr(0)

var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess              = modKernel32.NewProc("OpenProcess")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	
	modPsapi                     = syscall.NewLazyDLL("psapi.dll")
	procEnumProcesses            = modPsapi.NewProc("EnumProcesses")
	procGetProcessMemoryInfo     = modPsapi.NewProc("GetProcessMemoryInfo")
	procGetModuleBaseName        = modPsapi.NewProc("GetModuleBaseNameW")
	procGetProcessImageFileName  = modPsapi.NewProc("GetProcessImageFileNameW")
)

type PROCESS_MEMORY_COUNTERS_EX struct {
	cb                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr
}

func w2s(w []uint16) string{
	r := make([]rune,len(w))
	for i,ww := range w {
		r[i]=rune(ww)
	}
	return string(r)
}


func listProcesses() []int{
	lst := make([]int32,1<<16)
	rs := 0
	lp := unsafe.Pointer(&lst[0])
	rp := unsafe.Pointer(&rs)
	procEnumProcesses.Call(uintptr(lp),uintptr(len(lst)*4),uintptr(rp))
	rs/=4
	dlst := make([]int,rs)
	for i := 0 ; i<rs ; i++ {
		dlst[i] = int(lst[i])
	}
	return dlst[:rs]
}

func infoProcess(pid int,i *prinfo) error {
	name := make([]uint16,1<<11)
	var counters PROCESS_MEMORY_COUNTERS_EX
	counters.cb = uint32(unsafe.Sizeof(counters))
	cp := uintptr(unsafe.Pointer(&counters))
	// PROCESS_QUERY_INFORMATION
	h,_,e := procOpenProcess.Call(uintptr(0x0400),NULL,uintptr(pid))
	if h==0 && e!=nil { return e }
	defer procCloseHandle.Call(h)
	ok,_,e := procGetProcessMemoryInfo.Call(h,cp,unsafe.Sizeof(counters))
	if ok==0 && e!=nil { return e }
	if ok==0 { return errors.New("failed") }
	if counters.PrivateUsage==0 { counters.PrivateUsage = counters.PagefileUsage }
	i.procmem = uint64(counters.PrivateUsage)
	namep := uintptr(unsafe.Pointer(&name[0]))
	namel := uintptr(len(name))
	//namel,_,_ = procGetModuleBaseName.Call(h,NULL,namep,uintptr(len(name)))
	//if namel==0 {
		namel,_,_ = procGetProcessImageFileName.Call(h,namep,uintptr(len(name)))
	//}
	i.procimg = w2s(name[:int(namel)])
	i.procname = filepath.Base(i.procimg)
	return nil
}

