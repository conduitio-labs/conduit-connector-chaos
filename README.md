# Conduit Connector Chaos

<!-- readmegen:description -->
A [Conduit](https://conduit.io) connector that can be configured to behave in
unexpected ways to figure out how Conduit handles it.

## Modes

Each method of the connector can be configured to act in a certain way.
Currently supported modes are:

- `success` - The connector method acts as expected (returns successfully).
- `error` - The connector method returns an error with the message "chaos".
- `block` - The connector method blocks forever and never returns.
- `context-done` - The connector method blocks until the supplied context is
  closed, after that it returns the context error.
- `panic` - The connector method panics.

Each source and destination method can be configured to act as described in one
of the modes above. For a list of methods and their descriptions see the [
`sdk.Source`](https://pkg.go.dev/github.com/conduitio/conduit-connector-sdk#Source)
and [
`sdk.Destination`](https://pkg.go.dev/github.com/conduitio/conduit-connector-sdk#Destination)
interfaces.<!-- /readmegen:description -->

### Source Configuration Parameters

<!-- readmegen:source.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "chaos"
        settings:
          # AckMode controls what the Ack method should do.
          # Type: string
          # Required: no
          ackMode: "success"
          # ConfigureMode controls what the Configure method should do.
          # Type: string
          # Required: no
          configureMode: "success"
          # OpenMode controls what the Open method should do.
          # Type: string
          # Required: no
          openMode: "success"
          # ReadMode controls what the Read method should do.
          # Type: string
          # Required: no
          readMode: "success"
          # TeardownMode controls what the Teardown method should do.
          # Type: string
          # Required: no
          teardownMode: "success"
          # Maximum delay before an incomplete batch is read from the source.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets read from the source.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Specifies whether to use a schema context name. If set to false, no
          # schema context name will be used, and schemas will be saved with the
          # subject name specified in the connector (not safe because of name
          # conflicts).
          # Type: bool
          # Required: no
          sdk.schema.context.enabled: "true"
          # Schema context name to be used. Used as a prefix for all schema
          # subject names. If empty, defaults to the connector ID.
          # Type: string
          # Required: no
          sdk.schema.context.name: ""
          # Whether to extract and encode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # The subject of the key schema. If the record metadata contains the
          # field "opencdc.collection" it is prepended to the subject name and
          # separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.key.subject: "key"
          # Whether to extract and encode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
          # The subject of the payload schema. If the record metadata contains
          # the field "opencdc.collection" it is prepended to the subject name
          # and separated with a dot.
          # Type: string
          # Required: no
          sdk.schema.extract.payload.subject: "payload"
          # The type of the payload schema.
          # Type: string
          # Required: no
          sdk.schema.extract.type: "avro"
```
<!-- /readmegen:source.parameters.yaml -->

## Destination Configuration Parameters

<!-- readmegen:destination.parameters.yaml -->
```yaml
version: 2.2
pipelines:
  - id: example
    status: running
    connectors:
      - id: example
        plugin: "chaos"
        settings:
          # ConfigureMode controls what the Configure method should do.
          # Type: string
          # Required: no
          configureMode: "success"
          # OpenMode controls what the Open method should do.
          # Type: string
          # Required: no
          openMode: "success"
          # TeardownMode controls what the Teardown method should do.
          # Type: string
          # Required: no
          teardownMode: "success"
          # WriteMode controls what the Write method should do.
          # Type: string
          # Required: no
          writeMode: "success"
          # Maximum delay before an incomplete batch is written to the
          # destination.
          # Type: duration
          # Required: no
          sdk.batch.delay: "0"
          # Maximum size of batch before it gets written to the destination.
          # Type: int
          # Required: no
          sdk.batch.size: "0"
          # Allow bursts of at most X records (0 or less means that bursts are
          # not limited). Only takes effect if a rate limit per second is set.
          # Note that if `sdk.batch.size` is bigger than `sdk.rate.burst`, the
          # effective batch size will be equal to `sdk.rate.burst`.
          # Type: int
          # Required: no
          sdk.rate.burst: "0"
          # Maximum number of records written per second (0 means no rate
          # limit).
          # Type: float
          # Required: no
          sdk.rate.perSecond: "0"
          # The format of the output record. See the Conduit documentation for a
          # full list of supported formats
          # (https://conduit.io/docs/using/connectors/configuration-parameters/output-format).
          # Type: string
          # Required: no
          sdk.record.format: "opencdc/json"
          # Options to configure the chosen output record format. Options are
          # normally key=value pairs separated with comma (e.g.
          # opt1=val2,opt2=val2), except for the `template` record format, where
          # options are a Go template.
          # Type: string
          # Required: no
          sdk.record.format.options: ""
          # Whether to extract and decode the record key with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.key.enabled: "true"
          # Whether to extract and decode the record payload with a schema.
          # Type: bool
          # Required: no
          sdk.schema.extract.payload.enabled: "true"
```
<!-- /readmegen:destination.parameters.yaml -->

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit tests.
