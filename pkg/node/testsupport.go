package node

import (
	"sync"
)

// ImportDB is sync point for tests when importing DB schemas.
var ImportDB = sync.Once{}
