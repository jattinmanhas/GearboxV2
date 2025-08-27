# Validation Package

This package provides comprehensive validation for the Auth Service using `go-playground/validator/v10` with custom validation rules.

## Features

- **Struct-based validation** using Go struct tags
- **Custom validation rules** for business logic
- **Individual field validation** functions
- **User-friendly error messages**
- **Comprehensive validation rules** for common data types

## Quick Start

### 1. Basic Struct Validation

```go
type User struct {
    Username    string    `validate:"required,username"`
    Email       string    `validate:"required,email"`
    Password    string    `validate:"required,password"`
    DateOfBirth time.Time `validate:"required,date_of_birth"`
}

user := User{...}
if errors := validation.ValidateStruct(user); len(errors) > 0 {
    // Handle validation errors
    for _, err := range errors {
        fmt.Printf("Field: %s, Error: %s\n", err.Field, err.Message)
    }
}
```

### 2. Individual Field Validation

```go
if err := validation.ValidateUsername("john_doe"); err != nil {
    fmt.Println(err.Error())
}

if err := validation.ValidatePassword("SecurePass123"); err != nil {
    fmt.Println(err.Error())
}

if err := validation.ValidateEmail("user@example.com"); err != nil {
    fmt.Println(err.Error())
}
```

## Available Validation Tags

### Built-in Tags
- `required` - Field must be present and not empty
- `min=X` - Minimum length/value
- `max=X` - Maximum length/value
- `email` - Must be valid email format
- `url` - Must be valid URL format
- `oneof=val1 val2 val3` - Must be one of the specified values

### Custom Tags
- `username` - 3-30 characters, alphanumeric + underscore + hyphen
- `password` - At least 6 characters with uppercase, lowercase, and number
- `phone` - 10-15 digits (international format)
- `postal_code` - Valid postal code format
- `date_of_birth` - Valid date, not future, user must be 13+ years old

## Validation Rules

### Username
- Length: 3-30 characters
- Characters: Letters, numbers, underscores, hyphens only
- No spaces allowed

### Password
- Minimum length: 8 characters
- Must contain: uppercase letter, lowercase letter, number
- Special characters allowed but not required

### Email
- Standard email format validation
- Supports common email patterns

### Phone Number
- 10-15 digits
- Automatically strips non-digit characters
- International format support

### Date of Birth
- Cannot be in the future
- User must be at least 13 years old
- Cannot be more than 150 years ago

## Error Handling

Validation errors are returned as `ValidationErrors` which implements the `error` interface:

```go
type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Value   string `json:"value"`
    Message string `json:"message"`
}

type ValidationErrors []ValidationError
```

### Example Error Response

```json
{
    "field": "password",
    "tag": "password",
    "value": "weak",
    "message": "password must be at least 8 characters with uppercase, lowercase, and number"
}
```

## Best Practices

1. **Validate early**: Validate input data as soon as it's received
2. **Use struct tags**: Leverage struct tags for automatic validation
3. **Custom rules**: Create custom validation rules for business logic
4. **Error messages**: Provide clear, user-friendly error messages
5. **Performance**: Validation is fast, but avoid unnecessary validation calls

## Examples

See `examples.go` for comprehensive usage examples including:
- Basic user validation
- Business user validation with nested structs
- Error handling examples
- Custom validation scenarios

## Adding New Validation Rules

To add a new custom validation rule:

1. Create a validation function:
```go
func validateCustomField(fl validator.FieldLevel) bool {
    // Your validation logic here
    return true // or false if validation fails
}
```

2. Register it in the `init()` function:
```go
func init() {
    validate = validator.New()
    _ = validate.RegisterValidation("custom_field", validateCustomField)
}
```

3. Use it in struct tags:
```go
type MyStruct struct {
    Field string `validate:"required,custom_field"`
}
```
