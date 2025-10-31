package generator

import (
	"sort"
	"strings"
)

// DetectScopeWithListImproved は最も多いスコープを自動選択、複数なら選択肢も作る
func DetectScopeWithListImproved(files []string) (string, []string) {
	counts := map[string]int{}
	fullScopes := map[string]bool{}

	excludeScopes := map[string]bool{
		"tagoly": true,
		"":       true,
	}

	for _, f := range files {
		parts := strings.Split(f, "/")
		var scope string
		if len(parts) < 2 {
			scope = "root"
		} else {
			switch parts[0] {
			case "internal", "cmd":
				scope = parts[0] + "/" + parts[1]
			default:
				scope = parts[0]
			}
		}
		if !excludeScopes[scope] {
			counts[scope]++
			fullScopes[scope] = true
		}
	}

	// スコープ一覧
	scopes := []string{}
	for s := range fullScopes {
		scopes = append(scopes, s)
	}

	// 複数スコープなら最後に "multiple"
	if len(scopes) > 1 {
		scopes = append(scopes, "multiple")
	}

	// 最多出現スコープをデフォルトに設定
	var defaultScope string
	if len(counts) > 0 {
		type kv struct {
			Key   string
			Value int
		}
		var ss []kv
		for k, v := range counts {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool { return ss[i].Value > ss[j].Value })
		defaultScope = ss[0].Key
	} else {
		defaultScope = "root"
	}

	return defaultScope, scopes
}
