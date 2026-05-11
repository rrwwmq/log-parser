package domain

type Port struct {
	ID        int
	NodeID    int
	PortGUID  string
	PortNum   int
	PortState int
	LID       int
}

func NewPort(id int, nodeID int, portGUID string, portNum int, portState int, lid int) Port {
	return Port{
		ID:        id,
		NodeID:    nodeID,
		PortGUID:  portGUID,
		PortNum:   portNum,
		PortState: portState,
		LID:       lid,
	}
}

func NewUninitializedPort(nodeID int, portGUID string, portNum int, portState int, lid int) Port {
	return NewPort(
		UninitializedID,
		nodeID,
		portGUID,
		portNum,
		portState,
		lid,
	)
}
