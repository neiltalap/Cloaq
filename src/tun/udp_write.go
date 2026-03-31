package tun

func WritePacket(device Device, packet []byte) error {
	_, err := device.Write(packet)
	return err
}
