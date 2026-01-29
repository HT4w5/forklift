package patch

import "strings"

func Patch(base any, patch any) any {
	baseMap, baseOK := base.(map[string]any)
	patchMap, patchOK := patch.(map[string]any)

	if !baseOK || !patchOK {
		return patch
	}

	for k, v := range patchMap {
		realKey := realKey(k)

		// Force overwrite
		if strings.HasSuffix(k, "!") {
			baseMap[realKey] = v
			continue
		}

		// Prepend
		if strings.HasPrefix(k, "+") {
			if baseArr, ok := baseMap[realKey].([]any); ok {
				if patchArr, ok := v.([]any); ok {
					baseMap[realKey] = append(patchArr, baseArr...)
					continue
				}
			}
		}

		// Append
		if strings.HasSuffix(k, "+") {
			if baseArr, ok := baseMap[realKey].([]any); ok {
				if patchArr, ok := v.([]any); ok {
					baseMap[realKey] = append(baseArr, patchArr...)
					continue
				}
			}
		}

		// Recursive patch
		if subBase, ok := baseMap[realKey]; ok {
			baseMap[realKey] = Patch(subBase, v)
		} else {
			baseMap[realKey] = v
		}
	}

	return base
}

func realKey(k string) string {
	if len(k) == 0 {
		return k
	}
	start := strings.Index(k, "<")
	if start >= 0 {
		end := strings.LastIndex(k, ">")
		if end >= 0 && end > start {
			return k[start+1 : end]
		}
	}

	// Strip modifiers
	end := len(k) - 1
	if k[end] == '!' || k[end] == '+' {
		return k[:end]
	}
	if k[0] == '+' {
		return k[1:]
	}

	return k
}
