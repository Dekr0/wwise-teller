package decoder

import "fmt"

func WrongBKHDPosition(p string) error {
	return fmt.Errorf("Sound bank %s does not have BKHD data chunk to be the first data chunk.", p)
}

func NoBKHD(p string) error {
	return fmt.Errorf("Sound bank %s is missing BKHD section", p)
}

func NoDIDX(p string) error {
	return fmt.Errorf("Sound bank %s is missing DIDX section", p)
}

func NoDATA(p string) error {
	return fmt.Errorf("Sound bank %s is missing DATA section", p)
}

func NoHIRC(p string) error {
	return fmt.Errorf("Sound bank %s is missing HIRC section", p)
}

func AKBKBank(p string) error {
	return fmt.Errorf("Detected AKBK chunk. Unwise does not support legacy Wwise sound bank %s.", p)
}

func LegacyBank(p string, version u32) error {
	return fmt.Errorf("Unwise does not support legacy Wwise sound bank %s with version %d", p, version)
}

func CustomBank(p string, version u32) error {
	return fmt.Errorf("Unwise yet support Wwise custom sound bank %s with version %d", p, version)
}

func UnknownCustomBank(p string, version u32) error {
	return fmt.Errorf("Unwise yet support Wwise unknown custom sound bank %s with version %d", p, version)
}

func EncryptionBank(p string) error {
	return fmt.Errorf("Unwise yet support decryption of encrypted Wwise sound bank %s", p)
}

func UnsupportKnownBank(p string, version u32) error {
	return fmt.Errorf("Unwwise yet support Wwise sound bank with version %d", version)
}
