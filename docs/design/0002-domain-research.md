# Domain Research: Configuration Objects and Parameters for Simulation

- Feature Name: `domain-research`
- Status: Draft

# Summary

This document presents domain research on the network objects and configuration parameters required to model Layer 3+ (IP and transport layer) forwarding and flow-based traffic simulation.

This document establishes:

1. **Which object types** must exist in configuration files.
2. **Which parameters** each object requires, backed by an analysis of production routers (Cisco IOS, Juniper JunOS) and modern data center architectures.
3. **Evaluation criteria for RFC 0001 (`0001-configuration-language.md`)**: whether the proposed configuration language can express all required domain parameters, and what adjustments RFC 0001 needs before acceptance.

# Motivation

Before finalizing and accepting the configuration language design in RFC 0001, we must establish the exact domain model the language is meant to configure.

For `usernet`, our scope is **L3+ forwarding and flow-level traffic simulation**.
We need a clear inventory of domain objects and their parameters so we can:

- Avoid bloating the simulator with legacy Layer 1/Layer 2 mechanics or hardware-specific plumbing.
- Verify that the value types, literals, and structural constructs in RFC 0001 are necessary and sufficient to express real-world L3+ topologies.
- Provide a direct specification for generating the future `schema.json` contract.

# 1. Domain Research

## 1.1. Pruning Production Interface Types

To determine which interface abstractions are genuinely needed for an L3+ simulator, we audited the full list of interface types available in the two dominant production network operating systems: **Cisco IOS** and **Juniper JunOS**.

### A. Cisco IOS Audit

- **Legacy WAN & Modem Interfaces (Excluded):** `Async`, `Group-Async`, `Dialer`, `CDMA-Ix`, `MFR` (Frame Relay), `LongReachEthernet`, `Virtual-TokenRing`. These represent obsolete 1980s/1990s physical/data-link technologies irrelevant to modern IP simulation.
- **Subscriber Templates & CLI Helpers (Excluded):** `Auto-Template`, `Virtual-Template`, `Vif`, `vmi`, `range` (which is a CLI macro, not an interface).
- **Layer 2 Bridging & MPLS L2VPNs (Excluded):** `BVI` (Bridge-Group Virtual Interface) and `Pseudowire`. Because `usernet` operates strictly at L3+, L2 switching domains and VLAN bridging are out of scope.
- **MPLS / Advanced Overlays (Deferred):** `GMPLS`, `Lspvif`, `LISP`. Can be added in future iterations if MPLS Traffic Engineering will be required.
- **Core L3 Interfaces (Retained):** `GigabitEthernet` (`ethernet`), `Loopback` (`loopback`), `Tunnel` / `CTunnel` (`tunnel`), `Null` (`blackhole`), and `Virtual-PPP` / `Multilink` (merged into `ethernet`).

### B. Juniper JunOS Audit

Auditing Juniper JunOS (`MX` / `PTX` / `QFX` / `SRX` series) reveals a different set of hardware and software abstractions that we also prune:

- **Speed-Prefixed Physical Ethernet (`fe-`, `ge-`, `xe-`, `et-`, `fxp0`, `em0`) -> Merged into `ethernet`:**
  JunOS uses distinct interface prefixes for 100M (`fe-`), 1G (`ge-`), 10G (`xe-`), 40G/100G/400G (`et-`), and management (`fxp0`/`em0`). In `usernet`, all of these collapse into a single `ethernet` type with a `bandwidth` parameter.
- **Legacy SONET/ATM/TDM WAN (`so-`, `at-`, `t1-`, `t3-`, `e1-`, `e3-`, `sp-`, `ls-`, `pp0`) -> Excluded:**
  Legacy optical/TDM leased lines (`SONET/SDH`, `ATM`, `T1/E1`) and PPPoE subscriber termination (`pp0`).
- **Internal ASIC / Packet Forwarding Engine (PFE) Plumbing (`pime`, `pimd`, `pd-`, `pe-`, `pfh-`, `mt-`, `cbp0`, `demux0`, `jsrv`, `tap`, `lt-`) -> Excluded:**
  JunOS exposes internal hardware recirculation pseudo-interfaces (used by router ASICs to encapsulate/decapsulate multicast PIM packets, subscriber demuxing, or inter-instance stitching). A software simulator performs these operations directly in memory without fake hardware interfaces.
