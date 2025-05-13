package utils

import "time"

// GetZeroUnix 获取北京时间0点时间戳，精确到秒
func GetZeroUnix() int64 {
	t := time.Now()
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	addTime := time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), 0, 0, 0, 0, loc)
	timeSamp := addTime.Unix()
	return timeSamp
}

// GetTimeZeroUnix 获取北京时间0点时间戳，精确到秒
func GetTimeZeroUnix(t time.Time) int64 {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	addTime := time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), 0, 0, 0, 0, loc)
	timeSamp := addTime.Unix()
	return timeSamp
}

// ConvertToBJTime 转成北京时间
func ConvertToBJTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	return time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), t.In(loc).Hour(), t.In(loc).Minute(), t.In(loc).Second(), t.In(loc).Nanosecond(), loc)
}

// GetTimeZeroUnixMilli 获取北京时间0点时间戳，精确到毫秒
func GetTimeZeroUnixMilli(t time.Time) int64 {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	addTime := time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), 0, 0, 0, 0, loc)
	timeSamp := addTime.UnixMilli()
	return timeSamp
}

func GetTimeStr() string {
	t := time.Now()
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	date := time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), t.In(loc).Hour(), t.In(loc).Minute(), t.In(loc).Second(), 0, loc)
	return date.Format("2006-01-02 15:04:05")
}

// TimeStampToDate 时间戳转日期
func TimeStampToDate(timeStamp int64) string {
	timeLayout := "2006-01-02"
	return time.Unix(timeStamp, 0).Format(timeLayout)
}

// DateToTimeStamp 日期转时间戳
func DateToTimeStamp(timeStr string) (int64, error) {

	timeLayout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	_, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return 0, err
	}
	timeStamp, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	return timeStamp.Unix(), nil
}
