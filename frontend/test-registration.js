// Test script to demonstrate registration form functionality
// This simulates what happens when the form is submitted

const mockRegistrationData = {
  username: "test_user",
  password: "SecurePass123",
  email: "test@example.com",
  first_name: "John",
  middle_name: "M",
  last_name: "Doe"
};

// Simulate API call
async function testRegistration() {
  console.log("🚀 Testing Registration Form");
  console.log("📝 Form Data:", mockRegistrationData);
  
  try {
    // Simulate API call delay
    console.log("⏳ Sending registration request...");
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Simulate successful response
    const response = {
      success: true,
      message: "User registered successfully",
      data: {
        id: 1,
        username: mockRegistrationData.username,
        email: mockRegistrationData.email
      }
    };
    
    console.log("✅ Registration successful!");
    console.log("📊 Response:", response);
    console.log("🔄 Redirecting to login page...");
    
  } catch (error) {
    console.error("❌ Registration failed:", error.message);
  }
}

// Run the test
testRegistration();
