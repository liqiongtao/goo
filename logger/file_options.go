package logger

import "time"

type FileOption struct {
	Name  string
	Value interface{}
}

const (
	optNameFilePath       = "file-path"
	optNameFileName       = "file-name"
	optNameFileDateFormat = "file-date-format"
	optNameFileMaxAge     = "file-max-age"
	optNameFileMaxSize    = "file-max-size"
)

func FilePathOption(filePath string) FileOption {
	return FileOption{Name: optNameFilePath, Value: filePath}
}

func FileNameOption(fileName string) FileOption {
	return FileOption{Name: optNameFileName, Value: fileName}
}

func FileDateFormatOption(fileDateFormat string) FileOption {
	return FileOption{Name: optNameFileDateFormat, Value: fileDateFormat}
}

func FileMaxAgeOption(fileMaxAge time.Duration) FileOption {
	return FileOption{Name: optNameFileMaxAge, Value: fileMaxAge}
}

func FileMaxSizeOption(fileMaxSize int64) FileOption {
	return FileOption{Name: optNameFileMaxSize, Value: fileMaxSize}
}
