package crossfunction

// DBClient to work with DBWrapper
type DBClient interface {
	SaveToDB(interface{})
}
