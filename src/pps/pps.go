package pps

import "fmt"

func VersionString(version *Version) string {
	return fmt.Sprintf("%d.%d.%d%s", version.Major, version.Minor, version.Micro, version.Additional)
}

func GetNameToNode(pipeline *Pipeline) map[string]*Node {
	m := make(map[string]*Node)
	for name, element := range pipeline.NameToElement {
		if element.Node != nil {
			m[name] = element.Node
		}
	}
	return m
}