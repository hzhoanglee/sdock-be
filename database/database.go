package database

import (
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

var RethinkDB *RdbSess
