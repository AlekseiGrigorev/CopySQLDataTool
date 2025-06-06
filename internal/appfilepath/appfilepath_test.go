package appfilepath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWithDateTimeFileName(t *testing.T) {
	fp := AppFilePath{Path: "test.txt"}
	actual := fp.GetWithDateTime()
	fmt.Println(actual)
	assert.Regexp(t, `^test_\d{8}_\d{6}\.txt$`, actual)
}

func TestGetWithDateTimePath(t *testing.T) {
	fp := AppFilePath{Path: "/var/log/test.txt"}
	actual := fp.GetWithDateTime()
	fmt.Println(actual)
	assert.Regexp(t, `^\Svar\Slog\Stest_\d{8}_\d{6}\.txt$`, actual)
}