- **L2 Bridging & Link Aggregation (`irb`, `vlan`, `ae`) -> Excluded / Deferred:**
  `irb` (Integrated Routing and Bridging, analogous to Cisco `BVI`) is excluded since L2 switching is omitted. `ae` (Aggregated Ethernet / LACP LAG) is deferred to future iterations.
- **Protocol-Specific Tunnels (`gr-`, `ip-`, `st0`, `vtep`) and Flexible Tunnel Interface (`fti`) -> Merged into `tunnel`:**
  Historically, JunOS required different interface types for GRE (`gr-`), IP-in-IP (`ip-`), IPsec (`st0`), and VXLAN (`vtep`). Modern JunOS replaced this fragmented model with **`fti` (Flexible Tunnel Interface)**, where a single tunnel interface type accepts an encapsulation mode parameter.
- **Loopback (`lo0`) and Discard (`dsc`) -> Retained as `loopback` and `blackhole`:**
  `lo0` maps directly to `loopback`, and JunOS's `dsc` (Discard) interface (as well as `discard` static route next-hops) maps directly to `blackhole`.

### Synthesis

| Concept                   | Cisco IOS                                                              | Juniper JunOS                                  | `usernet` Abstraction |
| :------------------------ | :--------------------------------------------------------------------- | :--------------------------------------------- | :-------------------- |
| **Wired L3 Port**         | `FastEthernet`, `GigabitEthernet`, `TenGigabitEthernet`, `Virtual-PPP` | `fe-`, `ge-`, `xe-`, `et-`                     | **`ethernet`**        |
| **Node Internal Address** | `Loopback`                                                             | `lo0`                                          | **`loopback`**        |
| **L3 Overlay Tunnel**     | `Tunnel`, `CTunnel`                                                    | `fti` (and legacy `gr-`, `ip-`, `st0`, `vtep`) | **`tunnel`**          |
| **Traffic Sink (Drop)**   | `Null0`                                                                | `dsc` (Discard)                                | **`blackhole`**       |

## 1.2. Interfaces Types `ethernet` and `p2p`

In real networks, `Ethernet` and `Point-to-Point (PPP/HDLC)` interfaces differ only because of Layer 2 mechanics:

1. **L2 Headers & ARP:** Ethernet frames require 6-byte MAC addresses, forcing routers to run ARP/NDP before sending packets. PPP/HDLC frames have no MAC addresses and never run ARP.
2. **Broadcast Multi-Access Pitfalls:** Real routers assume an Ethernet port might be plugged into an L2 switch with hundreds of hosts. Pointing a static route solely to an Ethernet port (`ip route 0.0.0.0 0.0.0.0 Gi0/0` in IOS, or `qualified-next-hop ge-0/0/0.0` in JunOS on a broadcast link) without a next-hop IP triggers catastrophic **Proxy ARP** storms or validation errors, and OSPF waits 40 seconds to elect a Designated Router (DR/BDR) unless manually overridden with `point-to-point` mode.

Because `usernet` models L3+ forwarding, **every physical link in topology graph connects exactly two interfaces (1-to-1)**. Therefore:

- MAC addresses and ARP are completely unnecessary (the simulator knows the peer interface on the other end of any link in $O(1)$ time).
- Separating `ethernet` and `p2p` is redundant. A single **`ethernet`** interface type combines physical Ethernet wire characteristics (`mtu = 1500B`, `l2_overhead = 18B`) with pure Point-to-Point L3 routing semantics (routes can specify `interface`, `next_hop`, or both; IP Unnumbered links work out of the box).

## 1.3. Traffic Simulation

For L3+ forwarding, routing lookups and ECMP hashing make identical decisions for all packets sharing the same header (`5-tuple`). Simulating individual packets at 10–100 Gbps creates massive CPU overhead without adding L3 routing insight.

Therefore, `usernet` should model traffic as **Flows** (fluid or discrete-event flow chunks characterized by a `5-tuple` header, data rate/volume, and average packet size).

---

# 2. Required Configuration Objects and Parameters

Based on our research, the configuration system needs **5 core object types**.

## Object 1: `Node`

Represents a router, L3 switch, firewall, or end-host.

