# Merkle Tree

See my article on [Merkle trees].

# Example

```
package main

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/btoll/merkle-tree"
)

var farmhands = []string{
	"Huck",
	"Utley",
	"Molly",
	"Ben",
	"Pete",
	"Lily",
	"Rupert",
	"Moose",
	"Annie",
	"Phoebe",
	"Ginger",
	"Noam",
	"Chomsky",
}

func verifyFarmhands(tree *merkle.Tree) {
	for _, name := range farmhands {
		fmt.Fprintf(os.Stdout, "tree.VerifyProof(\"%s\") \t=> %t\n", name, tree.VerifyProof([]byte(name)))
	}
}

func main() {
	tree, err := merkle.New(sha256.New(), [][]byte{
		[]byte("Huck"),
		[]byte("Utley"),
		[]byte("Molly"),
		[]byte("Pete"),
	})

	err = tree.GenerateTree()
	if err != nil {
		fmt.Println(err)
	}

	verifyFarmhands(tree)
	fmt.Println("\nAdd more farmhands!\n")

	tree.AppendBlocks([][]byte{
		[]byte("Lily"),
		[]byte("Rupert"),
		[]byte("Moose"),
		[]byte("Annie"),
		[]byte("Phoebe"),
		[]byte("Noam"),
		[]byte("Chomsky"),
	})

	err = tree.GenerateTree()
	if err != nil {
		fmt.Println(err)
	}

	verifyFarmhands(tree)
}
```

## License

[GPLv3](COPYING)

## Author

Benjamin Toll

[Merkle trees]: http://www.benjamintoll.com/2018/08/08/on-merkle-trees/
[onrik/gomerkle]: https://github.com/onrik/gomerkle

