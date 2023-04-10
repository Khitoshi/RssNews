package tables

import "rss_reader/database"

func CreateTable() error {

	var err error

	if err = database.CreateTable((*ITEMS)(nil)); err != nil {
		return err
	}

	if err = database.CreateTable((*RSS_URLS)(nil)); err != nil {
		return err
	}

	if err = database.CreateTable((*USER_FAVORITE_ITEMS)(nil)); err != nil {
		return err
	}

	if err = database.CreateTable((*USER_ITEMS)(nil)); err != nil {
		return err
	}

	if err = database.CreateTable((*USER)(nil)); err != nil {
		return err
	}

	return nil
}
