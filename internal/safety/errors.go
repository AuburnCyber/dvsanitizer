package safety

import (
	"log"
	"runtime"
)

const (
	ERROR_REPORT_MSG = "" +
		`################################################################################

An error has been encountered and the sanitizer can not proceed safely. The
above data describes the situation encountered and if the solution is
non-obvious, please copy-paste everything between the "BEGIN ERROR REPORT" and
"END ERROR REPORT" lines into an email and send to us at team@DVSorder.org so
that we can debug and correct. In addition to fixing for you, we'd like to fix
it for any other jurisdiction before they encounter it.

Some of the reasons this error may have encountered are:
		- A data format was encountered that we have not seen before
		- An combination of facts which was thought to not exist but found
		- A specific-check failed and while it may be safe to procede, we'd
		  prefer to confirm with you directly rather than try to explain via an
		  error message

You should review the data to ensure that you are comfortable sending it to us
and if not, please send us what you are comfortable sending and the fields
which you are not.

-DVSorder Team

################################################################################
`
)

var (
	storedErrDescs []string
)

// Due to the likelihood of having to debug misunderstood/never-seen data
// formats remotely (not to mention errors), this is a mechanism to immediately
// stop processing if not-known-good things occur. The goal is to create a
// straight-forward bug report that people can copy-paste to us while
// simultaneously:
// 		A) provide enough information to debug/understand what happened
//
//		B) Be in a readable format and understandable presentation that
//		non-technical people can look at it and feel confident that it's not
//		trying to leak any secrets like the seed or other information that
//		shouldn't be public
func ReportError(reason string, triggerErr error, descLines ...string) {
	msg := createErrorReportMessage(reason, triggerErr, descLines)
	log.Fatal(msg)
}

func createErrorReportMessage(reason string, triggerErr error, descLines []string) string {
	separator := "--------------------------------------------------------------------------------\n"

	log.Printf("error encountered, printing report")

	// Write details first b/c they are unconstrained lengths
	msg := "------------------------------ BEGIN ERROR REPORT ------------------------------\n"
	if len(descLines) != 0 || len(storedErrDescs) != 0 {
		msg += "SPECIFIC DETAILS:\n"
		for _, line := range storedErrDescs {
			msg += "\t" + line + "\n"
		}
		for _, line := range descLines {
			msg += "\t" + line + "\n"
		}
	} else {
		msg += "SPECIFIC DETAILS: None\n"
	}

	// Write the stack-trace next b/c if details explode, still want to know
	// the code-path of how arrived at the error
	msg += separator
	msg += "STACK TRACE:\n"
	stackBuffer := make([]byte, 1024*1024*1024, 1024*1024*1024) // Far larger than should ever need for safety and will be truncated.
	n := runtime.Stack(stackBuffer, false)
	msg += string(stackBuffer[:n])

	// Write the explicit reasons next so that they're right above the message
	msg += separator
	msg += "REASON: " + reason + "\n"
	if triggerErr == nil {
		msg += "INTERNAL ERROR OBJECT: None\n"
	} else {
		msg += "INTERNAL ERROR OBJECT: " + triggerErr.Error() + "\n"
	}

	// Write the message last so it's on-screen
	msg += ERROR_REPORT_MSG
	msg += "------------------------------- END ERROR REPORT -------------------------------\n"

	return msg
}

// Due to the likelihood of having to debug misunderstood/never-seen data
// formats remotely (not to mention errors), this creates a list of of things
// that might be useful to know so that it's A) all in 1 place and B) obvious
// that people need to send it to us. These will be included in the Error
// Report created above.
func StoreErrorDesc(line string) {
	storedErrDescs = append(storedErrDescs, line)
}
