package demo.ProductsService.controller;

import demo.ProductsService.model.Products;
import demo.ProductsService.service.ProductAnalyticsService;
import demo.ProductsService.service.ProductService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/products")
public class ProductController {

    @Autowired
    private ProductService productService;

    @Autowired
    private ProductAnalyticsService productAnalyticsService;

    @GetMapping("/getAllProducts")
    public List<Products> getAllProducts() {
        return productService.getAllProducts();
    }

    @GetMapping("/getProductById/{id}")
    public ResponseEntity<Products> getProductById(@PathVariable Long id) {
        return ResponseEntity.ok(productService.getProductById(id));
    }

    @PostMapping("/createProduct")
    public Products createProduct(@RequestBody Products products) {
        return productService.createProduct(products);
    }

    @PutMapping("/updateProduct/{id}")
    public ResponseEntity<Products> updateProduct(@PathVariable Long id, @RequestBody Products productsDetails) {
        return ResponseEntity.ok(productService.updateProduct(id, productsDetails));
    }

    @DeleteMapping("/deleteProduct/{id}")
    public ResponseEntity<?> deleteProduct(@PathVariable Long id) {
        productService.deleteProduct(id);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/analytics/get-most-expensive")
    public ResponseEntity<Products> getMostExpensiveProduct() {
        Products products = productAnalyticsService.findMostExpensiveProduct();
        if (products != null) {
            return ResponseEntity.ok(products);
        } else {
            return ResponseEntity.noContent().build();
        }
    }

    @GetMapping("/analytics/get-average-price")
    public double getAveragePrice() {
        double averagePrice = productAnalyticsService.calculateAveragePrice();
        return averagePrice;

    }
}
