package entity

const (
	SHAREFORCELOGOUT = "SHAREFORCELOGOUT"
	GRAYDATAUPDATE   = "GRAYDATAUPDATE"
	FINDCONNECTION   = "FINDCONNECTION"
)

type NodeEntity struct {
	NodeName string
	Ip       string
	Port     int
}

type NodeMessageEntity struct {
	Context ContextEntity `json:"context"`
	Body    string        `json:"body"`
	Type    string        `json:"type"`
	Async   bool          `json:"async"`
}
