package normalizer_test

import (
	"fmt"
	"testing"

	"github.com/honeycombio/mysqltools/query/normalizer"
)

var parserTests = []struct {
	ID                string
	Input             string
	ExpectedOutput    string
	ExpectedStatement string
	ExpectedTables    []string
	ExpectedComments  []string
}{
	{"multi-valued select",
		`SELECT (5, 2, "hi", foo) FROM tablename where id = 5`,
		"select (?, ?, ?, foo) from tablename where id = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"single order-by",
		`SELECT colname FROM tablename WHERE id = 5 ORDER BY colname2 ASC`,
		"select colname from tablename where id = ? order by colname2",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"empty strings work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = \"\"",
		"select `colname` from `tablename` where `tablename`.`textcol` = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"escaped quotes",
		`SELECT colname FROM tablename WHERE textCol = "an escaped \" doesn't end a string" ORDER BY colname2 ASC`,
		"select colname from tablename where textcol = ? order by colname2",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"id literals work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = 'hi there'",
		"select `colname` from `tablename` where `tablename`.`textcol` = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"floats work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159",
		"select `colname` from `tablename` where `tablename`.`floatcol` = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"floats without leading 0 work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = .14159",
		"select `colname` from `tablename` where `tablename`.`floatcol` = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"ints work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159",
		"select `colname` from `tablename` where `tablename`.`intcol` = ?",
		"select",
		[]string{"tablename"},
		[]string{},
	},
	{"union",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159 UNION SELECT `colname` FROM `tablename2` WHERE `tablename2`.`floatCol` = 3.14159",
		"select `colname` from `tablename` where `tablename`.`intcol` = ? union select `colname` from `tablename2` where `tablename2`.`floatcol` = ?",
		"union",
		[]string{"tablename", "tablename2"},
		[]string{},
	},
	{"union single table",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159 UNION SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159",
		"select `colname` from `tablename` where `tablename`.`intcol` = ? union select `colname` from `tablename` where `tablename`.`floatcol` = ?",
		"union",
		[]string{"tablename"},
		[]string{},
	},
	{"simple insert",
		"INSERT /* insert comment here */ INTO `tablename` (intCol, floatCol) VALUES (12345, 1.2345)",
		"insert into `tablename`(intcol,floatcol) values (?, ?)",
		"insert",
		[]string{"tablename"},
		[]string{"insert comment here"},
	},
	{"insert with subquery",
		"INSERT /* comment1 */ /* comment2 */ INTO `tablename` (intCol, floatCol) SELECT /* comment3 */ intCol2, floatCol2 FROM sourceTable WHERE id = 12345",
		"insert into `tablename`(intcol,floatcol) select intcol2,floatcol2 from sourcetable where id = ?",
		"insert",
		[]string{"sourcetable", "tablename"},
		[]string{"comment1", "comment2", "comment3"},
	},
	{"insert with subquery and constants",
		"INSERT INTO `tablename` (col1, col2, col3) SELECT 1, 2, intCol2, floatCol2 FROM sourceTable WHERE id = 12345",
		"insert into `tablename`(col1,col2,col3) select ?,?,intcol2,floatcol2 from sourcetable where id = ?",
		"insert",
		[]string{"sourcetable", "tablename"},
		[]string{},
	},
	{"simple delete",
		"DELETE FROM `tablename` WHERE intCol < 12345",
		"delete from `tablename` where intcol < ?",
		"delete",
		[]string{"tablename"},
		[]string{},
	},
	{"inner join",
		"SELECT `colname` FROM `tablename` INNER JOIN `tablename2` ON `tablename`.`colName` = `tablename2`.`colName2` WHERE `tablename`.`intCol` = 314159",
		"select `colname` from `tablename` inner join `tablename2` on `tablename`.`colname` = `tablename2`.`colname2` where `tablename`.`intcol` = ?",
		"select",
		[]string{"tablename", "tablename2"},
		[]string{},
	},
	{"parse error falls back to scan normalizer",
		"SELECT `colname` FROM `tablename` INNER JOIN `tablename2` ON `tablename`.`colName` = `tablename2`.`colName2` WHERE `tablename`.`intCol` = 314159 ORDER BY date",
		"select `colname` from `tablename` inner join `tablename2` on `tablename`.`colname` = `tablename2`.`colname2` where `tablename`.`intcol` = ? order by date",
		"",
		[]string{},
		[]string{},
	},
	//{"alter table", "ALTER TABLE `tablename` ADD COLUMN `text` VARCHAR(100) NOT NULL AFTER `before_text`", "alter table tablename add column text varchar(?) not null after before_text"},
}

func TestParserNormalization(t *testing.T) {
	n := &normalizer.Parser{}

	for _, test := range parserTests {
		actual := n.NormalizeQuery(test.Input)
		if test.ExpectedOutput != actual {
			t.Error("test '" + test.ID + "' failed normalization.  actual = " + actual)
		}

		if test.ExpectedStatement != n.LastStatement {
			t.Error("test '" + test.ID + "' failed statement.  actual = " + n.LastStatement)
		}

		tables := n.LastTables
		tablesEqual := true
		if len(test.ExpectedTables) != len(tables) {
			tablesEqual = false
		}
		if tablesEqual {
			for i := 0; i < len(tables); i++ {
				if tables[i] != test.ExpectedTables[i] {
					tablesEqual = false
					break
				}
			}
		}
		if !tablesEqual {
			t.Error("test '" + test.ID + "' failed table accumulation.  actual = " + fmt.Sprint(tables))
		}

		comments := n.LastComments
		commentsEqual := true
		if len(test.ExpectedComments) != len(comments) {
			commentsEqual = false
		}
		if commentsEqual {
			for i := 0; i < len(comments); i++ {
				if comments[i] != test.ExpectedComments[i] {
					commentsEqual = false
					break
				}
			}
		}
		if !commentsEqual {
			t.Error("test '" + test.ID + "' failed comment accumulation.  actual = " + fmt.Sprint(comments))
		}
	}

}
