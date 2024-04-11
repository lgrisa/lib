package pbxproj

// ProjectSection represent isa PBXNativeTarget
type FrameworkBuildPhase struct {
	Id    string
	Files []string
}

// parse PBXProject
func parseFrameworkBuildPhase(m map[string]interface{}) map[string]*FrameworkBuildPhase {
	ns := map[string]*FrameworkBuildPhase{}

	for id, mm := range m {
		obj := mm.(map[string]interface{})
		for k, v := range obj {
			if k == "isa" && v.(string) == "PBXFrameworksBuildPhase" {

				nt := &FrameworkBuildPhase{
					id,
					lookupStrSlices(obj, "files"),
				}
				ns[id] = nt
			}
		}
	}

	return ns
}
