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

- Avoid bloating the simulator with legacy Layer 1/Layer 2 mechanics.
- Verify that the value types, literals, and structural constructs in RFC 0001 are necessary and sufficient to express real-world L3+ topologies.
- Provide a direct specification for generating the future `schema.json` contract.
