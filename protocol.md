# Hanover Flipdot Display Protocol Description

## Overview

The Hanover flipdot display uses a serial protocol for communication, typically over RS485. This document describes the packet structure, commands, and data format used to control the display.

## Packet Structure

Each packet sent to the Hanover display follows this structure:

```
[STX][Command][Address][Resolution][Pixel Data][ETX][Checksum]
```

### Field Descriptions

1. **STX (Start of Text)**
   - 1 byte
   - Fixed value: 0x02
   - Indicates the start of a packet

2. **Command**
   - 1 byte
   - Specifies the type of command being sent
   - Known commands:
     - 0x31 ('1'): Write image data
   - Other command values may exist for different operations (e.g., starting/stopping test sequences)

3. **Address**
   - 1 byte
   - ASCII character representing the display's address
   - Range: '1' to '9' (0x31 to 0x39)
   - Set by a potentiometer inside the display

4. **Resolution**
   - 2 bytes
   - ASCII representation of a hexadecimal value
   - Represents (width * height) / 8
   - Example: For a 96x16 display, resolution would be 0x00C0 (192), sent as "C0"
   - Example: For a 128x16 display, resolution would be 0x0100 (256), sent as "00"

5. **Pixel Data**
   - Variable length
   - Each byte represents 8 pixels (1 bit per pixel)
   - Data is sent column-wise, from top to bottom, then left to right
   - Within each byte, bits represent pixels from bottom to top
   - Each byte is sent as two ASCII characters representing its hexadecimal value
   - Example: 0xAA would be sent as "AA"

6. **ETX (End of Text)**
   - 1 byte
   - Fixed value: 0x03
   - Indicates the end of the packet data

7. **Checksum**
   - 2 bytes
   - ASCII representation of a hexadecimal value
   - Calculated as follows:
     1. Sum all bytes from STX to ETX (inclusive)
     2. Subtract the STX value (0x02) from the sum
     3. Take the 8 least significant bits of the result
     4. Perform XOR with 0xFF
     5. Add 1 to the result
   - The final checksum is sent as two ASCII characters

## Pixel Data Format

The pixel data is organized as follows:

- Data is sent column-wise, starting from the top-left corner
- Each column is represented by (height / 8) bytes (rounded up)
- Bits in each byte represent pixels from bottom to top
- A '1' bit typically represents an "on" (visible) pixel
- A '0' bit typically represents an "off" (hidden) pixel

Example for a 16x16 display:
```
Column 1: [Byte 1: Pixels 1-8 (bottom to top)][Byte 2: Pixels 9-16 (bottom to top)]
Column 2: [Byte 1: Pixels 1-8 (bottom to top)][Byte 2: Pixels 9-16 (bottom to top)]
...
Column 16: [Byte 1: Pixels 1-8 (bottom to top)][Byte 2: Pixels 9-16 (bottom to top)]
```

## Communication Parameters

- Baud Rate: 4800 bps (typical, may vary)
- Data Bits: 8
- Stop Bits: 1
- Parity: None
- Flow Control: None

## Example Packet

For a 96x16 display with address 1, sending a pattern where all even columns are on:

```
02 31 31 43 30 AA 55 AA 55 ... AA 55 03 XX YY
```

- STX: 0x02
- Command: 0x31 ('1', Write image)
- Address: 0x31 ('1')
- Resolution: "C0" (96 * 16 / 8 = 192 = 0xC0)
- Pixel Data: "AA55AA55..." (repeated 96 times)
- ETX: 0x03
- Checksum: "XX YY" (calculated value)

## Implementation Notes

1. Always validate the packet structure, including STX, ETX, and checksum.
2. Handle partial packets and data reassembly if necessary.
3. Implement timeout mechanisms for incomplete packets.
4. Consider implementing error handling and retransmission strategies.
5. Ensure proper handling of different display sizes and resolutions.
6. Be aware that different commands (beyond writing image data) may exist for operations like starting or stopping test sequences. These would use different values for the Command byte.

## SDK Considerations

When building an SDK for the Hanover flipdot display, consider including the following features:

1. Packet construction and parsing functions
2. Checksum calculation and verification
3. High-level drawing primitives (e.g., lines, rectangles, text)
4. Display buffer management
5. Serial port handling and configuration
6. Error handling and logging
7. Examples and documentation for common use cases

By following this protocol description, you can create a robust SDK that accurately communicates with Hanover flipdot displays while providing a user-friendly interface for developers.
