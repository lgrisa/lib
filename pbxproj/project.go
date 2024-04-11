package pbxproj

import "strings"

// ProjectSection represent isa PBXNativeTarget
type ProjectSection struct {
	Id                      string
	TargetsUnityIphone      string
	TargetsUnityIphoneTests string
	TargetsUnityFramework   string
}

// parse PBXProject
func parseProjectSection(m map[string]interface{}, lines []string) map[string]*ProjectSection {
	ns := map[string]*ProjectSection{}

	for id, mm := range m {
		obj := mm.(map[string]interface{})
		for k, v := range obj {
			if k == "isa" && v.(string) == "PBXProject" {

				targets := lookupStrSlices(obj, "targets")

				nt := &ProjectSection{}

				for _, targetId := range targets {
					for _, line := range lines {
						if strings.Contains(line, targetId) && strings.Contains(line, "/* Unity-iPhone */") {
							nt.TargetsUnityIphone = targetId
						}

						if strings.Contains(line, targetId) && strings.Contains(line, "/* Unity-iPhone Tests */") {
							nt.TargetsUnityIphoneTests = targetId
						}

						if strings.Contains(line, targetId) && strings.Contains(line, "/* UnityFramework */") {
							nt.TargetsUnityFramework = targetId
						}

						if nt.TargetsUnityIphone != "" && nt.TargetsUnityIphoneTests != "" && nt.TargetsUnityFramework != "" {
							break
						}
					}
				}

				ns[id] = nt
			}
		}
	}

	return ns
}
