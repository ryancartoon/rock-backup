package host

type Host struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name     string `gorm:"column:name"`
	Location string
	IsActive bool
	Load     int
}

func TableName() string {
	return "hosts"
}
