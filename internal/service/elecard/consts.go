package elecard

var (
	XmlHeader = []byte(`<?xml version="1.0" encoding="windows-1251"?>` + "\n")
)

const (
	ConfiguringStatus   = "Configuring"
	RunningStatus       = "Running"
	WaitingStatus       = "Waiting"
	SuccessStatus       = "Success"
	FaultStatus         = "Fault"
	CriticalErrorStatus = "Critical Error"
)
