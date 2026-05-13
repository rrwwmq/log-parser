package ports_postgres_repository

type PortModel struct {
	ID        int
	NodeID    int
	PortGUID  string
	PortNum   int
	PortState int
	LID       int
}
