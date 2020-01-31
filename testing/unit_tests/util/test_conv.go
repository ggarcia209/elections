// SUCCESS
package main

import (
	"fmt"

	"github.com/elections/util"
)

func main() {
	smallInt := 255
	mediumInt := 99999
	bigInt := 1000000000

	smallBytes := util.Itob(smallInt)
	mediumBytes := util.Itob(mediumInt)
	bigBytes := util.Itob(bigInt)

	fmt.Println("s: ", smallBytes)
	fmt.Println("m: ", mediumBytes)
	fmt.Println("b: ", bigBytes)
	fmt.Println()

	smallInt = util.Btoi(smallBytes)
	mediumInt = util.Btoi(mediumBytes)
	bigInt = util.Btoi(bigBytes)

	superInt := util.Btoi([]byte{255, 255, 255, 255})

	fmt.Println("s: ", smallInt)
	fmt.Println("m: ", mediumInt)
	fmt.Println("b: ", bigInt)
	fmt.Println("super: ", superInt)
	fmt.Println()
}
