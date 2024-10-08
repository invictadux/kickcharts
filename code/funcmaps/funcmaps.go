package funcmaps

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"time"
)

func NewTemplate(files ...string) *template.Template {
	segments := strings.Split(files[0], "/")
	fileName := segments[len(segments)-1]
	return template.Must(template.New(fileName).Funcs(DefaultFunctions()).ParseFiles(files...))
}

func DefaultFunctions() template.FuncMap {
	return template.FuncMap{
		"ContainsString":        ContainsString,
		"FormatWithCommas":      FormatNumberCommas,
		"NearestThousandFormat": NearestThousandFormat,
		"Percent":               Percent,
		"FormatDate":            FormatDate,
		"ElapsedTime":           ElapseTime,
		"TimeSince":             TimeSince,
		"FormatStringDate":      FormatStringDate,
		"ToLower":               strings.ToLower,
		"CalculatePercent":      CalculatePercent,
		"Add":                   Add,
	}
}

func Add(a, b int) int {
	return a + b
}

func CalculatePercent(a, b int) float64 {
	return float64(b) / float64(a) * 100
}

func ContainsString(s, v string) bool {
	return strings.Contains(s, v)
}

func Percent(n float64, d int) string {
	format := fmt.Sprintf("%%.%df", d)
	formattedNumber := fmt.Sprintf(format, n)
	return formattedNumber
}

func FormatDate(date time.Time, format string) string {
	return date.Format(format)
}

func ElapseTime(sec float64) string {
	duration := time.Second * time.Duration(sec)
	return duration.String()
}

func TimeSince(t time.Time) string {
	duration := time.Since(t)
	years := int(duration.Hours() / (24 * 365))
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	switch {
	case years > 0:
		return fmt.Sprintf("%v years ago", years)
	case days > 0:
		return fmt.Sprintf("%v days ago", days)
	case hours > 0:
		return fmt.Sprintf("%v hours ago", hours)
	case minutes > 0:
		return fmt.Sprintf("%v minutes ago", minutes)
	default:
		return fmt.Sprintf("%v seconds ago", duration.Seconds())
	}
}

func FormatStringDate(dateStr, format string) string {
	date, err := time.Parse("2006-01-02 15:04:05", dateStr)

	if err != nil {
		return ""
	}

	return date.Format(format)
}

func RoundPrec(x float64, prec int) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow * sign
}

func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
	if math.IsNaN(number) || math.IsInf(number, 0) {
		number = 0
	}

	var ret string
	var negative bool

	if number < 0 {
		number *= -1
		negative = true
	}

	d, fract := math.Modf(number)

	if decimals <= 0 {
		fract = 0
	} else {
		pow := math.Pow(10, float64(decimals))
		fract = RoundPrec(fract*pow, 0)
	}

	if thousandsSep == "" {
		ret = strconv.FormatFloat(d, 'f', 0, 64)
	} else if d >= 1 {
		var x float64
		for d >= 1 {
			d, x = math.Modf(d / 1000)
			x = x * 1000
			ret = strconv.FormatFloat(x, 'f', 0, 64) + ret
			if d >= 1 {
				ret = thousandsSep + ret
			}
		}
	} else {
		ret = "0"
	}

	fracts := strconv.FormatFloat(fract, 'f', 0, 64)

	// "0" pad left
	for i := len(fracts); i < decimals; i++ {
		fracts = "0" + fracts
	}

	ret += decPoint + fracts

	if negative {
		ret = "-" + ret
	}
	return ret
}

func RoundInt(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

func FormatNumber(input float64) string {
	x := RoundInt(input)
	xFormatted := NumberFormat(float64(x), 2, ".", ",")
	return xFormatted
}

func NearestThousandFormat(n int) string {
	num := float64(n)

	if math.Abs(num) < 999.5 {
		xNum := FormatNumber(num)
		xNumStr := xNum[:len(xNum)-3]
		return string(xNumStr)
	}

	xNum := FormatNumber(num)
	xNumStr := xNum[:len(xNum)-3]
	xNumCleaned := strings.Replace(xNumStr, ",", " ", -1)
	xNumSlice := strings.Fields(xNumCleaned)
	count := len(xNumSlice) - 2
	unit := [4]string{"K", "M", "B", "T"}
	xPart := unit[count]

	afterDecimal := ""
	if xNumSlice[1][0] != 0 {
		afterDecimal = "." + string(xNumSlice[1][0])
	}
	final := xNumSlice[0] + afterDecimal + xPart
	return final
}

func FormatNumberCommas(n int) string {
	numStr := strconv.Itoa(n)
	length := len(numStr)
	separatorCount := (length - 1) / 3
	formattedNumber := make([]string, 0, length+separatorCount)

	if length%3 != 0 {
		formattedNumber = append(formattedNumber, numStr[:length%3])
	}

	for i := length % 3; i < length; i += 3 {
		formattedNumber = append(formattedNumber, numStr[i:i+3])
	}

	return strings.Join(formattedNumber, ",")
}
