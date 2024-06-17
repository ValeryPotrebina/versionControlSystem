package object

// End leaves of tree (files)
type Blob struct {
	Data []byte
}

// Create object for blob
func (b *Blob) CreateObject() *Object {
	return &Object{
		TypeBlob,
		b.Data,
	}
}


