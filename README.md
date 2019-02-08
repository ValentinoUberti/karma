# This project is a fork of Graphqlator CLI
__*This project is a WIP.*__
# 

Graphqlator takes your existing database schema and generates code for a GraphQL-Go server. Type 'graphqlator help' to see usage.

## Supported Data Stores:

- mysql
- mariadb
- mariadb
- postgres

## Installation:

```
go get github.com/ValentinoUberti/karma
```

## Prerequisites

[grahpql-go](https://github.com/graphql-go/graphql) - Generated code uses graphql-go

[GORM](https://github.com/jinzhu/gorm) - Generated code uses GORM

## Usage
```
  graphqlator [flags]
  graphqlator [command]
```
Available Commands:
```
  init        Create a graphqlator-pkg.json file.
  describe    Describe database or table
  generate    Generate GraphQL-Go API implementation using grapqhlator-pkg.json.
  help        Help about any command
  version     Print the version number of Graphqlator
```
Flags:
  -h, --help   help for graphqlator

Use "graphqlator [command] --help" for more information about a command.

## Example Usage:

Please visit the [graphqlator website](https://ahmedalhulaibi.github.io/graphqlator-website/) for a short tutorial.

# External Libraries Used
[goreturns](https://github.com/sqs/goreturns) - Generator uses goreturns to remove unnecessary generated imports

[Substance](https://github.com/ahmedalhulaibi/substance) - This library is used to introspect on the database information schema and generate the graphql-go code.

[grahpql-go](https://github.com/graphql-go/graphql) - The generated code is using this implementation of GraphQL in Go.

[GORM](https://github.com/jinzhu/gorm) - The generated code is using GORM.

