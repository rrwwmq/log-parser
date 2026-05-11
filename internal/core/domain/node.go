package domain

type NodeType int

var (
	NodeTypeHost   NodeType = 1
	NodeTypeSwitch NodeType = 2
)

type Node struct {
	ID       int
	LogID    int
	NodeGUID string
	NodeDesc string
	NodeType NodeType
	NumPorts int
	Info     *NodeInfo
}

type NodeInfo struct {
	ID                     int
	NodeID                 int
	SerialNumber           *string
	PartNumber             *string
	Revision               *string
	ProductName            *string
	Endianness             *int
	EnableEndiannessPerJob *int
	ReproducibilityDisable *int
}

func NewNode(id int, logID int, nodeGUID string, nodeDesc string, nodeType NodeType, numPorts int, info *NodeInfo) Node {
	return Node{
		ID:       id,
		LogID:    logID,
		NodeGUID: nodeGUID,
		NodeDesc: nodeDesc,
		NodeType: nodeType,
		NumPorts: numPorts,
		Info:     info,
	}
}

func NewUninitializedNode(logID int, nodeGUID string, nodeDesc string, nodeType NodeType, numPorts int) Node {
	return NewNode(
		UninitializedID,
		logID,
		nodeGUID,
		nodeDesc,
		nodeType,
		numPorts,
		nil,
	)
}
