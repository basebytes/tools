package tools

import (
	"regexp"
	"strconv"
	"time"
)

var reg = regexp.MustCompile(`(((\d{2})\d{2})\d{2})((\d{4})(\d{4}))(?:\d{2}(\d))[\dxX]`)

// ExtractCerInfo isMale 1：女性 2：男性
func ExtractCerInfo(cerNo string) (province, city, district, birthStr string, birth int64, age int, isMale int8) {
	if len(cerNo) < 18 {
		return
	}
	if subs := reg.FindStringSubmatch(cerNo); len(subs) == 8 {
		district = subs[1]
		city = subs[2]
		province = subs[3]
		birthStr = subs[4]
		if t, err := time.ParseInLocation("20060102", birthStr, time.Local); err == nil {
			birth = t.Unix()
		}
		if n, err := strconv.Atoi(subs[7]); err == nil {
			isMale = int8(n%2 + 1)
		}
	}
	return
}

func ComputeAge(birthStr string, baseDate *time.Time) (age int) {
	if len(birthStr) < 4 {
		return
	}
	if y, err := strconv.Atoi(birthStr[:4]); err == nil {
		age = baseDate.Year() - y
		if baseDate.Format("0102") < birthStr[4:] {
			age--
		}
	}
	return age
}

func IsValidCardNo18(CardNo18 *[]byte) bool {
	nLen := len(*CardNo18)
	if nLen != 18 {
		return false
	}

	nSum := 0
	for i := 0; i < nLen-1; i++ {
		n, _ := strconv.Atoi(string((*CardNo18)[i]))
		nSum += n * weight[i]
	}
	mod := nSum % 11
	if validValue[mod] == (*CardNo18)[17] {
		return true
	}

	return false
}

var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var validValue = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
