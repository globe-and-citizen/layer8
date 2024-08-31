package utils

import (
	stdBytes "bytes"
	"io"
	"log"
)

func Equal(fst []byte, snd []byte) bool {
	if len(fst) != len(snd) {
		return false
	}

	for i := 0; i < len(fst); i++ {
		if fst[i] != snd[i] {
			return false
		}
	}

	return true
}

func WriteBytes[T io.WriterTo](key T) []byte {
	var byteWriter stdBytes.Buffer

	_, err := key.WriteTo(&byteWriter)
	if err != nil {
		log.Fatalf("Error while writing key to byte buffer: %e", err)
	}

	return byteWriter.Bytes()
}

func ReadBytes[T io.ReaderFrom](key T, keyBytes []byte) {
	byteBuffer := stdBytes.NewBuffer([]byte{})

	_, err := byteBuffer.Write(keyBytes)
	if err != nil {
		log.Fatalf("Error while writing bytes to byte buffer: %e", err)
	}

	_, err = key.ReadFrom(byteBuffer)
	if err != nil {
		log.Fatalf("Error while decoding key from bytes: %e", err)
	}
}
