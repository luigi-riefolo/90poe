package store

// TODO:
// get psql and mongo docker
// let env decide which store
// initialise clients with client lib

/*
   Postgres has a strongly typed schema that leaves very little room for errors. You first create the schema for a table and then add rows to the table. You can also define relationships between different tables with rules so that you can store related data across several tables and avoid data duplication. All this means someone on the team can act as a 'Database architect' and control the schema which acts as a standard for everyone else to follow.
   You can change tables in PostgreSQL without requiring to lock it for every operation. For example, you can add a column and set a default value quickly without locking the entire table. This ensures that every row in a table has the column and your codebase remains clean without needing to check if the column exists at every stage. It is also much quicker to update every row since Postgres doesn't need to retrieve each row, update, and put it back.
   Postgres also supports JSONB, which lets you create unstructured data, but with data constraint and validation functions to help ensure that JSON documents are more meaningful. The folks at Sisense have written a great blog with a detailed comparison of Postgres vs MongoDB for JSON documents.
   Our database size reduced by 10x since Postgres stores information more efficiently and data isn't unnecessarily duplicated across tables.
   As was shown in previous studies, we found that Postgres performed much better for indexes and joins and our service became faster and snappier as a result.
*/

// PostgresStore represents a PostgreSQL store
// that implements the the Storer interface.
type PostgresStore struct {
}

// NewPostresStore returns a PostgresStore client.
func NewPostresStore() *PostgresStore {
	store := &PostgresStore{}

	return store
}

// GetEntry returns an Entry by ID.
func (p *PostgresStore) GetEntry(ID string) Entry {

	return Entry{}
}

// UpsertEntry inserts a new Entry or updates an existing one.
func (p *PostgresStore) UpsertEntry(entry Entry) Entry {

	return Entry{}
}
