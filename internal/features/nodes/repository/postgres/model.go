package nodes_postgres_repository

type NodeModel struct {
	ID       int
	LogID    int
	NodeGUID string
	NodeDesc string
	NodeType int
	NumPorts int
}

type NodeInfoModel struct {
	ID                     *int
	NodeID                 *int
	SerialNumber           *string
	PartNumber             *string
	Revision               *string
	ProductName            *string
	Endianness             *int
	EnableEndiannessPerJob *int
	ReproducibilityDisable *int
}
