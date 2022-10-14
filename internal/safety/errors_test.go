package safety

import (
	"errors"
	"fmt"
	"testing"
)

func Test_ErrorReportMessage(t *testing.T) {
	got := createErrorReportMessage("this is the reason", errors.New("this is the error"), []string{"line 1", "line 2"})
	fmt.Println("ERROR Example:")
	fmt.Println(got)
}
