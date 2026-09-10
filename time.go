package tools

import (
	"container/ring"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dateFormat              = "2006-01-02"
	compactDateFormat       = "20060102"
	compactTimeFormat       = "150405"
	DatetimeFormat          = "2006-01-02 15:04:05"
	maxTime                 = "23h59m59s"
	minute                  = 60
	hour                    = 60 * minute
	day               int64 = 24 * hour
	dayMillis         int64 = day * 1000
)

func GetCurrentDate() *time.Time {
	return ExtractDate(GetLocalTime())
}

func GetLocalTime() *time.Time {
	t := time.Now().In(location)
	return &t
}

func GetTodayCompactString() string {
	return GetDateCompactString(GetLocalTime())
}

func GetDateCompactString(time *time.Time) string {
	return GetTimeString(compactDateFormat, time)
}

func GetTimeString(format string, time *time.Time) string {
	return time.Format(format)
}

func FormatMilliSecondsToStr(millSec int64, fmt string) string {
	return GetTimeString(fmt, GetTimeWithMilliSeconds(millSec))
}

func FormatMilliSecondsToTimeStr(millSec int64) string {
	return FormatMilliSecondsToStr(millSec, DatetimeFormat)
}

func FormatMilliSecondsToDateStr(millSec int64) string {
	return FormatMilliSecondsToStr(millSec, dateFormat)
}

func GetTimeWithMilliSeconds(millSec int64) *time.Time {
	t := time.UnixMilli(millSec)
	return &t
}

func ParseMilliSecondsFromStr(timeStr string) (millSec int64) {
	if t, err := ParseTimeFromStr(timeStr); err == nil {
		millSec = GetMilliSeconds(t)
	}
	return
}

func ParseTimeFromStr(timeStr string) (*time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	fmtStr := dateFormat
	if len(timeStr) > 10 {
		fmtStr = DatetimeFormat
	}
	return ParseTimeInLocation(timeStr, fmtStr)
}

func GetMilliSeconds(t *time.Time) (millSec int64) {
	if t != nil {
		millSec = t.UnixMilli()
	}
	return
}

func GetCurrentSeconds() int64 {
	return GetLocalTime().Unix()
}

func GetCurrentMillSeconds() int64 {
	return GetLocalTime().UnixMilli()
}

func ParseTimeInLocation(timeStr, timeFmt string) (*time.Time, error) {
	t, err := time.ParseInLocation(timeFmt, timeStr, location)
	return &t, err
}

func GetWeekStartDate(now *time.Time) *time.Time {
	weekStart := now.AddDate(0, 0, -(6+int(now.Weekday()))%7)
	return ExtractDate(&weekStart)
}

func GetWeekEndDate(now *time.Time) *time.Time {
	weekEnd := now.AddDate(0, 0, (7-int(now.Weekday()))%7)
	return ExtractDate(&weekEnd)
}

func ExportToUTCMills(localMillSec int64) (millSec int64) {
	return localMillSec + diff
}

func ImportFromUTCMills(utcMillSec int64) (localMillSec int64) {
	return utcMillSec - diff
}

func ExtractDate(value *time.Time) *time.Time {
	millis := ((value.UnixMilli()+diff)/dayMillis)*dayMillis - diff
	t := time.UnixMilli(millis)
	return &t
}

func DaysAgo(base *time.Time, n int) *time.Time {
	t := base.AddDate(0, 0, -n)
	return &t
}

func DaysAfter(base *time.Time, n int) *time.Time {
	t := base.AddDate(0, 0, n)
	return &t
}

func MonthsAgo(base *time.Time, n int) *time.Time {
	t := base.AddDate(0, -n, 0)
	return &t
}

func CreateDate(year, month, day int) *time.Time {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
	return &t
}

func CreateTime(year, month, day, hour, min, sec int) *time.Time {
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, location)
	return &t
}

func MonthDays(t *time.Time) (days int) {
	month := int(t.Month())
	if month == 2 {
		year := t.Year()
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			days = 29
		} else {
			days = 28
		}
	} else if month > 0 && month < 13 {
		days = monthDays[month]
	}
	return
}

func TransTimeSec(timeStr string) (sec int) {
	if subs := timeReg.FindStringSubmatch(timeStr); len(subs) > 2 {
		if sec, _ = strconv.Atoi(subs[1]); sec > 0 {
			switch subs[2] {
			case "h", "H":
				sec *= hour
			case "m", "M":
				sec *= minute
			}
		}
	}
	return
}

