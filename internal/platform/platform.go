package platform

import "os"

type Diagnostic interface {
	FullDiagnostics(file *os.File)
	BaseDiagnostics(file *os.File)
	NetowrDiagnosics(file *os.File)
}
