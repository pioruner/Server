package devices

// ModbusDevice описывает интерфейс для всех Modbus устройств.
type ModbusDevice interface {
	Connect(address string, slaveID byte) error
	ReadData() ([]byte, error)
	Disconnect() error
}
