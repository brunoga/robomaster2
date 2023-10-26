package dji

type DJIFunctionEnableParamValue struct {
	List []DJIFunctionEnableInfo `json:"list"`
}

func NewDJIFunctionEnableParamValue() *DJIFunctionEnableParamValue {
	return &DJIFunctionEnableParamValue{
		List: []DJIFunctionEnableInfo{},
	}
}

func (p *DJIFunctionEnableParamValue) AddFunctionType(id DJIRMEnableFunctionType, enabled bool) {
	for _, info := range p.List {
		if info.ID == id {
			info.Enabled = enabled
			return
		}
	}

	p.List = append(p.List, DJIFunctionEnableInfo{
		ID:      id,
		Enabled: enabled,
	})
}
