# go-agent

`Agent` aims to simplify concurrent resource management by developing an `efficient`, `user-friendly` task scheduler `with shared-resource support` in Go.

## Current Status

MVP was developed with the help of AI-assisted development. We are currently developing an AI-assisted project generator ([Solus](github.com/CSXL/solus)), so we are using this project as a testing ground for manual experimentation. In other words, we are using this project to actuate the language models ourselves so we can create the autonomous agents better.

## Requirements

### 1. Efficient Task Scheduler

- 1.1. Optimize task scheduling to minimize resource usage and wait times
- 1.2. Implement prioritization of tasks based on user-defined parameters
- 1.3. Support dynamic adjustment of task priorities during runtime

### 2. User-friendly Interface

- 2.1. Develop a clear and concise API for interacting with the task scheduler
- 2.2. Provide detailed documentation, including examples and best practices
- 2.3. Ensure smooth integration with existing Go projects

### 3. Shared-resource Support

- 3.1. Implement a mechanism for managing shared resources among concurrent tasks
- 3.2. Ensure safe and efficient resource allocation and deallocation
- 3.3. Support user-defined resource constraints, such as maximum resource usage per task
- 3.4. Provide a mechanism for extending the resource management system to support new resource types

### 4. Robustness and Reliability

- 4.1. Implement error handling and recovery mechanisms for task failures
- 4.2. Provide monitoring and logging capabilities for better understanding of task execution
- 4.3. Ensure compatibility with various Go versions and platforms

### 5. Testing and Validation

- 5.1. Develop comprehensive test cases covering allmajor functionalities and edge cases
- 5.2. Perform performance benchmarking against existing task schedulers
- 5.3. Validate correct resource allocation and deallocation under various scenarios

### 6. Maintainability and Extensibility

- 6.1. Implement a modular and clean code structure for easy maintenance
- 6.2. Ensure adherence to Go coding standards and best practices
- 6.3. Design the system to support future enhancements and feature additions

### 7. Security

- 7.1. Ensure safe execution of tasks by implementing proper isolation and sandboxing mechanisms
- 7.2. Implement access control mechanisms for shared resources
- 7.3. Protect against common security vulnerabilities, such as race conditions and deadlocks

### 8. Deployment and Distribution

- 8.1. Package the library for easy installation and integration with existing projects
- 8.2. Provide clear instructions for deployment, configuration, and usage
- 8.3. Publish the library to relevant package repositories for easy discovery and adoption
