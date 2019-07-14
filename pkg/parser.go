package pkg

import (
	"github.com/tivvit/yap/pkg/structs"
	"log"
)

func ParsePipeline(raw *structs.PipelineRaw) *structs.Pipeline {
	// todo check duplicate keys (is it even possible here? - yaml solves that?)
	p := structs.NewPipeline(raw.Version, raw.Settings, raw.Deps)
	// todo check version and determine loader
	for k, v := range (*raw).Pipeline {
		//log.Println(k, v)
		vm, ok := v.(map[string]interface{})
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
				if deps, ok := vm[DepsKeyword]; ok {
					for _, i := range deps.([]interface{}) {
						ip.Deps = append(ip.Deps, i.(string))
					}
				}
				p.Pipeline[k] = ParsePipeline(ip)
			case "yaml":
				ip := LoadFile(name)
				//log.Println(ip)
				pp := ParsePipeline(ip)
				//log.Println(pp)
				if deps, ok := vm[DepsKeyword]; ok {
					for _, i := range deps.([]interface{}) {
						pp.Deps = append(pp.Deps, i.(string))
					}
				}
				p.Pipeline[k] = pp
			default:
				log.Println("Unknown import type", t)
			}
		} else {
			p.Pipeline[k] = structs.NewBlockFromMap(k, vm)
		}
	}
	return p
}
