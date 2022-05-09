package messages

const (
	// handlers
	ErrUnsupportedHttpMethod = "Unsupported http method: %s."
	ErrReqRead               = "Error reading request body: %v."
	ErrReqParse              = "Error parsing request body: %v."
	ErrCloseBody             = "Error closing request body: %v"

	ErrInvalidResource = "invalid resource format"

	//dbs
	ErrDbCreateTable = "Error while creating %s table: %v."
	ErrDbRead        = "Reading from DB error: %v."
	ErrDbAdd         = "Adding in DB error: %v."
	ErrDbDelete      = "Deleting from DB error: %v."
)
