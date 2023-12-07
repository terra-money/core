package height_driver

import tmdb "github.com/cometbft/cometbft-db"

var _ tmdb.Iterator = (*HeightDBIterator)(nil)

type HeightDBIterator struct {
	oit      tmdb.Iterator
	atHeight int64
}

func NewHeightLimitedIterator(atHeight int64, oit tmdb.Iterator) tmdb.Iterator {
	return &HeightDBIterator{
		oit:      oit,
		atHeight: atHeight,
	}
}

func (h *HeightDBIterator) Domain() (start []byte, end []byte) {
	// TODO: fix me
	return h.oit.Domain()
}

func (h *HeightDBIterator) Valid() bool {
	return h.oit.Valid()
}

func (h *HeightDBIterator) Next() {
	h.oit.Next()
}

func (h *HeightDBIterator) Key() (key []byte) {
	return h.oit.Key()[:len(key)-9]
}

func (h *HeightDBIterator) Value() (value []byte) {
	return h.oit.Value()
}

func (h *HeightDBIterator) Error() error {
	return h.oit.Error()
}

func (h *HeightDBIterator) Close() error {
	return h.oit.Close()
}
