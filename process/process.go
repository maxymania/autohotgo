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

import "github.com/yuin/gopher-lua"
import "os"
import "fmt"

type prinfo struct{
	procname string
	procimg  string
	procmem  uint64
}

func Install(L *lua.LState) {
	L.Register("processKill",processKill)
	L.Register("processList",processList)
	L.Register("processInfo",processInfo)
}

//CreateTable

// processKill(pid)
func processKill(L *lua.LState) int {
	pid := L.CheckInt(1)
	p,e := os.FindProcess(pid)
	if e!=nil { return 0 }
	p.Kill()
	p.Release()
	return 0
}

func processList(L *lua.LState) int {
	L.SetTop(1)
	lp := listProcesses()
	tab := L.CreateTable(len(lp),0)
	for i,v := range lp {
		L.RawSetInt(tab,i+1,lua.LNumber(v))
	}
	L.Push(tab)
	return 1
}

func processInfo(L *lua.LState) int {
	var i prinfo
	pid := L.CheckInt(1)
	L.SetTop(1)
	e := infoProcess(pid,&i)
	if e!=nil {
		L.Push(lua.LString(fmt.Sprint(e)))
		return 1
	}
	L.Push(lua.LNil)
	tab := L.NewTable()
	tab.RawSetString("procname",lua.LString(i.procname))
	tab.RawSetString("procimg",lua.LString(i.procimg))
	tab.RawSetString("procmem",lua.LNumber(i.procmem))
	L.Push(tab)
	return 2
}

