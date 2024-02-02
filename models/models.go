package models

type Agent struct {
    ID       string `gorm:"primaryKey"`
    Hostname string
    HostOS   string
    IP       string
    Status   string
}
