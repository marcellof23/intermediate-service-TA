package boot

var (
	RequestCommand = make(map[string]CountRequest, 0)
)

type CountRequest struct {
	TotalCommand  int
	TotalExecuted int
}
