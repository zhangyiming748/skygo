package custom_util

import (
	"regexp"
	"strconv"
	"time"
)

var ShanghaiLoc *time.Location

func init() {
	var err error
	ShanghaiLoc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}

/**
 * end_day int 查询结束时间, 单位为天, 默认0表示当前时间为查询结束时间 查询结束时间为(当前时间-$end)
 * time_span int 查询时间跨度, 单位为天,  默认为1,查询时间区域为($end-$timeSpan, $end)
 */
func GetTimeRange(endDay int, timeSpan int) (startTime, endTime time.Time) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	startTime = time.Date(now.Year(), now.Month(), now.Day()-timeSpan, now.Hour(), now.Minute(), now.Second(), 0, loc)
	endTime = time.Date(now.Year(), now.Month(), now.Day()-endDay, now.Hour(), now.Minute(), now.Second(), 0, loc)
	return
}

func GetThisWeekStartTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	return time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, loc)
}

// 时间戳 to 时间
func TimestampToString(timestamp int64) string {
	var loc, _ = time.LoadLocation("Asia/Shanghai") // 上海时区
	return time.Unix(timestamp, 0).In(loc).Format("2006-01-02 15:04:05")
}

func TimestampToStringNoSpace(timestamp int64) string {
	var loc, _ = time.LoadLocation("Asia/Shanghai") // 上海时区
	return time.Unix(timestamp, 0).In(loc).Format("20060102150405")
}

func DurationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

/**
 * CalculateTimeRangeString
 *
 * 计算一个时间范围,返回时间字符串
 *
 * @param {int}      [endTime=0]            查询结束时间(单位:取决于参数time_interval),查询结束时间计算方法为:(当前时间)-$end*$time_interval
 * @param {int}      [timeSpan=0]           查询时间跨度
 * @param {string}   [timeInterval=day]     查询时间粒度(支持时间单位:day、week、month、quarter、year)
 */
func CalculateTimeRangeString(endTime, timeSpan int, timeInterval string) (startDateStr, endDateStr string) {
	startDate, endDate := CalculateTimeRange(endTime, timeSpan, timeInterval)
	startDateStr = startDate.Format("2006-01-02")
	endDateStr = endDate.Format("2006-01-02")
	return
}

/**
 * CalculateTimeRange
 *
 * 计算一个时间范围
 *
 * @param {int}      [endTime=0]            查询结束时间(单位:取决于参数time_interval),查询结束时间计算方法为:(当前时间)-$end*$time_interval
 * @param {int}      [timeSpan=0]           查询时间跨度
 * @param {string}   [timeInterval=day]     查询时间粒度(支持时间单位:min、hour、day、week、month、quarter、year)
 */
func CalculateTimeRange(endTime, timeSpan int, timeUnit string) (startDate, endDate time.Time) {
	var loc, _ = time.LoadLocation("Asia/Shanghai") // 上海时区
	now := time.Now().In(loc)
	year, month, day := now.Date()
	timeNum, timeInterval := SplitTimeUnit(timeUnit)
	switch timeInterval {
	case "min":
		now = time.Unix(now.Unix()-now.Unix()%(60*int64(timeNum)), 0).In(loc)
		year, month, day = now.Date()
		hour := now.Hour()
		min := now.Minute()
		startDate = time.Date(year, month, day, hour, min-(endTime+timeSpan-1)*timeNum, 0, 0, loc)
		endDate = time.Date(year, month, day, hour, min-(endTime-1)*timeNum, 0, 0, loc)
	case "hour":
		hour := now.Hour()
		startDate = time.Date(year, month, day, hour-(endTime+timeSpan-1), 0, 0, 0, loc)
		endDate = time.Date(year, month, day, hour-(endTime-1), 0, 0, 0, loc)
	case "day":
		currentDate := time.Date(year, month, day, 0, 0, 0, 0, loc)
		startDate = currentDate.AddDate(0, 0, -(endTime + timeSpan - 1))
		endDate = startDate.AddDate(0, 0, timeSpan)
	case "week":
		spaceDay := (int(now.Weekday()) + 6) % 7
		currentDate := time.Date(year, month, day, 0, 0, 0, 0, loc)
		startDate = currentDate.AddDate(0, 0, -7*(endTime+timeSpan-1)-spaceDay)
		endDate = startDate.AddDate(0, 0, 7*timeSpan)
	case "month":
		currentDate := time.Date(year, month, 1, 0, 0, 0, 0, loc)
		startDate = currentDate.AddDate(0, -(endTime + timeSpan - 1), 0)
		endDate = startDate.AddDate(0, timeSpan, 0)
	case "quarter":
		spaceMonth := (int(month) + 2) % 3
		currentDate := time.Date(year, month, 1, 0, 0, 0, 0, loc)
		startDate = currentDate.AddDate(0, -spaceMonth-3*(endTime+timeSpan-1), 0)
		endDate = startDate.AddDate(0, timeSpan*3, 0)
	case "year":
		currentDate := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		startDate = currentDate.AddDate(-(endTime + timeSpan - 1), 0, 0)
		endDate = startDate.AddDate(timeSpan, 0, 0)
	}
	return
}

func SplitTimeUnit(timeUnit string) (int, string) {
	reg := regexp.MustCompile(`\s*(\d*)(\w*)`)
	if unitInfo := reg.FindStringSubmatch(timeUnit); len(unitInfo) == 3 {
		if i, err := strconv.Atoi(unitInfo[1]); err == nil {
			return i, unitInfo[2]
		}
	}
	return 1, timeUnit
}

func GetCurrentMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func MillisecondToTime(mtimestamp int64) time.Time {
	second := int64(mtimestamp / 1000)
	msecond := int64(mtimestamp % 1000)
	return time.Unix(second, msecond*1e6)
}

// 判断时间是当年的第几周
func WeekByDate(t time.Time) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	// 今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	return week
}

// 获取当天零点的时间
func GetTodayFormatDate() time.Time {
	var loc, _ = time.LoadLocation("Asia/Shanghai") // 上海时区
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func StringToTime(dateStr string) time.Time {
	defaultTimeTemplate := "2006-01-02 15:04:05" // 时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
	var loc, _ = time.LoadLocation("Asia/Shanghai")
	if t, err := time.ParseInLocation(defaultTimeTemplate, dateStr, loc); err == nil {
		return t
	} else {
		panic(err)
	}
}

/*
*
获取上一周的时间范围，从周一的0点，到周日的最后，
比如
2020-02-10 00:00:00 +0800 CST
2020-02-16 23:59:59.999999999 +0800 CST
*/
func GetLastWeekRange(t time.Time) (MondayTime, SundayTime time.Time) {
	diffDay := 0
	weekDiff := time.Hour*24*7 - 1

	w := t.In(ShanghaiLoc).Weekday()
	if w > time.Sunday {
		diffDay = int(w)
	} else {
		diffDay = 7
	}

	MondayTime = time.Date(t.Year(), t.Month(), t.Day()-diffDay-6, 0, 0, 0, 0, ShanghaiLoc)
	SundayTime = MondayTime.Add(weekDiff)

	return MondayTime, SundayTime
}
