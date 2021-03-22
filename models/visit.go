package models

// Create a db entry
func (visit *Visit) Create() error {
	return db.Create(visit).Error
}

func GetVisits(link_id, start_time string) *[]VisitsPerDay {
	output := &[]VisitsPerDay{}
	db.Raw(`SELECT count(*), date(created_at) FROM visits
		WHERE link_id = ?
		AND
		date(created_at) > ?
		GROUP BY date(created_at);`, link_id, start_time).Scan(&output)
	return output
}

// Update a db entry by ID
func (visit *Visit) Update() error {
	return db.Save(visit).Error
}

// Delete a db entry2
func (visit *Visit) Delete() error {
	return db.Delete(visit).Error
}
