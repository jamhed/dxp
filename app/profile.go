package app

func profileGroups(p map[string]interface{}) []string {
	var groups []string
	if ifaces, ok := p["groups"].([]interface{}); ok {
		for _, iface := range ifaces {
			groups = append(groups, iface.(string))
		}
		return groups
	}
	return []string{}
}

func profileUser(p map[string]interface{}) string {
	if sub, ok := p["sub"].(string); ok {
		return sub
	}
	return ""
}
