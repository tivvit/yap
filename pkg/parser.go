package pkg

import (
	"github.com/tivvit/yap/pkg/structs"
	"log"
)

func ParsePipeline(raw *structs.PipelineRaw) *structs.Pipeline {
	// todo check duplicate keys
	p := structs.Pipeline{}
	for k, v := range *raw {
		//log.Println(k, v)
		vm, ok := v.(structs.PipelineRaw)
		// todo unify structure
		if !ok {
			//log.Println("PipelineRaw conversion err")
			vmi, ok := v.(map[interface{}]interface{})
			if !ok {
				//log.Println("map conversion err")
			}
			tmp := map[string]interface{}{}
			for ik, iv := range vmi {
				tmp[ik.(string)] = iv
			}
			vm = tmp
		}
		block := structs.Block{}
		if _, ik := vm[IncludeKeyword]; ik {
			//log.Println("\tdetected import")
			name := vm[IncludeKeyword].(string)
			t := "yaml"
			if _, iik := vm[TypeKeyword]; iik {
				t = vm[TypeKeyword].(string)
			}
			//log.Println("\t", t)
			switch t {
			case "script":
				ip := LoadScript(name)
				//log.Println(ip)
				p[k] = ParsePipeline(ip)
			case "yaml":
				ip := LoadFile(name)
				//log.Println(ip)
				pp := ParsePipeline(ip)
				//log.Println(pp)
				p[k] = pp
			default:
				log.Println("Unknown import type", t)
			}
		} else {
			block = *structs.NewBlockFromMap(k, vm)
			p[k] = block
		}
	}
	return &p
}
