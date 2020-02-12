# Nibbler User Group Extension

Provides a simple way of managing group-based authorization access.

Groups are collections of users.  A GroupMembership associates a user to a group.  GroupPrivilege records establish 
which actions a group can take on a specific resource (e.g. "Hospital #4").  If no resource is specified in a privilege, 
it is a global privilege (e.g. "create-customer").

## Required Nibbler extensions

Beyond what ships with Nibbler, these extensions are required:

- github.com/markdicksonjr/nibbler-sql

Currently, this extension can only use the SQL extension for persistence.  To expand support, a persistence 
encapsulation needs to be employed (as other Nibbler extensions have done).

## Running the sample app

The simplest run command uses an in-memory sqlite db:

```bash
PORT=3000 go run example/main.go
```

To use a file-based sqlite db:

```bash
PORT=3000 SQL_URL=./test.db go run example/main.go
```

## Playing with the sample app

Once running, open your browser and go to http://localhost:3000/api/noresource and confirm that authorization is not
granted.  Specifically, this API endpoint sh