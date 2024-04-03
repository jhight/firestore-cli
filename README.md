<p align="center">
    <img width="128" src="icon.png" align="center" alt="Stash" />
    <h1 align="center">firestore-cli</h1>
    <p align="center">A command-line utility for viewing and managing Firestore data.</p>
    <p><br/></p>
</p>

## Getting documents
Here are a few examples of how to retrieve data:
```bash
# list collections
firestore collections

# list documents in a collection
firestore get users

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

## Creating documents
Here are a few examples of creating document:
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
# set an entire document, specifying data manually
firestore set users/user-1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": false}'

# set an entire document, specifying data from a file
firestore set users/user-1234 <path/to/data.json
```

### Update a document
```bash
# update a single field in a document (others untouched)
firestore update users/user-1234 '{"age": 31}'

# update a few fields in a document (others untouched)
firestore update users/user-1234 '{"age": 31, "active": false}'
```

### Deleting a collection, document, or field
```bash
# delete a collection and all its documents (this will prompt you to confirm)
firestore delete users

# delete a document
firestore delete users/user-1234

# delete a field from a document
firestore delete users/user-1234 age
```

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