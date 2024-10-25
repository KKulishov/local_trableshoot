package platform

import "os"

type Diagnostic interface {
	BaseDiagnostics(file *os.File)
}
