package s32cs_test

import (
	"os"
	"testing"

	"github.com/fujiwara/s32cs"
)

func TestApexHandler(t *testing.T) {
	os.Setenv("AWS_EXECUTION_ENV", "Test_AWS_Lambda_go")
	defer os.Setenv("AWS_EXECUTION_ENV", "")
	s32cs.ApexRun()
}
