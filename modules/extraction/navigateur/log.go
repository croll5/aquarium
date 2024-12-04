package navigateur

import(
    "time"
)

type Log struct{
    Time_string string
    Time_date time.Time
    Url string
    Title string
    Domain_name string
    Visit_count int
    
}


func (l *Log) ConvertStringToTime() {
    l.Time_date,_ = time.Parse("2006-01-02 15:04:05", l.Time_string)
}