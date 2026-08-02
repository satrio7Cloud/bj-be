package qris

import "strings"

const (
	crcPolynomial = 0x1021
	crcInitial    = 0xffff
)

func CRC16CCITTFalse(input string) string {
	crc := crcInitial
	for _, char := range input {
		crc ^= int(char) << 8
		for bit := 0; bit < 8; bit++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ crcPolynomial
			} else {
				crc <<= 1
			}
			crc &= 0xffff
		}
	}
	return strings.ToUpper(hex4(crc))
}

func StripCRC(payload string) string {
	index := strings.LastIndex(payload, "6304")
	if index == -1 {
		return payload
	}
	return payload[:index]
}

func AppendCRC(payloadWithoutCRC string) string {
	base := StripCRC(payloadWithoutCRC)
	crc := CRC16CCITTFalse(base + "6304")
	return base + "6304" + crc
}

func IsCRCValid(payload string) bool {
	current := PayloadCRC(payload)
	if current == "" {
		return false
	}
	return current == CRC16CCITTFalse(StripCRC(payload)+"6304")
}

func PayloadCRC(payload string) string {
	index := strings.LastIndex(payload, "6304")
	if index == -1 || index+8 != len(payload) {
		return ""
	}
	return strings.ToUpper(payload[index+4:])
}

func hex4(value int) string {
	const chars = "0123456789ABCDEF"
	out := []byte{'0', '0', '0', '0'}
	for i := 3; i >= 0; i-- {
		out[i] = chars[value&0xF]
		value >>= 4
	}
	return string(out)
}
