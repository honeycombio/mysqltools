package normalizer_test

import (
	"testing"

	"github.com/honeycombio/mysqltools/query/normalizer"
)

var scannerTests = []struct{ ID, Input, Expected string }{
	{"multi-valued select", `SELECT (5, 2, "hi", foo) FROM tablename where id = 5`, "select (?, ?, ?, foo) from tablename where id = ?"},
	{"single order-by", `SELECT colname FROM tablename WHERE id = 5 ORDER BY colname2 ASC`, "select colname from tablename where id = ? order by colname2"},
	{"empty strings work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`text` = \"\"", "select `colname` from `tablename` where `tablename`.`text` = ?"},
	{"escaped quotes", `SELECT colname FROM tablename WHERE text = "an escaped \" doesn't end a string" ORDER BY colname2 ASC`, "select colname from tablename where text = ? order by colname2"},
	{"id literals work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`text` = 'hi there'", "select `colname` from `tablename` where `tablename`.`text` = ?"},
	{"floats work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`float` = 3.14159", "select `colname` from `tablename` where `tablename`.`float` = ?"},
	// fails(issue #1) {"floats without leading 0 work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`float` = .14159", "select `colname` from `tablename` where `tablename`.`float` = ?"},
	{"ints work", "SELECT `colname` FROM `tablename` WHERE `tablename`.`int` = 314159", "select `colname` from `tablename` where `tablename`.`int` = ?"},
	{"alter table", "ALTER TABLE `tablename` ADD COLUMN `text` VARCHAR(100) NOT NULL AFTER `before_text`", "alter table `tablename` add column `text` varchar(?) not null after `before_text`"},
}

func TestScannerNormalization(t *testing.T) {
	n := &normalizer.Scanner{}

	for _, test := range scannerTests {
		actual := n.NormalizeQuery(test.Input)
		if test.Expected != actual {
			t.Error("test '" + test.ID + "' failed normalization.  actual = " + actual)
		}
	}

}
