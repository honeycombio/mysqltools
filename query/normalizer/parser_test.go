package normalizer_test

import (
	"testing"

	"github.com/honeycombio/mysqltools/query/normalizer"
)

var parserTests = []struct{ ID, Input, Expected string }{
	{"multi-valued select", `SELECT (5, 2, "hi", foo) FROM tablename where id = 5`, "select (?, ?, ?, foo) from tablename where id = ?"},
	{"single order-by", `SELECT colname FROM tablename WHERE id = 5 ORDER BY colname2 ASC`, "select colname from tablename where id = ? order by colname2 asc"},
	{"empty strings work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = \"\"", "select colname from tablename where tablename.textcol = ?"},
	{"escaped quotes", `SELECT colname FROM tablename WHERE textCol = "an escaped \" doesn't end a string" ORDER BY colname2 ASC`, "select colname from tablename where textcol = ? order by colname2 asc"},
	{"id literals work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = 'hi there'", "select colname from tablename where tablename.textcol = ?"},
	{"floats work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159", "select colname from tablename where tablename.floatcol = ?"},
	{"floats without leading 0 work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = .14159", "select colname from tablename where tablename.floatcol = ?"},
	{"ints work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159", "select colname from tablename where tablename.intcol = ?"},
	{"union",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159 UNION SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159",
		"select colname from tablename where tablename.intcol = ? union select colname from tablename where tablename.floatcol = ?"},
	{"simple insert",
		"INSERT INTO `tablename` (intCol, floatCol) VALUES (12345, 1.2345)", "insert into tablename(intcol,floatcol) values (?, ?)"},
	{"simple delete",
		"DELETE FROM `tablename` WHERE intCol < 12345", "delete from tablename where intcol < ?"},
	//{"alter table", "ALTER TABLE `tablename` ADD COLUMN `text` VARCHAR(100) NOT NULL AFTER `before_text`", "alter table tablename add column text varchar(?) not null after before_text"},
}

func TestParserNormalization(t *testing.T) {
	n := &normalizer.Parser{}

	for _, test := range parserTests {
		actual := n.NormalizeQuery(test.Input)
		if test.Expected != actual {
			t.Error("test '" + test.ID + "' failed normalization.  actual = " + actual)
		}
	}

}
