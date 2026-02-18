# UDAL (Universal Decentralized Anonymity Layer)

A protocol-agnostic privacy transport separating identity from network location.

## Core Architecture

               Applications
           (Browser, Game, API)
                 ▲
                 |
          [ UDAL Local Node ]
                 |
      (Identity | Routing | Encryption)
                 |
                 ▼
          [ UDAL Mesh Network ]

## Key Concepts

1.  **Identity != IP**: Routing uses cryptographic keys, not IP addresses.
2.  **Decentralized Mesh**: Every node can act as a client, relay, or gateway. Meaning there will be a single binary to run everywhere.
3.  **Onion Routing**: Traffic is multi-hop encrypted using the **HORNET** protocol (Client → Entry → Middle → Exit).
4.  **Transport Agnostic**: Carries TCP, UDP, P2P, and hidden services.
5.  **WireGuard-class Speed**: Uses Noise protocol + UDP for high performance.