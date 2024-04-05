<p align="center">
    <img width="128" src="icon.png" align="center" alt="Stash" />
    <h1 align="center">firestore-cli</h1>
    <p align="center">A command-line utility for viewing and managing Firestore data.</p>
    <p><br/></p>
</p>

## Retrieving data
```bash
# note: see firestore get --help for a lot more information
firestore get <path> [<field1>,<field2>,...] [--filter <json>] [--order <field>:<asc|desc>] [--limit <n>] [--offset <n>] [--count]
```
Here, `<path>` can be one of the following:

* a collection path, like `users`
* a document path, like `users/user-1234` or `users/user-1234/projects/project-5678`
* a subcollection path, like `users/user-1234/projects`

Additionally, you can specify a list of fields to display, a filter to apply, how to sort, and paging options.

If you want to reference a nested field, use the dot notation. For example, to get the `city` field from the `address` object, you would use `address.city`.

### Examples
Getting a document by its ID:
```bash
# get a user document
firestore get users/user-1234

# output:
{
  "address": {
    "city": "Chicago",
    "state": "Illinois",
    "zip": 60606
  },
  "age": 30,
  "firstName": "John",
  "lastName": "Doe",
  "interests": [
    "finance",
    "engineering"
  ]
}
```

Listing a collection, displaying only specific fields:
```bash
# show specific user fields ($id refers to document ID)
firestore get users \$id,first,lastName

# output:
[
  {
    "$id": "user-1234",
    "firstName": "John",
    "lastName": "Doe"
  },
  {
    "$id": "user-5678",
    "firstName": "Jane",
    "lastName": "Smith"
  }
]
```

Filtering documents:
```bash
# get users with a lastName of "Doe"
firestore get users firstName,lastName --filter '{"lastName":"Doe"}'

# output:
[
  {
    "firstName": "John",
    "lastName": "Doe"
  }
]
```

Filtering on nested properties:
```bash
# get users with a city of "Chicago"
firestore get users \$id,lastName,address --filter '{"address.city":"Chicago"}'

# output:
[
  {
    "$id": "user-1234",
    "address": {
      "city": "Chicago",
      "state": "Illinois",
      "zip": 60606
    },
    "lastName": "Doe"
  }
]
```

Here are a few more command examples:
```bash
# list documents in a collection
firestore get users

# list document IDs in a collection
firestore get users \$id

# get an entire user document by its ID
firestore get users/user-1234

# list the user's projects subcollection
firestore get users/user-1234/projects

# get specific fields from a user document
firestore get users/user-1234 name,age

# get specific fields from a user document, limit to first 10 documents
firestore get users/user-1234 name,age --limit 10

# get all users with a firstName of "John"
firestore get users --filter '{"firstName":"John"}'

# get all users with a firstName of "John" and a lastName of "Doe", ordered by age desc, then title asc
firestore get users --filter '{"firstName":"John","lastName":"Doe"}' -o age:desc,title:asc

# get all users with a firstName of "John" and a lastName of "Doe" or an age >= 30
firestore get users --filter '{"$or":{"$and":{"firstName":"John","lastName":"Doe"}},"age":{">=":30}}'

# get all users where address city (nested property) is one of: "New York", "Los Angeles", or "Chicago"
firestore get users --filter '{"address.city":{"$in":["New York","Los Angeles","Chicago"]}}'

# get a count of the users where address city is one of: "New York", "Los Angeles", or "Chicago"
firestore get users --filter '{"address.city":{"$in":["New York","Los Angeles","Chicago"]}}' --count
```

### Filter syntax
Let's look at one of the previous examples in more detail:
```bash
# get all users with a firstName of "John" and a lastName of "Doe" or an age >= 30
firestore get users --filter '{"$or":{"$and":{"firstName":"John","lastName":"Doe"},"age":{">=":30}}}'
```

Made more readable, the filter JSON is:
```json
{
  "$or": {
    "$and": {
      "firstName": "John",
      "lastName": "Doe"
    },
    "age": {
      ">=": 30
    }
  }
}
```
The `$and` and `$or` keys are logical operators. You can probably see that this translates into the following SQL pseudo-code:
```sql
... WHERE (firstName = "John" AND lastName = "Doe") OR age >= 30
```
## Listing collections
```bash
# note: see firestore collections --help for a lot more information
firestore collections [<path>]
```

