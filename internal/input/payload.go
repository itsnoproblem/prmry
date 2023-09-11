package input

const (
	PayloadTypeText PayloadType = "text"
	//PayloadTypeImage PayloadType = "image"
	//PayloadTypeAudio PayloadType = "audio"
	//PayloadTypeVideo PayloadType = "video"
)

type PayloadType string

type Payload struct {
	Type  PayloadType
	Value []byte
	Tags  map[string]string
}

func Types() []PayloadType {
	return []PayloadType{
		PayloadTypeText,
		//PayloadTypeImage,
		//PayloadTypeAudio,
		//PayloadTypeVideo,
	}
}

func NewTextPayload(value string, tags map[string]string) Payload {
	if tags == nil {
		tags = make(map[string]string)
	}

	return Payload{
		Type:  PayloadTypeText,
		Value: []byte(value),
		Tags:  tags,
	}
}
