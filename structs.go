package trojanTechTickets

import "time"

type BaseTemplateData struct {
	CurrentURL  string
	LogoutURL   string
	User        string
	UserIsAdmin bool
}

type Customer struct {
	Email string
	Grade int
	Name  string
}

type Step struct {
	Action     string
	Diagnostic string
}

type Ticket struct {
	Category     string
	ClassName    string
	ClassType    string
	ClassTeacher string
	Customer     Customer
	Description  string
	EscalatedTo  string
	Number       int
	Steps        []Step
	Time         time.Time
	Worker       string
}