func TransMinuteToSec(timeStr string) (sec int) {
	if d, err := time.ParseDuration(timeStr + "m"); err == nil {
		sec = int(d.Milliseconds() / 1000)
	} else {
		sec = -1
	}
	return
}

func TransSecondStrToSec(timeStr string) (sec int) {
	durations := strings.Split(timeStr, ".")
	if d, err := time.ParseDuration(durations[0] + "s"); err == nil {
		sec = int(d.Milliseconds() / 1000)
	} else {
		sec = -1
	}
	return
}

func init() {
	if location == nil {
		location = time.FixedZone("CST", 8*60*60)
	}
	_, offset := time.Now().Zone()
	diff = int64(offset * 1000)
	t, _ := time.ParseInLocation(dateFormat, time.Now().Format(dateFormat), location)
	data[0] = t.Unix()
	data[1] = data[0] + day
}

// CurrentDateInSeconds return the seconds representation of the current date in local
func CurrentDateInSeconds() int64 {
	if time.Now().Unix() > data[idx]+day {
		lock.Lock()
		defer lock.Unlock()
		if time.Now().Unix() > data[idx]+day {
			data[idx] = data[idx] + day*2
			idx = 1 - idx
		}
	}
	return data[idx]
}

// DaysBefore return the seconds representation of the start of days before now in local
func DaysBefore(days int) int64 {
	cur := CurrentDateInSeconds()
	return cur - int64(days)*day
}

// TransTimeFormat Trans time format from hh:mm:ss to hhhmmmsss.If timeString format invalid return itself
//
// eg. 23:10:34 => 23h10m34s
func TransTimeFormat(timeString string) string {
	if timeString == "24:00:00" {
		return maxTime
	}
	return timeRegexp.ReplaceAllString(timeString, "${1}h${2}m${3}s")
}

// GetTimeInSeconds return the seconds representation of the timeString
//
// timeStr in format hhhmmmsss ,e.g. 23h10m34s ,otherwise return -1
func GetTimeInSeconds(timeStr string) int64 {
	var result int64 = -1
	if d, e := time.ParseDuration(timeStr); e == nil {
		result = d.Milliseconds() / 1e3
	}
	return result
}

// CurrentDateTimeStr return current datetime in format yyyy-MM-dd HH:mm:ss
func CurrentDateTimeStr() string {
	return formatCurrentTime(DatetimeFormat)
}

// CurrentDateStr return current date in format yyyy-MM-dd
func CurrentDateStr() string {
	return formatCurrentTime(dateFormat)
}

// CurrentDataCompactStr return current date in format yyyyMMdd
func CurrentDataCompactStr() string {
	return formatCurrentTime(compactDateFormat)
}

// CurrentTimeCompactStr return current time in format HHmmss
func CurrentTimeCompactStr() string {
	return formatCurrentTime(compactTimeFormat)
}

func formatCurrentTime(format string) string {
	return time.Now().Format(format)
}

// TimeRingFromDaysAgo return time ring from days ago util next day.
func TimeRingFromDaysAgo(days int) *ring.Ring {
	timeRing := ring.New(days + 2)
	cur := timeRing
	for i := 0; i <= days; i, cur = i+1, cur.Next() {
		cur.Value = DaysBefore(i)
	}
	cur.Value = DaysBefore(-1)
	return timeRing
}

// SplitSeconds return unix timestamp in seconds with two part:begin of a day and the offset
func SplitSeconds(seconds int64) (int64, int64) {
	dayStart := CurrentDateInSeconds()
	offset := seconds - dayStart
	for offset < 0 || offset >= day {
		d := (offset+day)/day - 1
		dayStart, offset = dayStart+d*day, offset-d*day
	}
	return dayStart, offset
}

var (
	location, _ = time.LoadLocation("Asia/Shanghai")
	timeRegexp  = regexp.MustCompile("^([0,1][\\d]|2[0-3]):([0-5][\\d]):([0-5][\\d])$")
	timeReg     = regexp.MustCompile("^(\\d+)([HhMm])")
	monthDays   = []int{0, 31, 0, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	lock        sync.Mutex
	idx         = 0
	data        = [2]int64{}
	diff        int64
)
