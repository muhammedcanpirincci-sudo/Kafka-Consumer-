package demo.ProductsService.model;

import jakarta.persistence.*;

import java.math.BigDecimal;

@Entity
@Table(name = "Products")
public class Products {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "ProductID")
    private int productID;

    @Column(name = "ProductName", nullable = false, length = 40)
    private String productName;

    @Column(name = "SupplierID")
    private int supplierID; // Use Integer for nullable integer columns

    @Column(name = "CategoryID")
    private int categoryID; // Use Integer for nullable integer columns

    @Column(name = "QuantityPerUnit", length = 20)
    private String quantityPerUnit;

    @Column(name = "UnitPrice", precision = 19, scale = 4)
    // For 'money' type, precision and scale need to be compatible
    private BigDecimal unitPrice;

    @Column(name = "UnitsInStock")
    private short unitsInStock; // Use Short for smallint type

    @Column(name = "UnitsOnOrder")
    private short unitsOnOrder; // Use Short for smallint type

    @Column(name = "ReorderLevel")
    private short reorderLevel; // Use Short for smallint type

    @Column(name = "Discontinued")
    private boolean discontinued; // Use Boolean for bit type


    public String getProductName() {
        return productName;
    }

    public void setProductName(String productName) {
        this.productName = productName;
    }

    public int getSupplierID() {
        return supplierID;
    }

    public void setSupplierID(int supplierID) {
        this.supplierID = supplierID;
    }

    public int getCategoryID() {
        return categoryID;
    }

    public void setCategoryID(int categoryID) {
        this.categoryID = categoryID;
    }

    public String getQuantityPerUnit() {
        return quantityPerUnit;
    }

    public void setQuantityPerUnit(String quantityPerUnit) {
        this.quantityPerUnit = quantityPerUnit;
    }

    public BigDecimal getUnitPrice() {
        return unitPrice;
    }

    public void setUnitPrice(BigDecimal unitPrice) {
        this.unitPrice = unitPrice;
    }

    public short getUnitsInStock() {
        return unitsInStock;
    }

    public void setUnitsInStock(short unitsInStock) {
        this.unitsInStock = unitsInStock;
    }

    public short getUnitsOnOrder() {
        return unitsOnOrder;
    }

    public void setUnitsOnOrder(short unitsOnOrder) {
        this.unitsOnOrder = unitsOnOrder;
    }

    public short getReorderLevel() {
        return reorderLevel;
    }

    public void setReorderLevel(short reorderLevel) {
        this.reorderLevel = reorderLevel;
    }

    public boolean isDiscontinued() {
        return discontinued;
    }

    public void setDiscontinued(boolean discontinued) {
        this.discontinued = discontinued;
    }
}
