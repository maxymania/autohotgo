package main

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

import "fmt"
import "os"
import "github.com/yuin/gopher-lua"
import "github.com/maxymania/autohotgo/process"

func main(){
	if len(os.Args)<2 {
		fmt.Println("usage: ahg file")
		return
	}
	L := lua.NewState()
	process.Install(L)
	e := L.DoFile(os.Args[1])
	if e!=nil {
		fmt.Println("error",e)
	}
}


