package pbxproj

import "strings"

// NativeTarget represent isa PBXNativeTarget
type NativeTarget struct {
	Id                     string
	buildConfigurationList string
	productReference       string
	productType            string
	productName            string
	BuildPhases            []string
	dependencies           []string
	name                   string
	buildRules             []string
	BuildPhasesFrameworks  string
}

// parse PBXNativeTarget
func parseNativeTargets(m map[string]interface{}, lines []string) map[string]*NativeTarget {
	ns := map[string]*NativeTarget{}

	for id, mm := range m {
		obj := mm.(map[string]interface{})
		for k, v := range obj {
			if k == "isa" && v.(string) == "PBXNativeTarget" {

				buildPhase := lookupStrSlices(obj, "buildPhases")

				buildPhaseFrameworksId := ""
				//对应解析Frameworks
				for _, buildPhaseId := range buildPhase {
					for _, line := range lines {
						if strings.Contains(line, buildPhaseId) && strings.Contains(line, "/* Frameworks */") {
							buildPhaseFrameworksId = buildPhaseId
							break
						}
					}

					if buildPhaseFrameworksId != "" {
						break
					}
				}

				nt := NativeTarget{
					id,
					lookupStr(obj, "buildConfigurationList"),
					lookupStr(obj, "productReference"),
					lookupStr(obj, "productType"),
					lookupStr(obj, "productName"),
					buildPhase,
					lookupStrSlices(obj, "dependencies"),
					lookupStr(obj, "name"),
					lookupStrSlices(obj, "buildRules"),
					buildPhaseFrameworksId,
				}
				ns[id] = &nt
			}
		}
	}
	return ns
}
