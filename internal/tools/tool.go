package tools

import (
	"strconv"
	"strings"
)

func RemoveName(content string) string {
	strs := strings.Split(content, "//@")
	if len(strs) == 1 {
		return content
	}
	for i := range strs {
		ss := strings.Split(strs[i], ":")
		if len(ss) < 2 {
			continue
		}
		strs[i] = strings.Join(ss[1:], "")
	}
	return strings.Join(strs, ",")
}
func RemoveMark(slices []string) string {
	result := make([]string, 0, len(slices))
	toRemove := []string{
		"《",
		"》",
		"#",
		"@",
		"$",
		"%",
		"^",
		"&",
		"*",
		"(",
		")",
		"?",
		"/",
		"`",
		"~",
		"-",
		"=",
		"+",
		"{",
		"}",
		"[",
		"]",
		"\\",
		"+",
		"\"",
		"|",
		"【",
		"】",
		"：",
		"；",
		"‘",
		"’",
		"“",
		"”",
		":",
		";",
		"'",
		"'",
		".",
		"<",
		">",
		"·",
	}
	for _, v := range slices {
		if CheckExistance(v, toRemove) {
			continue
		}
		result = append(result, v)
	}
	return strings.Join(result, " ")
}

func CheckExistance(target string, set []string) bool {
	for _, v := range set {
		if target == v {
			return true
		}
	}
	return false
}
func CheckExistanceInMap(target string, imap map[string][]string) []string {
	result := make([]string, 0)
	for k, v := range imap {
		if CheckExistance(target, v) {
			result = append(result, k)
		}
	}
	return result
}
func AddSlices(targets ...[]int) []int {
	if len(targets) == 0 {
		return []int{}
	}
	result := make([]int, len(targets[0]))
	for _, v := range targets {
		AddTwoSlices(result, v)
	}
	return result
}
func AddTwoSlices(t1, t2 []int) {
	for i := range t1 {
		t1[i] = t1[i] + t2[i]
	}
}
func ConvIs2Ss(slices []int) []string {
	result := make([]string, 0, len(slices))
	for _, v := range slices {
		result = append(result, strconv.Itoa(v))
	}
	return result
}
