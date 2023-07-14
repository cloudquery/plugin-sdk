package plugin

import (
	"fmt"
	"time"
)

// tableNameForTest returns a table name that is unique to the test run. It adds the current unix second
// and a random number between 0 and 1000 as a suffix.
func (s *WriterTestSuite) tableNameForTest(name string) string {
	return fmt.Sprintf("cq_%s_test_%d_%04d", name, time.Now().Unix(), s.rand.Intn(1000))
}
