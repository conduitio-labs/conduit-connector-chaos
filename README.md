# Conduit Connector Chaos

A [Conduit](https://conduit.io) connector that can be configured to behave in
unexpected ways to figure out how Conduit handles it.

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests.

## Modes

Each method of the connector can be configured to act in a certain way.
Currently supported modes are:

- `success` - The connector method acts as expected (returns successfully).
- `error` - The connector method returns an error with the message "chaos".
- `block` - The connector method blocks forever and never returns.
- `context-done` - The connector method blocks until the supplied context is
  closed, after that it returns the context error.
- `panic` - The connector method panics.

## Source

Each source method can be configured to act as described in one of the
[modes](#Modes).

For a list of methods and their descriptions
see [`sdk.Source` interface](https://pkg.go.dev/github.com/conduitio/conduit-connector-sdk#Source)
. Note that the behavior of `Parameters` can't be adjusted, because that method
is called before the connector receives the configuration.

### Configuration

| name            | description                                           | required | default value |
|-----------------|-------------------------------------------------------|----------|---------------|
| `configureMode` | Controls the mode that the Configure method executes. | false    | success       |
| `openMode`      | Controls the mode that the Open method executes.      | false    | success       |
| `readMode`      | Controls the mode that the Read method executes.      | false    | success       |
| `ackMode`       | Controls the mode that the Ack method executes.       | false    | success       |
| `teardownMode`  | Controls the mode that the Teardown method executes.  | false    | success       |

## Destination

Each destination method can be configured to act as described in one of the
[modes](#Modes).

For a list of methods and their descriptions
see [`sdk.Destination` interface](https://pkg.go.dev/github.com/conduitio/conduit-connector-sdk#Destination)
. Note that the behavior of `Parameters` can't be adjusted, because that method
is called before the connector receives the configuration.

### Configuration

| name            | description                                           | required | default value |
|-----------------|-------------------------------------------------------|----------|---------------|
| `configureMode` | Controls the mode that the Configure method executes. | false    | success       |
| `openMode`      | Controls the mode that the Open method executes.      | false    | success       |
| `writeMode`     | Controls the mode that the Write method executes.     | false    | success       |
| `teardownMode`  | Controls the mode that the Teardown method executes.  | false    | success       |
