package decode

import (
	"bufio"
	"bytes"
	"image"
	"io"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"github.com/schollz/progressbar/v3"
)

func ExtractJPEGsFromMJPEG(r io.Reader, imgCount int) ([]image.Image, error) {
	log.PrefixKey = "image decode"
	bar := progressbar.Default(int64(imgCount), "Decoding images")
	var images []image.Image
	bufReader := bufio.NewReader(r)

	for {
		// Ищем начало JPEG
		start, err := bufReader.Peek(2)
		if err != nil {
			log.Debugf("Failed to peek: %s", err)
			break
		}
		if start[0] != 0xFF || start[1] != 0xD8 {
			_, _ = bufReader.ReadByte()
			continue
		}

		// Считываем JPEG в буфер до EOI (0xFFD9)
		var jpegBuf bytes.Buffer
		for {
			b, err := bufReader.ReadByte()
			if err != nil {
				return nil, err
			}
			jpegBuf.WriteByte(b)
			if jpegBuf.Len() >= 2 {
				lastTwo := jpegBuf.Bytes()[jpegBuf.Len()-2:]
				if lastTwo[0] == 0xFF && lastTwo[1] == 0xD9 {
					break
				}
			}
		}

		img, err := imaging.Decode(&jpegBuf)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
		bar.Add(1)
	}

	return images, nil
}
