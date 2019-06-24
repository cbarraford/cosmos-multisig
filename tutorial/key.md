# Key

Start by creating a key.go file withing in the types folder. Within your key.go file, you will set your keys to be used throughout the creation of the module.

Defining the keys that will be used throughout the application helps with writing DRY code.

```go
package types

const (
	// module name
	ModuleName = "nameservice"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)
```

### Now we move on to the writing the [Keeper for the module](./keeper.md).
