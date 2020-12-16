package models

var InsertTopBook = func(b TopBook) {
	b.Insert()
}

func ClearBooks() {
	db.Delete(&TopBook{})
}
