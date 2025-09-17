import { toast } from "sonner";

export interface NotificationOptions {
  title?: string;
  description?: string;
  duration?: number;
  action?: {
    label: string;
    onClick: () => void;
  };
}

/**
 * Show a success notification
 */
export function showSuccess(message: string, options?: NotificationOptions) {
  toast.success(message, {
    description: options?.description,
    duration: options?.duration || 4000,
    action: options?.action,
  });
}

/**
 * Show an error notification
 */
export function showError(message: string, options?: NotificationOptions) {
  toast.error(message, {
    description: options?.description,
    duration: options?.duration || 6000,
    action: options?.action,
  });
}

/**
 * Show an info notification
 */
export function showInfo(message: string, options?: NotificationOptions) {
  toast.info(message, {
    description: options?.description,
    duration: options?.duration || 4000,
    action: options?.action,
  });
}

/**
 * Show a warning notification
 */
export function showWarning(message: string, options?: NotificationOptions) {
  toast.warning(message, {
    description: options?.description,
    duration: options?.duration || 5000,
    action: options?.action,
  });
}

/**
 * Show a loading notification
 */
export function showLoading(message: string, options?: NotificationOptions) {
  return toast.loading(message, {
    description: options?.description,
  });
}

/**
 * Update a loading notification
 */
export function updateLoading(
  toastId: string | number,
  message: string,
  type: 'success' | 'error' | 'info' | 'warning' = 'success',
  options?: NotificationOptions
) {
  toast.dismiss(toastId);
  
  switch (type) {
    case 'success':
      return showSuccess(message, options);
    case 'error':
      return showError(message, options);
    case 'info':
      return showInfo(message, options);
    case 'warning':
      return showWarning(message, options);
    default:
      return showSuccess(message, options);
  }
}

/**
 * Dismiss a specific notification
 */
export function dismissNotification(toastId: string | number) {
  toast.dismiss(toastId);
}

/**
 * Dismiss all notifications
 */
export function dismissAllNotifications() {
  toast.dismiss();
}

// Predefined notification messages for common operations
export const NotificationMessages = {
  // Blog operations
  blog: {
    postCreated: "Blog post created successfully",
    postUpdated: "Blog post updated successfully", 
    postDeleted: "Blog post deleted successfully",
    postCreateError: "Failed to create blog post",
    postUpdateError: "Failed to update blog post",
    postDeleteError: "Failed to delete blog post",
    
    categoryCreated: "Category created successfully",
    categoryUpdated: "Category updated successfully",
    categoryDeleted: "Category deleted successfully",
    categoryCreateError: "Failed to create category",
    categoryUpdateError: "Failed to update category", 
    categoryDeleteError: "Failed to delete category",
  },
  
  // Product operations
  product: {
    created: "Product created successfully",
    updated: "Product updated successfully",
    deleted: "Product deleted successfully",
    createError: "Failed to create product",
    updateError: "Failed to update product",
    deleteError: "Failed to delete product",
  },
  
  // User operations
  user: {
    profileUpdated: "Profile updated successfully",
    profileUpdateError: "Failed to update profile",
    passwordChanged: "Password changed successfully",
    passwordChangeError: "Failed to change password",
  },
  
  // Cart operations
  cart: {
    itemAdded: "Item added to cart",
    itemRemoved: "Item removed from cart",
    cartCleared: "Cart cleared",
    addError: "Failed to add item to cart",
    removeError: "Failed to remove item from cart",
  },
  
  // Wishlist operations
  wishlist: {
    itemAdded: "Item added to wishlist",
    itemRemoved: "Item removed from wishlist",
    addError: "Failed to add item to wishlist",
    removeError: "Failed to remove item from wishlist",
  },
  
  // General operations
  general: {
    loading: "Loading...",
    success: "Operation completed successfully",
    error: "An error occurred",
    networkError: "Network error. Please check your connection.",
    unauthorized: "You are not authorized to perform this action",
    notFound: "Resource not found",
  }
} as const;
