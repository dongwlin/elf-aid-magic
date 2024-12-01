package logic

import "os"

type PidLogic struct{}

func NewPidLogic() *PidLogic {
	return &PidLogic{}
}

func (p *PidLogic) GetPid() int {
	return os.Getpid()
}

func (p *PidLogic) ValidatePid(pid int) bool {
	return p.GetPid() == pid
}
