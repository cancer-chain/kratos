package chaindb

import (
	"strings"
	"testing"
)

func TestErrSql(t *testing.T) {
	s := `INSERT INTO coins (max, init, supply,supply_float, module, symbol, \n\t\t\t\tcan_lock, creator, issue_create_height, height, can_issue, issue_to_height, _desc ,time) VALUES  ('10000000000000000000000000000000kratos/kts','0kratos/kts',0,0,'asset','kts','true','kratos','1',0,'true','0','for staking','0001-01-01 00:00:00')`
	QuerySec := "'" + strings.ReplaceAll(s, "'", "''") + "'" + ","

	t.Logf("%s", QuerySec)
}
