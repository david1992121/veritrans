package veritrans

import (
	"time"
)

// get the time after one month from now
func getAfterOneMonth() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(0, 1, 0)
	return expiredAt.Format("01/06")
}

// get the time after one year from now
func getAfterOneYear() string {
	nowTime := time.Now()
	expiredAt := nowTime.AddDate(1, 0, 0)
	return expiredAt.Format("01/06")
}
