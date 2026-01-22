package convert

import (
	"strconv"
	"strings"
)

func StringSliceToInt64Slice(s []string) ([]int64, error) {
	res := make([]int64, len(s))
	for i, v := range s {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = n
	}
	return res, nil
}

func ParseTagsStringSlice(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	// 统一处理分隔符：将中文逗号统一替换为英文逗号，便于分割
	raw = strings.NewReplacer(
		"，", ",",
	).Replace(raw)

	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		tag := strings.TrimSpace(p)
		if tag == "" {
			continue
		}
		// 检查重复
		if _, ok := seen[tag]; !ok {
			seen[tag] = struct{}{}
			tags = append(tags, tag)
		}
	}
	return tags
}
