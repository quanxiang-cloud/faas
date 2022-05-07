package consts

// status of function doc
const (
	DocNotExists = iota
	DocTaskExists
	DocSucceed
	DocFailed
)

var DocStatusMapping = map[string]int{
	"False": DocFailed,
	"True":  DocSucceed,
}
