package devices

import "github.com/goburrow/modbus"

// DeviceTypeB — структура, представляющая устройство типа B.
type DeviceTypeB struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

// Connect подключается к устройству по указанному адресу и устанавливает Slave ID.
func (d *DeviceTypeB) Connect(address string, slaveID byte) error {
	d.handler = modbus.NewTCPClientHandler(address)
	d.handler.SlaveId = slaveID
	err := d.handler.Connect()
	if err == nil {
		d.client = modbus.NewClient(d.handler)
	}
	return err
}

// ReadData читает регистры устройства типа B.
func (d *DeviceTypeB) ReadData() ([]byte, error) {
	return d.client.ReadInputRegisters(0, 5) // Устройство B читает другие регистры
}

// Disconnect закрывает соединение.
func (d *DeviceTypeB) Disconnect() error {
	return d.handler.Close()
}
