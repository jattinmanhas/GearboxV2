import { httpClient } from "./http-client";
import type {
    Category,
    Product,
    CreateCategoryRequest,
    UpdateCategoryRequest,
    CreateProductRequest,
    UpdateProductRequest,
    ListProductsResponse,
    ListCategoriesResponse,
    ProductFilters,
    CategoryFilters,
    ProductVariant,
    CreateProductVariantRequest,
    UpdateProductVariantRequest,
    ApiResponse,
} from "../types";

export const productApi = {
    // Categories
    async getCategories(
        filters: CategoryFilters = {}
    ): Promise<ListCategoriesResponse> {
        const params = new URLSearchParams();

        if (filters.page) params.append("page", filters.page.toString());
        if (filters.limit) params.append("limit", filters.limit.toString());
        if (filters.search) params.append("search", filters.search);
        if (filters.parent_id)
            params.append("parent_id", filters.parent_id.toString());
        if (filters.is_active !== undefined)
            params.append("is_active", filters.is_active.toString());

        const query = params.toString();
        return httpClient.get<ListCategoriesResponse>(
            `/products/categories${query ? `?${query}` : ""}`
        );
    },

    async getCategory(id: number): Promise<Category> {
        return httpClient.get<Category>(`/products/categories/${id}`);
    },

    async createCategory(categoryData: CreateCategoryRequest): Promise<Category> {
        return httpClient.post<Category>("/products/categories", categoryData);
    },

    async updateCategory(
        id: number,
        categoryData: UpdateCategoryRequest
    ): Promise<Category> {
        return httpClient.put<Category>(
            `/products/categories/${id}`,
            categoryData
        );
    },

    async deleteCategory(id: number): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/products/categories/${id}`);
    },

    // Products
    async getProducts(
        filters: ProductFilters = {}
    ): Promise<ListProductsResponse> {
        const params = new URLSearchParams();

        if (filters.page) params.append("page", filters.page.toString());
        if (filters.limit) params.append("limit", filters.limit.toString());
        if (filters.search) params.append("search", filters.search);
        if (filters.category_id)
            params.append("category_id", filters.category_id.toString());
        if (filters.is_active !== undefined)
            params.append("is_active", filters.is_active.toString());
        if (filters.is_digital !== undefined)
            params.append("is_digital", filters.is_digital.toString());
        if (filters.min_price !== undefined)
            params.append("min_price", filters.min_price.toString());
        if (filters.max_price !== undefined)
            params.append("max_price", filters.max_price.toString());
        if (filters.in_stock !== undefined)
            params.append("in_stock", filters.in_stock.toString());
        if (filters.on_sale !== undefined)
            params.append("on_sale", filters.on_sale.toString());
        if (filters.tags && filters.tags.length > 0)
            params.append("tags", filters.tags.join(","));
        if (filters.sort_by) params.append("sort_by", filters.sort_by);
        if (filters.sort_order) params.append("sort_order", filters.sort_order);

        const query = params.toString();
        return httpClient.get<ListProductsResponse>(
            `/products${query ? `?${query}` : ""}`
        );
    },

    async getProduct(id: number): Promise<Product> {
        return httpClient.get<Product>(`/products/${id}`);
    },

    async createProduct(productData: CreateProductRequest): Promise<Product> {
        console.log("[API] Creating product with data:", productData);
        console.log("[API] Images being sent:", productData.images);
        return httpClient.post<Product>("/products", productData);
    },

    async updateProduct(
        id: number,
        productData: UpdateProductRequest
    ): Promise<Product> {
        return httpClient.put<Product>(`/products/${id}`, productData);
    },

    async deleteProduct(id: number): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/products/${id}`);
    },

    async searchProducts(query: string): Promise<Product[]> {
        const params = new URLSearchParams();
        params.append("q", query);

        const data = await httpClient.get<{ data: { products: Product[] } }>(
            `/products/search?${params}`
        );

        console.log("Search products response:", data);

        // Handle different response structures
        if (data.data?.products) {
            return data.data.products;
        } else if (Array.isArray(data.data)) {
            return data.data;
        } else if (Array.isArray(data)) {
            return data;
        }
        return [];
    },

    // Product Variants
    async getProductVariants(productId: number): Promise<ProductVariant[]> {
        const result = await httpClient.get<{ data: ProductVariant[] }>(
            `/products/${productId}/variants`
        );
        return result && result.data ? result.data : [];
    },

    async getProductVariantsWithInventory(
        productId: number
    ): Promise<ProductVariant[]> {
        const result = await httpClient.get<{ data: ProductVariant[] }>(
            `/products/${productId}/variants-with-inventory`
        );
        console.log("API response for variants:", result);
        console.log("Extracted data:", result.data);
        return result && result.data ? result.data : [];
    },

    async getProductVariant(variantId: number): Promise<ProductVariant> {
        const result = await httpClient.get<{ data: ProductVariant }>(
            `/products/variants/${variantId}`
        );
        return result.data;
    },

    async createProductVariant(
        productId: number,
        variantData: CreateProductVariantRequest
    ): Promise<ProductVariant> {
        const result = await httpClient.post<{ data: ProductVariant }>(
            `/products/${productId}/variants`,
            variantData
        );
        return result.data;
    },

    async updateProductVariant(
        variantId: number,
        variantData: UpdateProductVariantRequest
    ): Promise<ProductVariant> {
        const result = await httpClient.put<{ data: ProductVariant }>(
            `/products/variants/${variantId}`,
            variantData
        );
        return result.data;
    },

    async deleteProductVariant(variantId: number): Promise<ApiResponse> {
        return httpClient.delete<ApiResponse>(`/products/variants/${variantId}`);
    },
};
