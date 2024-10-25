package main

import (
	"github.com/goburrow/modbus"
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

func queryDevice(name, address string, slaveID byte, interval time.Duration, logger *zap.Logger) {
	for {
		handler := modbus.NewTCPClientHandler(address)
		handler.SlaveId = slaveID // Установка Slave ID
		logger.Info("Attempting to connect", zap.String("device", name), zap.String("address", address))

		err := handler.Connect()
		if err != nil {
			logger.Warn("Failed to connect", zap.String("device", name), zap.Error(err))
			handler.Close() // Убедитесь, что соединение закрыто при неудаче
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Info("Connected successfully", zap.String("device", name))

		client := modbus.NewClient(handler)
		defer handler.Close()

		ticker := time.NewTicker(interval)
		for range ticker.C {
			results, err := client.ReadHoldingRegisters(0, 10)
			if err != nil {
				logger.Error("Failed to read registers", zap.String("device", name), zap.Error(err))
				handler.Close()
				break
			} else {
				logger.Info("Read registers", zap.String("device", name), zap.ByteString("data", results))
			}
		}
	}
}

func main() {
	logger := setupLogger()
	defer logger.Sync()

	// Укажите адреса, порты и Slave ID для каждого устройства
	go queryDevice("Modbus1", "192.168.1.64:502", 1, 2*time.Second, logger) // Slave ID 1
	//go queryDevice("Modbus2", "192.168.1.65:502", 2, 5*time.Second, logger) // Slave ID 2

	select {}
}
