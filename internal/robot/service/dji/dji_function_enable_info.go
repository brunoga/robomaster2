package dji

type DJIFunctionEnableInfo struct {
	ID      DJIRMEnableFunctionType `json:"id"`
	Enabled bool                    `json:"enabled"`
}
