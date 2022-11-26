package tools

import (
	"math"
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
func GetMaxS(x, y string) string {
	if x < y {
		return y
	}
	return x
}
func getGrade(xstr, maxstr string, grade int, trans func(float64) float64) (result string, err error) {
	x, err := strconv.ParseFloat(xstr, 64)
	if err != nil {
		return "", err
	}
	max, err := strconv.ParseFloat(maxstr, 64)
	if err != nil {
		return "", err
	}
	unit := trans(float64(max)) / float64(grade)
	x = trans(x)
	for i := 1; i <= grade; i++ {
		if x < float64(i)*unit {
			return strconv.Itoa(i), nil
		}
	}
	return strconv.Itoa(grade), nil
}
func GetGrades(xstr, maxstr []string, grade int, trans func(float64) float64) ([]string, error) {
	result := make([]string, 0, len(maxstr))
	for i, _ := range maxstr {
		grade, err := getGrade(xstr[i], maxstr[i], grade, trans)
		if err != nil {
			return nil, err
		}
		result = append(result, grade)
	}
	return result, nil
}
func Linear(x float64) float64 {
	return x
}
func Sigmoid(x float64) float64 {
	return 1/(1+math.Exp(-x)) - 0.5
}
func Tanh(x float64) float64 {
	return 2/(1+math.Exp(-2*x)) - 1
}
func Log(x float64) float64 {
	return math.Log(x + 1)
}
func F(f string) func(x float64) float64 {
	switch f {
	case "linear":
		return Linear
	case "log":
		return Log
	case "sigmoid":
		return Sigmoid
	case "tanh":
		return Tanh
	default:
		return Linear
	}
}
