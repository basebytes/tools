package tools

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
)

const (
	BoxSizeRateMin           float32 = 0.1
	ImageQueryBoxSizeRateMin float32 = 0.1
)

var u64Arr = [64]uint64{}

func ImgBytesDecode(bs []byte) (image.Image, error) {
	fp := bytes.NewBuffer(bs)
	return jpeg.Decode(fp)
}

func u64Allow(base, v uint64) bool {
	return (v & base) == v
}

func ImageHashNull(a uint64) uint8 {
	i := uint8(0)
	for _, v := range u64Arr {
		if u64Allow(a, v) {
			i += 1
		}
	}
	if i > 32 {
		i = 64 - i
	}
	return i
}

func ImageHashContentRate(a uint64) float32 {
	i := float32(0)
	for _, v := range u64Arr {
		if u64Allow(a, v) {
			i += 1
		}
	}
	if i > 32 {
		i = 64 - i
	}
	return i / 64
}
func ImageMaskAllow(masks [][]float32, x, y int, w, h int) bool {

	if len(masks) < 1 || len(masks) != len(masks[0]) {
		return false
	}

	mn := float32(len(masks))

	xn := int(float32(x) / (float32(w) / mn))
	yn := int(float32(y) / (float32(h) / mn))

	if masks[yn][xn] >= 0.1 {
		return true
	}

	return false
}

func ImagePixGray(im image.Image, x, y int) uint8 {
	p := im.At(x, y)
	r, g, b, _ := p.RGBA()
	gc := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
	return uint8(gc / 4)
}

func ImageHashId(im image.Image, bs []byte) (image.Image, image.Image, uint64, error) {

	var err error

	if im == nil {
		r := bytes.NewBuffer(bs)

		im, err = jpeg.Decode(r)
		if err != nil {
			return nil, nil, 0, err
		}
	}

	var (
		// width   = im.Bounds().Max.X
		// height  = im.Bounds().Max.Y
		i8w     = int(8)
		im8     = resize.Resize(uint(i8w), uint(i8w), im, resize.Lanczos3)
		im8g    = image.NewGray(image.Rect(0, 0, i8w, i8w))
		hashs   = []uint8{}
		hashAll = uint64(0)
		hashId  = uint64(0)
	)

	//

	for x := 0; x < i8w; x++ {
		for y := 0; y < i8w; y++ {
			cg := ImagePixGray(im8, x, y)
			hashs = append(hashs, cg)
			hashAll += uint64(cg)
			im8g.SetGray(x, y, color.Gray{Y: cg * 4})
		}
	}

	hashAvg := uint8(hashAll / uint64(i8w*i8w))
	for x := 0; x < i8w; x++ {
		for y := 0; y < i8w; y++ {
			if hashs[(x*i8w)+y] >= hashAvg {
				//
			} else {
				hashId = hashId | (1 << uint64((x*i8w)+y))
			}
		}
	}

	return im, im8, hashId, nil
}

func ImageHashDiff(a, b uint64) float32 {

	score := float32(64)

	for _, v := range u64Arr {

		if u64Allow(a, v) != u64Allow(b, v) {
			score -= 1

			if score < 32 {
				break
			}
		}
	}

	return score / 64
}

func ImageColorId(im image.Image) uint64 {

	var (
		id = uint64(0)
		ac = im.Bounds().Max.X * im.Bounds().Max.Y
		ar = []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
	)

	if ac > 0 {

		for x := 0; x < im.Bounds().Max.X; x++ {
			for y := 0; y < im.Bounds().Max.Y; y++ {
				r, g, b, _ := im.At(x, y).RGBA()

				n := int(r / 16384)
				ar[n] += 1

				n = int(g / 16384)
				ar[4+n] += 1

				n = int(b / 16384)
				ar[8+n] += 1
			}
		}

		for i, v := range ar {
			n := ((v * 100) / ac) / 20
			id = id | 1<<uint64((i*5)+n)
		}
	}

	return id
}

func Float32Round(f float32, pOffset int64) float32 {
	if n := pOffset % 2; n > 0 {
		pOffset += n
	}
	if pOffset < 2 {
		pOffset = 2
	} else if pOffset > 8 {
		pOffset = 8
	}
	pa_fix := float32(1e4)
	switch pOffset {
	case 2:
		pa_fix = 1e2
	case 4:
		pa_fix = 1e4
	case 6:
		pa_fix = 1e6
	case 8:
		pa_fix = 1e8
	default:
		pa_fix = 1e4
	}
	return float32(int64(f*pa_fix+0.5)) / pa_fix
}
//func ImageHashPrint(a uint64) {
//	msg := ""
//	for i := uint64(0); i < 8; i++ {
//		for j := uint64(0); j < 8; j++ {
//			if u64Allow(a, 1<<((i*8)+j)) {
//				msg += "1 "
//			} else {
//				msg += "0 "
//			}
//		}
//		msg += "\n"
//	}
//	fmt.Println(msg)
//}

func ImgShmBytes(file string) ([]byte, error) {
	fp, err := os.Open(file)
	if err == nil {
		defer fp.Close()
		return ioutil.ReadAll(fp)
	}
	return []byte{}, err
}

func ImageTag(m image.Image, tagSpace int) []uint8 {

	if tagSpace < 8 {
		tagSpace = 8
	} else if tagSpace > 128 {
		tagSpace = 128
	}

	var (
		bounds = m.Bounds()
		dx     = bounds.Dx()
		dxs    = dx / tagSpace
		dy     = bounds.Dy()
		dys    = dy / tagSpace
		tags   []uint8
	)

	for i := dxs; i < dx; i += dxs {
		for j := dys; j < dy; j += dys {
			p := m.At(i, j)
			r, g, b, _ := p.RGBA()
			gc := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
			tags = append(tags, uint8(gc/16))
		}
	}

	return tags
}

func ImagEmpty(a []uint8) float32 {
	if len(a) < 1 {
		return 0
	}
	hit := float32(0)
	for i := 0; i < (len(a) - 1); i++ {
		if a[i] == a[i+1] {
			hit += 1
		}
	}
	return hit / float32(len(a))
}

func ImageSimilarity(a, b []uint8) float32 {
	hit := float32(0)
	if len(a) < 1 || len(a) != len(b) {
		return 0
	}
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			hit += 1
		}
	}
	return hit / float32(len(a))
}

func init() {
	for i := uint64(0); i < 64; i++ {
		u64Arr[i] = 1 << i
	}
}