package demo.ProductsService.service;

import demo.ProductsService.model.Products;
import demo.ProductsService.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ProductService {

    @Autowired
    private ProductRepository productRepository;

    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;

    public List<Products> getAllProducts() {
        kafkaTemplate.send("products-topic", "getAllProducts");
        return productRepository.findAll();
    }

    public Products getProductById(Long id) {
        return productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found"));
    }

    public Products createProduct(Products products) {
        // Save the product to the database
        Products savedProducts = productRepository.save(products);

        // Construct the message
        String message = "New Product Created: " + savedProducts.getProductName();

        // Send the message to Kafka
        kafkaTemplate.send("products-topic", message);

        return savedProducts;
    }


    public Products updateProduct(Long id, Products productsDetails) {
        Products products = productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found"));

        products.setProductName(productsDetails.getProductName());
        products.setCategoryID(productsDetails.getCategoryID());
        products.setUnitPrice(productsDetails.getUnitPrice());

        Products updatedProducts = productRepository.save(products);

        // Send a message to Kafka
        String updateMessage = "Product Updated: " + updatedProducts.getProductName();
        kafkaTemplate.send("products-topic", updateMessage);

        return updatedProducts;
    }


    public void deleteProduct(Long id) {
        Products products = productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found"));

        productRepository.delete(products);

        // Send a delete message to Kafka
        String deleteMessage = "Product Deleted: " + products.getProductName();
        kafkaTemplate.send("products-topic", deleteMessage);
    }

}
