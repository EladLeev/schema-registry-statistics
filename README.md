# Schema-registry-statistics
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/EladLeev/schema-registry-statistics/build.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/eladleev/schema-registry-statistics)](https://goreportcard.com/report/github.com/eladleev/schema-registry-statistics)  
Schema Registry Statistics Tool is a small utility that allows you to easily identify the usage of different schema versions within a topic.  
Using this tool, you can consume from a topic, while calculating the percentage of each schema version.  

Table of Contents
-----------------

- [schema-registry-statistics](#schema-registry-statistics)
  - [Table of Contents](#table-of-contents)
  - [Flags](#flags)
  - [Usage](#usage)
  - [How does it work?](#how-does-it-work)
  - [Local testing](#local-testing)
  - [License](#license)

Example output:
```bash
[sr-stats] 2022/12/28 10:02:12 Starting to consume from payments-topic
[sr-stats] 2022/12/28 10:02:12 Consumer up and running!...
[sr-stats] 2022/12/28 10:02:12 Use SIGINT to stop consuming.
[sr-stats] 2022/12/28 10:02:14 terminating: via signal
[sr-stats] 2022/12/28 10:02:14 Total messages consumed: 81
Schema ID 1 => 77%
Schema ID 3 => 23%
```
As you can see, in the `payments-topic`, 77% of the messages are produced using schema ID 1, while the remaining messages are produced using schema ID 3.

You can get the schema by ID:
```bash
curl -s http://<SCHEMA_REGISTRY_ADDR>/schemas/ids/1 | jq .
```

For further offsets analysis, you can store the results into a JSON file:
```json
{"1":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61],"3":[62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80]}
```

## Flags
| Name          | Description                                                                 | Require    | Type     | default             |
| ------------- | --------------------------------------------------------------------------- | ---------- | -------- | ------------------- |
| `--bootstrap` | The Kafka bootstrap servers.                                                | `V`        | `string` | "localhost:9092"    |
| `--topic`     | The topic name to consume.                                                  | `V`        | `string` | ""                  |
| `--version`   | The Kafka client version to be used.                                        |            | `string` | "2.1.1"             |
| `--group`     | The consumer group name.                                                    |            | `string` | schema-stats        |
| `--user`      | The Kafka username for authentication.                                      |            | `string` | ""                  |
| `--password`  | The Kafka password for authentication.                                      |            | `string` | ""                  |
| `--tls`       | Use TLS communication.                                                      |            | `bool`   | `false`             |
| `--cert`      | When TLS communication is enabled, specify the path for the CA certificate. | when `tls` | `string` | ""                  |
| `--store`     | Store results into a file.                                                  |            | `bool`   | `false`             |
| `--path`      | If `store` flag is set, the path to store the file.                         |            | `string` | "/tmp/results.json" |
| `--oldest`    | Consume from oldest offset.                                                 |            | `bool`   | `true`              |
| `--limit`     | Limit consumer to X messages, if different than 0.                          |            | `int`    | 0                   |
| `--verbose`   | Raise consumer log level.                                                   |            | `bool`   | `false`             |

## Usage
```bash
./schema-registry-statistics --bootstrap kafka1:9092 --group stat-consumer --topic payments-topic --store --path ~/results.json
```
Consume from `payments-topic` of `kafka1` and store the results. The consumer will run until `SIGINT` (`CMD + C`) will be used.

## How does it work?
According the Kafka [wire format](https://docs.confluent.io/platform/current/schema-registry/serdes-develop/index.html#wire-format), has only a couple of components:
| Bytes | Area       | Description                                                        |
| ----- | ---------- | ------------------------------------------------------------------ |
| 0     | Magic Byte | Confluent serialization format version number; currently always 0. |
| 1-4   | Schema ID  | 4-byte schema ID as returned by Schema Registry.                   |
| 5..   | Data       | Serialized data for the specified schema format (Avro, Protobuf).  |

The tool leverage this format, and reads the binary format of the each message in order to extract the schema ID and store it.

## Local testing
You can use the `docker-compose.yml` file to create a local environment from scratch.  
In the `/scripts` directory, there are 2 versions of the same schema, and a simple Python Avro producer.
## License
This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.
