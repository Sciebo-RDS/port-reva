# Sciebo RDS to Reva connector
This connector allows Sciebo RDS to talk to Reva through the CS3APIs.
Implements this interface: https://github.com/Sciebo-RDS/Sciebo-RDS/blob/master/RDS/circle2_use_cases/interface_port_file_storage.yml

## Configuration
The connector can be configured via command-line and/or environment variables:

| Value | Command-line | Environment | Description | Mandatory |
| --- | --- | --- | --- | --- |
| Webserver port | `-port <uint>` | -- | The webserver port the connector will listen on | No (default=80) |
| Reva host | `-host <string>` | `RDS_REVA_HOST` | The Reva host address, including its gRPC port | Yes |
| Reva user | `-user <string>` | `RDS_REVA_USER` | The Reva user name | Yes |
| Reva user password | `-password <string>` | `RDS_REVA_PASSWORD` | The Reva user password | Yes |

Settings passed via command-line have precedence over environment variables. All mandatory settings have to be specified one way or the other for the connector to work.
