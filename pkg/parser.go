package pkg

import (
	"github.com/tivvit/yap/pkg/structs"
	"github.com/tivvit/yap/pkg/utils"
	"log"
)

func parsePipelinev1(raw *structs.PipelineRaw) *structs.Pipeline {
	// todo check duplicate keys (is it even possible here? - yaml solves that?)
	p := structs.NewPipeline(raw.Version, raw.Settings, raw.Deps)
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
		if _, ik := vm[utils.IncludeKeyword]; ik {
			//log.Println("\tdetected import")
			name := vm[utils.IncludeKeyword].(string)
			t := "yaml"
			if _, iik := vm[utils.TypeKeyword]; iik {
				t = vm[utils.TypeKeyword].(string)
			}
			//log.Println("\t", t)
			switch t {
			case "script":
				ip, err := LoadScript(name)
				if err != nil {
					continue
				}
				//log.Println(ip)
				pp := ParsePipeline(ip)
				if deps, ok := vm[utils.DepsKeyword]; ok {
					for _, i := range deps.([]interface{}) {
						pp.Deps = append(pp.Deps, i.(string))
					}
				}
				p.Pipeline[k] = pp
			case "yaml":
				ip := LoadFile(name)
				//log.Println(ip)
				pp := ParsePipeline(ip)
				//log.Println(pp)
				if deps, ok := vm[utils.DepsKeyword]; ok {
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

func ParsePipeline(raw *structs.PipelineRaw) *structs.Pipeline {
	switch raw.Version {
	case 1.0:
		return parsePipelinev1(raw)
	}
	log.Fatal("Unsupported pipeline version")
	return nil
}
