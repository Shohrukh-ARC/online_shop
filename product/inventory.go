package product

import "fmt"

func DisplayStock(product Product) {
	fmt.Println("Product:", product.Name)
	fmt.Printf("Remaining stock: %d\n", product.Stock)
}
