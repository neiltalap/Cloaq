# Synapse Encapsulation Protocol (v1)

Each packet sent over the UDP transport follows this 4-byte header format followed by the encrypted payload.

## Header Structure

| Offset | Field   | Size    | Description                          |
|--------|---------|---------|--------------------------------------|
| 0      | Version | 1 byte  | Protocol version (currently 0x01)    |
| 1      | Type    | 1 byte  | Message type (0: Handshake, 1: Data) |
| 2-3    | Length  | 2 bytes | Payload size (Big-Endian, max 64KB)  |
| 4-N    | Payload | N bytes | Encrypted data (Symmetric Key)       |

## Constraints
- **Max Payload Size:** 65,535 bytes (limited by 16-bit length field).
- **Endianness:** Big-Endian (Network Byte Order) for the length field.