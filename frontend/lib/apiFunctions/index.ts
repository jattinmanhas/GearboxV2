// Export utilities
export { ApiError } from "./api-error";
export { httpClient } from "./http-client";
export { API_BASE_URL } from "./constants";

// Export API modules
export { authApi } from "./auth.api";
export { productApi } from "./product.api";
export { profileApi } from "./profile.api";
export { inventoryApi } from "./inventory.api";
export { cartApi } from "./cart.api";
export { wishlistApi } from "./wishlist.api";
export { userApi } from "./user.api";
export { roleApi } from "./role.api";
export { oauthApi } from "./oauth.api";
export { orderApi } from "./order.api";
export { paymentApi } from "./payment.api";

// Re-export coupon API from existing file
export { couponApi } from "../coupon-api";