### Examples
```bash
# list all collections
firestore collections

# list subcollections in a document
firestore collections users/user-1234
```

## Creating documents
```bash
# note: see firestore create --help for a lot more information
firestore create <path> <json>
```

### Examples
```bash
# create a document, specifying data manually
firestore create users/user-1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'

# create a document, specifying data from a file
firestore create users/user-1234 <path/to/data.json
```
Create will fail if the document already exists.

## Modifying documents
Modifying documents comes in two forms: `set` and `update`. The `set` command will overwrite the entire document, while `update` will only update the fields you specify.

### Set (e.g., create or replace) a document
```bash
# note: see firestore set --help for a lot more information
firestore set <path> <json>
```

#### Examples
```bash
# set an entire document, specifying data manually
firestore set users/user-1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": false}'

# set an entire document, specifying data from a file
firestore set users/user-1234/projects/project-1234 <path/to/data.json
```

### Update a document
```bash
# note: see firestore update --help for a lot more information
firestore update <path> <json>
```

#### Examples
```bash
# update a single field in a document (others untouched)
firestore update users/user-1234 '{"age": 31}'

# update a few fields in a document (others untouched)
firestore update users/user-1234/projects/project-5678 '{"active": false, "endDate": "2023-12-31"}'
```

## Deleting data
```bash
# note: see firestore delete --help for a lot more information
firestore delete <path> [<field1>,<field2>,...]
```
Firestore CLI can delete documents, fields, or entire collections. If `<path>` is a collection or subcollection path, Firestore CLI will prompt you to confirm before deleting the collection and all its documents. If one or more fields are specified, only those fields will be deleted from the document. Otherwise, the document referenced will be deleted.

### Examples
```bash
# delete a collection and all its documents (this will prompt you to confirm)
firestore delete users

# delete a document
firestore delete users/user-1234

# delete a field from a document
firestore delete users/user-1234 age
```

## Special tokens
<a name="special-tokens"></a>
### Filtering
| Token                 | Purpose                  | Example filter                                  |
|-----------------------|--------------------------|-------------------------------------------------|
| `$and`                | Logical AND              | `{"$and":{"k1":"v1","k2":"v2"}}`                |
| `$or`                 | Logical OR               | `{"$or":{"k1":"v1","k2":"v2"}}`                 |
| `>`                   | Greater than             | `{"field":{">":3}}`                             |
| `>=`                  | Greater than or equal to | `{"field":{">=":3}}`                            |
| `<`                   | Less than                | `{"field":{"<":3}}`                             |
| `<=`                  | Less than or equal to    | `{"field":{"<=":3}}`                            |
| `!=`                  | Not equal to             | `{"field":{"!=":3}}`                            |
| `$in`                 | In array                 | `{"field":{"$in":["v1","v2"]}}`                 |
| `$not-in`             | Not in array             | `{"field":{"$not-in":["v1","v2"]}}`             |
| `$array-contains`     | Array contains           | `{"field":{"$array-contains":"v1"}}`            |
| `$array-contains-any` | Array contains any       | `{"field":{"$array-contains-any":["v1","v2"]}}` |

### Values
| Token               | Represents                        | Example use                                                                        |
|---------------------|-----------------------------------|------------------------------------------------------------------------------------|
| `$id`               | Document ID                       | `firestore get users \$id`                                                         |
| `$now()`            | Current time function             | `firestore set users/user-1234 '{"lastUpdated":"$now()"}'`                         |
| `$timestamp(value)` | ISO-8601 timestamp parse function | `firestore set users/user-1234 {"lastUpdated":"$timestamp(2023-12-31T23:59:59Z)"}` |

## Configuration
You can move some of the boilerplate configuration out of CLI flags by storing them in a file. By default, Firestore CLI will look for `~/.firestore-cli.yaml`. You can specify a different configuration file with the `--config` flag.

Here's what the configuration file might look like:
```yaml
service-account: ~/path/to/service-account.json
project-id: your-project-id
pretty-print: true
spacing: 2
flatten: true
backup:
  collection: backup
  commands:
    - set
    - update
    - delete
```