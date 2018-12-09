package store

// Storer interface describes the methods
// available for a storer implementation.
type Storer interface {

	// GetEntry returns an Entry by ID.
	GetEntry(ID string) Entry

	// UpsertEntry inserts a new Entry or updates an existing one.
	UpsertEntry(entry Entry) Entry
}
