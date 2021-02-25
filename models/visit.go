package models

// Create a db entry
func (visit *Visit) Create() error {
	return db.Create(visit).Error
}

// Read a db entry by ID
func (visit *Visit) Read() error {
	return db.First(visit).Error
}

// Update a db entry by ID
func (visit *Visit) Update() error {
	return db.Save(visit).Error
}

// Delete a db entry2
func (visit *Visit) Delete() error {
	return db.Delete(visit).Error
}
