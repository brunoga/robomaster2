package event

// Type represents the specific Unity event type.
type Type uint64

const (
	TypeSetValue           = Type(0)
	TypeGetValue           = Type(1)
	TypeGetAvailableValue  = Type(2)
	TypePerformAction      = Type(3)
	TypeStartListening     = Type(4)
	TypeStopListening      = Type(5)
	TypeActivation         = Type(6)
	TypeLocalAlbum         = Type(7)
	TypeFirmwareUpgrade    = Type(8)
	TypeConnection         = Type(100)
	TypeSecurity           = Type(101)
	TypePrintLog           = Type(200)
	TypeStartVideo         = Type(300)
	TypeStopVideo          = Type(301)
	TypeRender             = Type(302)
	TypeGetNativeTexture   = Type(303)
	TypeVideoTransferSpeed = Type(304)
	TypeAudioDataRecv      = Type(305)
	TypeVideoDataRecv      = Type(306)
	TypeNativeFunctions    = Type(500)
)

var typeNameMap = map[Type]string{
	0:   "TypeSetValue",
	1:   "TypeGetValue",
	2:   "TypeGetAvailableValue",
	3:   "TypePerformAction",
	4:   "TypeStartListening",
	5:   "TypeStopListening",
	6:   "TypeActivation",
	7:   "TypeLocalAlbum",
	8:   "TypeFirmwareUpgrade",
	100: "TypeConnection",
	101: "TypeSecurity",
	200: "TypePrintLog",
	300: "TypeStartVideo",
	301: "TypeStopVideo",
	302: "TypeRender",
	303: "TypeGetNativeTexture",
	304: "TypeVideoTransferSpeed",
	305: "TypeAudioDataRecv",
	306: "TypeVideoDataRecv",
	500: "TypeNativeFunctions",
}

// IsValidType checks if the given Type is valid. It returns true if
// it is and false oherwise.
func IsValidType(typ Type) bool {
	_, ok := typeNameMap[typ]

	return ok
}

// TypeName returns the name associated with the given Type. If it
// is not known, returns an empty string.
func TypeName(typ Type) string {
	typeName, ok := typeNameMap[typ]
	if !ok {
		return ""
	}

	return typeName
}

func AllTypes() []Type {
	types := make([]Type, 0, len(typeNameMap))
	for typ, _ := range typeNameMap {
		types = append(types, typ)
	}

	return types
}
