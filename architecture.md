# UDAL (Universal Decentralized Anonymity Layer)

A protocol-agnostic, high-throughput privacy transport separating identity from network location. It is designed to be native to cryptographic financial systems and hardened against both global passive adversaries and active routing attacks.

## Core Architecture

               Applications
           (Browser, Game, DeFi/API)
                 ▲
                 |
          [ UDAL Local Node ]
                 |
      (Identity | Routing | Encryption)
                 |
                 ▼
          [ UDAL Mesh Network ]

## Key Architectural Principles

1. **Identity != IP**: Routing uses cryptographic keys, not IP addresses.
2. **Stateless Forwarding (HORNET)**: Eliminates per-flow state on intermediate nodes. Onion layer decryption keys and next-hop info are embedded into Forwarding Segments (FS) carried in the Anonymous Header (AHDR) of every packet, executing symmetric crypto at line rate (93+ Gb/s).
3. **Path-Aware Infrastructure (SCION)**: Replaces legacy BGP routing. Utilizes Isolation Domains (ISDs) to prevent hijacking and allows end-hosts to select specific path segments for geofencing, multi-path communication, and fast failover.
4. **Line-Rate Traffic Morphing (ditto principles)**: Hardware-accelerated or software-optimized traffic shaping using fixed, periodic patterns and "chaff" (dummy packets) to completely defeat AI/Transformer-based traffic classification (e.g., concealing metadata and traffic shape).
5. **Isomorphic Networks & Task Indistinguishability**: All nodes run the identical single binary, dynamically multiplexing roles (client, relay, auditor).
6. **Native Decentralized Incentives**: Deep integration with high-throughput L1s (e.g., Sui via Mysticeti V2, or Solana) for scalable "Proof-of-Relay" rewards to sustain network density and provide Sybil resistance.

## Technology Stack

- **Core Transport**: SCION topology integrated with fast UDP transport.
- **Onion Routing**: HORNET Sphinx-format packet headers.
- **Encryption**: AES-256-GCM (Hardware Accelerated via AES-NI/ARM Crypto).
- **Obfuscation**: Constant-size padding and mixnet delays ("ditto" principles).