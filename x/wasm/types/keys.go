package types

var (
	executedContractsKey = []byte("terra/executedContracts")
)

func GetExecutedContractsKey() []byte {
	return executedContractsKey
}
