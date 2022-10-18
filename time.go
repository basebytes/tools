package tools

import (
	"container/ring"
	"regexp"
	"sync"
	"time"
)

const (
	dateFormat              = "2006-01-02"
	compactDateFormat       = "20060102"
	compactTimeFormat       = "150405"
	datetimeFormat          = "2006-01-02 15:04:05"
	maxTime                 = "23h59m59s"
	day               int64 = 24 * 60 * 60
)

var (
	timeRegexp = regexp.MustCompile("^([0,1][\\d]|2[0-3]):([0-5][\\d]):([0-5][\\d])$")
	lock       sync.Mutex
	idx        = 0
	data       = [2]int64{}
)

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

func init() {
	t, _ := time.ParseInLocation(dateFormat, time.Now().Format(dateFormat), time.Local)
	data[0] = t.Unix()
	data[1] = data[0] + day
}

// TransTimeFormat Trans time format from hh:mm:ss to hhhmmmsss.If timeString format invalid return itself
//
//eg. 23:10:34 => 23h10m34s
func TransTimeFormat(timeString string) string {
	if timeString == "24:00:00" {
		return maxTime
	}
	return timeRegexp.ReplaceAllString(timeString, "${1}h${2}m${3}s")
}

// GetTimeInSeconds return the seconds representation of the timeString
//
//timeStr in format hhhmmmsss ,e.g. 23h10m34s ,otherwise return -1
func GetTimeInSeconds(timeStr string) int64 {
	var result int64 = -1
	if d, e := time.ParseDuration(timeStr); e == nil {
		result = d.Milliseconds() / 1e3
	}
	return result
}

// CurrentDateTimeStr return current datetime in format yyyy-MM-dd HH:mm:ss
func CurrentDateTimeStr() string {
	return formatCurrentTime(datetimeFormat)
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

//TimeRingFromDaysAgo return time ring from days ago util next day.
func TimeRingFromDaysAgo(days int) *ring.Ring {
	timeRing := ring.New(days + 2)
	cur := timeRing
	for i := 0; i <= days; i, cur = i+1, cur.Next() {
		cur.Value = DaysBefore(i)
	}
	cur.Value = DaysBefore(-1)
	return timeRing
}

//SplitSeconds return unix timestamp in seconds with two part:begin of a day and the offset
func SplitSeconds(seconds int64) (int64, int64) {
	dayStart := CurrentDateInSeconds()
	offset := seconds - dayStart
	for offset < 0 || offset >= day {
		d := (offset+day)/day - 1
		dayStart, offset = dayStart+d*day, offset-d*day
	}
	return dayStart, offset
}
