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

func (l *Log) Reverse_domain(){
  var s = l.Domain_name
  rs := []rune(s)
  for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
    rs[i], rs[j] = rs[j], rs[i]
  }
  l.Domain_name = string(rs) 
}