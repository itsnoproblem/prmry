package input

const (
	PayloadTypeText  PayloadType = "text"
	PayloadTypeImage PayloadType = "image"
	PayloadTypeAudio PayloadType = "audio"
	PayloadTypeVideo PayloadType = "video"
)

type PayloadType string

type Payload struct {
	Type  PayloadType
	Value []byte
	Tags  map[string]string
}

func SupportedTypes() []PayloadType {
	return []PayloadType{
		PayloadTypeText,
		PayloadTypeImage,
		PayloadTypeAudio,
		PayloadTypeVideo,
	}
}

func NewTextPayload(b []byte) Payload {
	return Payload{
		Type:  PayloadTypeText,
		Value: b,
	}
}
