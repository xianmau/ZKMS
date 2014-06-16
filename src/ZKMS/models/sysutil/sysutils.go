package sysutil

type SysUtils struct {
	Cpu  float64 `json:"cpu_util"`
	Disk float64 `json:"disk_util"`
	Net  float64 `json:"net_util"`
}
