package internal

type PhotoVariant struct {
	Url      string `json:"url"`
	SizeType string `json:"type"` // about types: https://dev.vk.com/reference/objects/photo-sizes
	Height   uint32 `json:"height"`
	Width    uint32 `json:"width"`
}

type Variants []PhotoVariant

var originalSizes = map[string]int{
	"s": 1,
	"m": 2,
	"x": 3,
	"y": 4,
	"z": 5,
	"w": 6,
}

func (vv Variants) GetMaxSize() PhotoVariant {
	var maxPhotoVariant PhotoVariant
	maxSize := originalSizes["s"]
	for _, variant := range vv {
		if originalSizes[variant.SizeType] >= maxSize {
			maxPhotoVariant = variant
			maxSize = originalSizes[variant.SizeType]
		}
	}

	return maxPhotoVariant
}

func (vv Variants) GetMinSize() PhotoVariant {
	var minPhotoVariant PhotoVariant
	minSize := originalSizes["w"]
	for _, variant := range vv {
		currentSize := originalSizes[variant.SizeType]
		if currentSize <= minSize && currentSize != 0 {
			minPhotoVariant = variant
			minSize = originalSizes[variant.SizeType]
		}
	}

	return minPhotoVariant
}
