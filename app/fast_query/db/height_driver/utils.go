package height_driver

type Height int64

const (
	// Max int64 value
	maxHeight = 9223372036854775807
)

// Cluster returns /1000 of the height; useful for clustering records in different partitions
func (h Height) Cluster() Height {
	return h / 1000
}

func (h Height) ToInt64() int64 {
	return int64(h)
}

func (h Height) IsLatestHeight() bool {
	if h == maxHeight {
		return true
	} else {
		return false
	}
}

func (h Height) CurrentOrLatest() Height {
	if h == 0 {
		return Height(maxHeight)
	} else {
		return h
	}
}

func (h Height) CurrentOrNever() Height {
	if h == 0 {
		return -1
	} else {
		return h
	}
}

type Key []byte

func (k Key) CurrentOrDefault() []byte {
	if k != nil {
		return k
	} else {
		return []byte{0x0}
	}
}