| Parameter    | Meaning / Purpose                                                                                  | Required Data Type | Default                             |
| :----------- | :------------------------------------------------------------------------------------------------- | :----------------- | :---------------------------------- |
| `name`       | Unique identifier of the hop in the topology.                                                      | `String`           | _(Required)_                        |
| `forwarding` | Whether the node routes transit traffic destined for other nodes. (`false` for end-hosts).         | `Bool`             | `true` for router, `false` for host |
| `ttl`        | Initial IP TTL (hop limit) for traffic originated by this node (`1..255`).                         | `Int`              | `64`                                |
| `ecmp_hash`  | Fields used to balance flows across equal-cost routes (`five_tuple`, `three_tuple`, `src_dst_ip`). | `Enum`             | `five_tuple`                        |

## Object 2: `Interface`

Represents a physical port or logical L3 interface on a node.

| Parameter     | Meaning / Purpose                                                            | Required Data Type                 | Default                                           |
| :------------ | :--------------------------------------------------------------------------- | :--------------------------------- | :------------------------------------------------ |
| `name`        | Interface identifier unique within the node (e.g., `"eth0"`, `"lo0"`).       | `String`                           | _(Required)_                                      |
| `type`        | Operational model: `ethernet`, `loopback`, `tunnel`, `blackhole`.            | `Enum`                             | `ethernet`                                        |
| `state`       | Administrative state (`up` or `down`).                                       | `Enum`                             | `up`                                              |
| `ipv4`        | One or more IPv4 addresses with prefix length.                               | `IPv4Prefix` or `List<IPv4Prefix>` | `[]`                                              |
| `ipv6`        | One or more IPv6 addresses with prefix length.                               | `IPv6Prefix` or `List<IPv6Prefix>` | `[]`                                              |
| `mtu`         | Maximum L3 packet size in bytes before fragmentation or drop.                | `Size`                             | `1500B` (`eth`), `65536B` (`lo`), `1476B` (`tun`) |
| `metric`      | Base routing cost of the interface (must be `>= 1`).                         | `Int`                              | `10`                                              |
| `bandwidth`   | Egress link serialization rate (e.g., `1Gbps`, `2.5Gbps`).                   | `Rate`                             | Unlimited                                         |
| `delay`       | Egress propagation latency (e.g., `2.5ms`, `500us`).                         | `Duration`                         | `0ms`                                             |
| `jitter`      | Latency variance on egress traffic.                                          | `Duration`                         | `0ms`                                             |
| `loss`        | Packet/flow drop probability (`0%` to `100%` or `0.0` to `1.0`).             | `Percentage / Float`               | `0.0`                                             |
| `queue_size`  | Egress buffer size in bytes for congestion modeling.                         | `Size`                             | `64KB`                                            |
| `l2_overhead` | Wire framing overhead per packet (14B Ethernet header + 4B FCS).             | `Size`                             | `18B` (`eth`)                                     |
| `tunnel_mode` | Encapsulation protocol if `type = tunnel` (`gre`, `ipip`, `vxlan`, `ipsec`). | `Enum`                             | `gre`                                             |
| `tunnel_src`  | Local underlay source IP (required if `type = tunnel`).                      | `IPv4Addr \| IPv6Addr`             | _(Conditional)_                                   |
| `tunnel_dst`  | Remote underlay destination IP (required if `type = tunnel`).                | `IPv4Addr \| IPv6Addr`             | _(Conditional)_                                   |

## Object 3: `Route`

Represents a static routing table entry inside a node.

| Parameter        | Meaning / Purpose                                                         | Required Data Type         | Default      |
| :--------------- | :------------------------------------------------------------------------ | :------------------------- | :----------- |
| `dst`            | Destination subnet prefix (e.g., `0.0.0.0/0` or `10.0.0.0/8`).            | `IPv4Prefix \| IPv6Prefix` | _(Required)_ |
| `next_hop`       | Neighbor IP address (optional if `interface` is specified).               | `IPv4Addr \| IPv6Addr`     | _(Optional)_ |
| `interface`      | Egress interface name on this node (optional if `next_hop` is specified). | `String`                   | _(Optional)_ |
| `metric`         | Preference among routes from the same source (equal metrics -> ECMP).     | `Int`                      | `1`          |
| `admin_distance` | Priority across routing sources (`1..255`, lower wins; `0` = connected).  | `Int`                      | `1`          |

