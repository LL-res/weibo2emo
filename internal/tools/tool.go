package tools

import (
	"encoding/csv"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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
func ConvStrs2Ints(strs []string) []float64 {
	result := make([]float64, len(strs))
	for _, v := range strs {
		t, _ := strconv.ParseFloat(v, 64)
		result = append(result, t)
	}
	return result
}
func DetectCSV(file *os.File) (*csv.Reader, error) {
	r := csv.NewReader(file)
	r.FieldsPerRecord = -1
	row, err := r.Read()
	if err != nil || len(row) == 0 {
		file.Seek(0, 0)
		r1 := csv.NewReader(file)
		return r1, errors.New("failed detection")
	}
	if utf8.ValidString(string(row[0])) {
		file.Seek(0, 0)
		r1 := csv.NewReader(file)
		return r1, nil
	}
	file.Seek(0, 0)
	reader := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())
	result := csv.NewReader(reader)
	result.FieldsPerRecord = -1
	return result, nil
}
func ENEmoKey(key string) string {
	switch key {
	case "死亡词":
		return "Death"
	case "仇视词":
		return "Hatred"
	case "脏话词":
		return "Swear"
	case "暂定词":
		return "Temporary"
	case "提及恐慌客体":
		return "Panic_object"
	case "提及紧缺物资":
		return "Short_supply"
	case "提及政府当局":
		return "About_regime"
	case "提及感官不适":
		return "Disturbing"
	default:
		strs := strings.Split(key, "(")
		if len(strs) < 2 {
			return key
		}
		return strings.TrimSuffix(strs[len(strs)-1], ")")
	}
}
