package main

import (
	"github.com/goburrow/modbus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "log"
	"os"
	"time"
)

func setupLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.Create("app.log") // Открыть файл для записи логов
	fileWriter := zapcore.AddSync(logFile)

	core := zapcore.NewCore(fileEncoder, fileWriter, zap.InfoLevel)
	logger := zap.New(core)
	return logger
}

func queryDevice(name string, address string, interval time.Duration, logger *zap.Logger) {
	for {
		handler := modbus.NewTCPClientHandler(address)
		err := handler.Connect()
		if err != nil {
			logger.Warn("Failed to connect", zap.String("device", name), zap.Error(err))
			time.Sleep(5 * time.Second) // Повторное подключение через 5 секунд
			continue
		}

		client := modbus.NewClient(handler)
		defer handler.Close()

		ticker := time.NewTicker(interval)
		for range ticker.C {
			results, err := client.ReadHoldingRegisters(0, 10)
			if err != nil {
				logger.Error("Failed to read registers", zap.String("device", name), zap.Error(err))
				handler.Close() // Закрываем соединение при ошибке
				break           // Переходим к переподключению
			} else {
				logger.Info("Read registers", zap.String("device", name), zap.ByteString("data", results))
			}
		}
	}
}

func main() {
	logger := setupLogger()
	defer logger.Sync()

	go queryDevice("Modbus1", "localhost:5020", 2*time.Second, logger)
	go queryDevice("Modbus2", "localhost:5021", 5*time.Second, logger)

	select {} // Ожидание вечно
}
