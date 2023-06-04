package util

// post process

type Output struct {
	Average      int
	OutputSuffix string
	OutputPrefix string               // dictionary
	OutputFormat []OutputFormatSample // 支持多种输出
}

type OutputFormatSample struct {
	OutputFormat string
	Start        int
	End          int
}

type Input struct {
	Start       int
	End         int
	InputType   string // heat / cool
	DataType    string // heat / energy
	InputPrefix string // dictionary
	InputSuffix string
}

type Data interface{}

type HeatData struct {
	TimeStep int64
	RouV     float64
	TempNow  float64
	Temp     float64
	Volume   float64
	Density  float64
	Enthalpy float64
}

type EnergyData struct {
	TimeStep int64
	Ke       float64
	Pe       float64
	ETotal   float64
	EMol     float64
	EPair    float64
	ELong    float64
}

type HEData struct {
	TimeStep int64
	RouV     float64
	TempNow  float64
	Temp     float64
	Volume   float64
	Density  float64
	Enthalpy float64
	Ke       float64
	Pe       float64
	ETotal   float64
	EMol     float64
	EPair    float64
	ELong    float64
}

// pre process

type PreInput struct {
	Path      string
	InputType string
}
