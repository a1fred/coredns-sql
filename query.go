package sql

import "github.com/jinzhu/gorm"

func FindSql(db *gorm.DB, qName string, qType string) ([]*Record, error) {
	var records []*Record
	err := db.Where(
		"LOWER(name)=LOWER(?) AND type=? AND disabled=false",
		qName, qType,
	).Find(&records).Error
	return records, err
}
