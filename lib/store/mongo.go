package store

// MongoStore represents a MongoSQL store
// that implements the the Storer interface.
type MongoStore struct {
}

// NewMongoStore returns a MongoStore client.
func NewMongoStore() *MongoStore {
	store := &MongoStore{}

	return store
}

// GetEntry returns an Entry by ID.
func (p *MongoStore) GetEntry(ID string) Entry {

	return Entry{}
}

// UpsertEntry inserts a new Entry or updates an existing one.
func (p *MongoStore) UpsertEntry(entry Entry) Entry {

	return Entry{}
}
