package devices

import "github.com/goburrow/modbus"

// DeviceTypeA — структура, представляющая устройство типа A.
type DeviceTypeA struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

// Connect подключается к устройству по указанному адресу и устанавливает Slave ID.
func (d *DeviceTypeA) Connect(address string, slaveID byte) error {
	d.handler = modbus.NewTCPClientHandler(address)
	d.handler.SlaveId = slaveID
	err := d.handler.Connect()
	if err == nil {
		d.client = modbus.NewClient(d.handler)
	}
	return err
}

// ReadData читает регистры устройства типа A.
func (d *DeviceTypeA) ReadData() ([]byte, error) {
	return d.client.ReadHoldingRegisters(0, 10)
}

// Disconnect закрывает соединение.
func (d *DeviceTypeA) Disconnect() error {
	return d.handler.Close()
}
