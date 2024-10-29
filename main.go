package main

import (
	"GServer/devices"
	_ "github.com/goburrow/modbus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func setupLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.Create("app.log")
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, consoleWriter, zap.InfoLevel),
	)
	logger := zap.New(core)
	return logger
}

func queryDevice(device devices.ModbusDevice, name, address string, slaveID byte, interval time.Duration, logger *zap.Logger) {
	for {
		err := device.Connect(address, slaveID)
		if err != nil {
			logger.Warn("Failed to connect", zap.String("device", name), zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Info("Connected successfully", zap.String("device", name))

		ticker := time.NewTicker(interval)
		for range ticker.C {
			data, err := device.ReadData()
			if err != nil {
				logger.Error("Failed to read data", zap.String("device", name), zap.Error(err))
				device.Disconnect()
				break
			} else {
				logger.Info("Data read", zap.String("device", name), zap.ByteString("data", data))
			}
		}
	}
}

func main() {
	logger := setupLogger()
	defer logger.Sync()

	deviceA := &devices.DeviceTypeA{}
	deviceB := &devices.DeviceTypeB{}

	go queryDevice(deviceA, "DeviceA", "192.168.1.64:502", 1, 2*time.Second, logger)
	go queryDevice(deviceB, "DeviceB", "192.168.1.65:503", 1, 5*time.Second, logger)

	select {}
}