## Object 4: `Link`

Connects two `ethernet` interfaces across two different nodes (`1-to-1` duplex wire).

| Parameter    | Meaning / Purpose                                                            | Required Data Type   | Default      |
| :----------- | :--------------------------------------------------------------------------- | :------------------- | :----------- |
| `src`, `dst` | The pair of connected node interfaces.                                       | `Endpoint` pair      | _(Required)_ |
| `bandwidth`  | Optional shorthand to set symmetric `bandwidth` on both attached interfaces. | `Rate (bps)`         | _(Optional)_ |
| `delay`      | Optional shorthand to set symmetric propagation `delay` on both interfaces.  | `Duration`           | _(Optional)_ |
| `loss`       | Optional shorthand to set symmetric packet `loss` on both interfaces.        | `Percentage / Float` | _(Optional)_ |
| `state`      | Physical link state (`up` or `down` for simulating link failures).           | `Enum`               | `up`         |

## Object 5: `Flow`

Represents a L3/L4 traffic stream injected into the network for simulation and verification.

| Parameter       | Meaning / Purpose                                                                         | Required Data Type     | Default                 |
| :-------------- | :---------------------------------------------------------------------------------------- | :--------------------- | :---------------------- |
| `name`          | Identifier for tracing and reporting statistics on this flow.                             | `String`               | _(Required)_            |
| `src_node`      | Originating node name (or resolved automatically from `src_ip`).                          | `String`               | _(Required)_            |
| `src_ip`        | Source IPv4 or IPv6 address.                                                              | `IPv4Addr \| IPv6Addr` | _(Required)_            |
| `dst_ip`        | Destination IPv4 or IPv6 address.                                                         | `IPv4Addr \| IPv6Addr` | _(Required)_            |
| `protocol`      | L4 protocol (`tcp`, `udp`, `icmp`, `gre`, `esp`).                                         | `Enum`                 | `tcp`                   |
| `src_port`      | L4 source port (`0..65535`, if `tcp` or `udp`).                                           | `Int`                  | `49152`                 |
| `dst_port`      | L4 destination port (`0..65535`, if `tcp` or `udp`).                                      | `Int`                  | `80`                    |
| `rate`          | Offered traffic load (e.g., `500Mbps`).                                                   | `Rate (bps)`           | `1Mbps`                 |
| `packet_size`   | Average L3 packet size (needed for MTU checks, tunnel overhead, and serialization delay). | `Size (Bytes)`         | `1500B`                 |
| `dont_fragment` | IP `DF` bit: if `true`, flow drops when `packet_size > interface.mtu`.                    | `Bool`                 | `true`                  |
| `behavior`      | Congestion behavior: `elastic` (TCP-like fair sharing) vs `inelastic` (UDP-like drop).    | `Enum`                 | `elastic`               |
| `start_time`    | When the flow starts during timed simulations.                                            | `Duration`             | `0s`                    |
| `duration`      | How long the flow remains active (or `total_volume` in `Size`).                           | `Duration`             | Infinite / Steady-state |

---

## Drawbacks

[drawbacks]: #drawbacks

Why should we _not_ do this?

## Rationale and alternatives

[rationale-and-alternatives]: #rationale-and-alternatives

- Why is this design the best in the space of possible designs?
- What other designs have been considered and what is the rationale for not choosing them?
- What is the impact of not doing this?

## Prior art

[prior-art]: #prior-art

Discuss prior art, both the good and the bad, in relation to this proposal.

## Unresolved questions

[unresolved-questions]: #unresolved-questions

- What parts of the design do you expect to resolve through the review process before this gets merged?
- What parts of the design do you expect to resolve through the implementation of this feature?
- What related issues do you consider out of scope for this document that could be addressed in the future independently of the solution that comes out of this design?

## Future possibilities

[future-possibilities]: #future-possibilities

Think about what the natural extension and evolution of your proposal would be and how it would affect the project as a whole in a holistic way.

This is also a good place to "dump ideas", if they are out of scope for the document you are writing but otherwise related.

If you have tried and cannot think of any future possibilities, you may simply state that you cannot think of anything.
