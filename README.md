## ospp-go-grpc

Go Plugin over gRPC is a lightweight Golang two-way communication plugin framework with built-in health check, plugin
management, logging and other extensions.

The framework is lightweight and very easy to use, you only need to spend a little time to migrate the existing plugin
implementation, or build your own plugin system from scratch.

For more details about the framework, please refer to the [examples](examples) directory.

## Core Concepts

Understanding the core concepts helps to grasp the whole. The core concepts in the solution are few and easy to understand, and the design solution strives to ensure the goal of "lightweight" and "easy to use".

### Core

`Core` is the host, it is the core implementation of a project, mainly containing the project logic, but basically does not contain specific business logic. Developers and users call functions in the plugin through `Core`. `Core` is designed to be a `Server`.

### Plugin

`Plugin` is a plugin that implements the specific business logic required by `Core`. The `Plugin` is designed to be a `Client`.

### Convention

The design solution favors "conventions" over "constraints". If you want to be as strongly constrained as `go-plugin`, the developer must spend a lot of time write the constraint. For a small, lightweight project, constraints should only be optional. The plugin name, version, function name, argument type, and result type between `Core` and `Plugin` should be discussed first and constrained later whenever possible.

### Interface

`Core` can define multiple `interface`. Similar to but different from `interface` in `Golang`. To enforce a single responsibility, a `Plugin` is allowed to implement only one interface of `Core` and perform one responsibility.

### Func/Handler

The business logic functions registered in a `Plugin` are called `Func/Handler`. The function signature is `func(ctx plugin.Context) (interface{}, error)`. It receives arguments from `Core`, processes them and returns the result to `Core`.
