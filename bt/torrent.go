package bt

// InfoHash type ------------------------------------------------------

type InfoHash string

func (hash InfoHash) Validate() bool {
	return len(hash) == 20
}
