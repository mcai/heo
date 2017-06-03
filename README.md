# README for Heo

# Heo README

This README briefly describes what Heo is and how to setup & run Heo.

## License Information

Copyright (c) 2010-2017 by Min Cai (<min.cai.china@bjut.edu.cn>)

This program is free software, licensed under the MIT license.

## About

Heo is a cycle-accurate CPU-GPU heterogeneous multicore architectural simulator.

Heo is written in Go.

## Key Features

- Heo facilitates interface-based hierarchical object orientation to improve code readability and extensibility.

- Heo provides **functional architectural simulation** of:
	- Program loading for statically-linked MIPS32 ELF executables (both Little Endian and Big Endian are supported);
	- Functional execution of common integer and floating-point instructions of MIPS32 ISA;
	- Emulation of common POSIX system calls;
	- Execution of Pthreads based multithreaded programs.

- Heo provides **cycle-accurate microarchitectural simulation** of:
	- Separate pipeline structures such as the reorder buffer, separate integer and floating point physical register files;
	- Explicit register renaming based on the merged architectural and rename register file model;
	- Single-threaded superscalar out-of-order execution, multithreaded SMT and CMP execution model;
	- Multi-level inclusive cache hierarchy with the directory-based MESI coherence protocol;
	- Simple cycle-accurate DRAM controller model;
	- Various kinds of static and dynamic branch predictors, checkpointing-based pipeline recovery on branch misprediction (**with bugs**).

- Heo supports the following **unclassified simulation features**:
	- Support measurement of instructions, pipeline structures and the memory hierarchy;
	- Support generation of instruction traces;
	- Support both execution-driven and trace-driven NoC subsystem simulation;
	- Support seamless switching between functional simulation and performance simulation mode;
	- Support the thread based data prefetching scheme: classification of good, bad and ugly prefetch requests.

- Heo provides the following **common infrastructure**:
	- Scheduling and dispatching framework for modeling synchronous (cycle accurate) and asynchronous activities;
	- Easy configuration and statistics reporting of the simulated architectures, workloads and simulations.

- Heo currently supports correct execution of all the **Olden benchmark suite** except incorrect output from "health", plus some of CPU2006 benchmarks. Other benchmarks are being tested.

## TODOs

1. Support CPU2006, PARSEC and Rodinia.

2. Support CPU-GPU heterogeneous multicore architectures.

3. Support CPU frequency control (multiple clock domain).

4. Support DRAM-NVM hybrid memory systems.

5. Support Command line options and Component Properites.

6. Support multiple hardware prefetchers.

7. Support application of machine learning algorithms in architectural components.

8. Support deep_anpr.

## Dependencies

1. `sudo apt-get install python3-pip python3-tk`

2. `pip3 install --upgrade pip`

3. `pip3 install matplotlib pandas seaborn objectpath`


## System Requirements

Heo has been tested on 64-bit MacOS Sierra (10.12.4) and **Ubuntu Linux 16.04** (with x86 machines).

For **developing and running Heo**, make sure that:

1. You have Ubuntu Linux 16.04 or higher (or other mainstream Linux distributions);

2. To compile and run Go programs, the following software must be installed:
	- golang 1.8.1

## Quick Start

git clone https://github.com/mcai/heo

## Customizing Heo for Your Needs

- As is always true for open source software, the existing Heo code are good examples for demonstrating its usage and power.

## Contact

Please report bugs and send suggestions to:

Min Cai (<min.cai.china@bjut.edu.cn>)

Faculty of Information Technology, Beijing University of Technology, Beijing, China
