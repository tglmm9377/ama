package jobstate

type Job struct {
	Name string `xorm:"name pk varchar(20)"`
	Status string `xorm:"status varchar(20)"`
}

