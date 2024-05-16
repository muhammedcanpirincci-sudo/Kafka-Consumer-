package demo.ProductsService.service;

import demo.ProductsService.model.Products;
import demo.ProductsService.repository.ProductRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.util.Comparator;
import java.util.List;

@Service
public class ProductAnalyticsService {

    @Autowired
    private ProductRepository productRepository;
    @Autowired
    private KafkaTemplate<String, String> kafkaTemplate;

    public double calculateAveragePrice() {
        List<Products> products = productRepository.findAll();
        double averagePrice = products.stream()
                // Convert BigDecimal to double
                .mapToDouble(product -> product.getUnitPrice().doubleValue())
                .average()
                .orElse(0.0);

        // Send a message to Kafka
        String avgPriceMessage = "Average price calculated: " + averagePrice;
        kafkaTemplate.send("products-topic", avgPriceMessage);

        return averagePrice;
    }


    public Products findMostExpensiveProduct() {
        Products mostExpensiveProducts = productRepository.findAll().stream()
                .max(Comparator.comparing(Products::getUnitPrice))
                .orElse(null);

        // Send a message to Kafka
        if (mostExpensiveProducts != null) {
            String expensiveProductMessage = "Most expensive product: " + mostExpensiveProducts.getProductName() +
                    " at price " + mostExpensiveProducts.getUnitPrice();
            kafkaTemplate.send("products-topic", expensiveProductMessage);
        }
        return mostExpensiveProducts;
    }

}
