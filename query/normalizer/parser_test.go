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
}{
	{"multi-valued select",
		`SELECT (5, 2, "hi", foo) FROM tablename where id = 5`,
		"select (?, ?, ?, foo) from tablename where id = ?",
		"select",
		[]string{"tablename"},
	},
	{"single order-by",
		`SELECT colname FROM tablename WHERE id = 5 ORDER BY colname2 ASC`,
		"select colname from tablename where id = ? order by colname2 asc",
		"select",
		[]string{"tablename"},
	},
	{"empty strings work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = \"\"",
		"select colname from tablename where tablename.textcol = ?",
		"select",
		[]string{"tablename"},
	},
	{"escaped quotes",
		`SELECT colname FROM tablename WHERE textCol = "an escaped \" doesn't end a string" ORDER BY colname2 ASC`,
		"select colname from tablename where textcol = ? order by colname2 asc",
		"select",
		[]string{"tablename"},
	},
	{"id literals work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`textCol` = 'hi there'",
		"select colname from tablename where tablename.textcol = ?",
		"select",
		[]string{"tablename"},
	},
	{"floats work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159",
		"select colname from tablename where tablename.floatcol = ?",
		"select",
		[]string{"tablename"},
	},
	{"floats without leading 0 work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = .14159",
		"select colname from tablename where tablename.floatcol = ?",
		"select",
		[]string{"tablename"},
	},
	{"ints work",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159",
		"select colname from tablename where tablename.intcol = ?",
		"select",
		[]string{"tablename"},
	},
	{"union",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159 UNION SELECT `colname` FROM `tablename2` WHERE `tablename2`.`floatCol` = 3.14159",
		"select colname from tablename where tablename.intcol = ? union select colname from tablename2 where tablename2.floatcol = ?",
		"union",
		[]string{"tablename", "tablename2"},
	},
	{"union single table",
		"SELECT `colname` FROM `tablename` WHERE `tablename`.`intCol` = 314159 UNION SELECT `colname` FROM `tablename` WHERE `tablename`.`floatCol` = 3.14159",
		"select colname from tablename where tablename.intcol = ? union select colname from tablename where tablename.floatcol = ?",
		"union",
		[]string{"tablename"},
	},
	{"simple insert",
		"INSERT INTO `tablename` (intCol, floatCol) VALUES (12345, 1.2345)",
		"insert into tablename(intcol,floatcol) values (?, ?)",
		"insert",
		[]string{"tablename"},
	},
	{"insert with subquery",
		"INSERT INTO `tablename` (intCol, floatCol) SELECT intCol2, floatCol2 FROM sourceTable WHERE id = 12345",
		"insert into tablename(intcol,floatcol) select intcol2,floatcol2 from sourcetable where id = ?",
		"insert",
		[]string{"sourceTable", "tablename"},
	},
	{"simple delete",
		"DELETE FROM `tablename` WHERE intCol < 12345",
		"delete from tablename where intcol < ?",
		"delete",
		[]string{"tablename"},
	},
	{"inner join",
		"SELECT `colname` FROM `tablename` INNER JOIN `tablename2` ON `tablename`.`colName` = `tablename2`.`colName2` WHERE `tablename`.`intCol` = 314159",
		"select colname from tablename join tablename2 on tablename.colname = tablename2.colname2 where tablename.intcol = ?",
		"select",
		[]string{"tablename", "tablename2"},
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
	}

}
