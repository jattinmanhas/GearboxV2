# API Setup for Registration Form

## Environment Configuration

Create a `.env.local` file in the frontend directory with the following configuration:

```bash
# Auth Service Configuration
AUTH_SERVICE_URL=http://localhost:8080

# Next.js Configuration
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

## Features Implemented

### 1. Form State Management
- Controlled form inputs with React state
- Real-time validation feedback
- Loading states during submission
- Error handling and display

### 2. Form Validation
- Required field validation
- Email format validation
- Password length validation (minimum 8 characters)
- Username length validation (minimum 3 characters)
- Real-time error clearing when user types

### 3. API Integration
- Next.js API routes that proxy requests to the auth service
- Centralized API utility (`lib/api.ts`) for maintainable code
- Proper error handling with custom `ApiError` class
- TypeScript interfaces for type safety

### 4. User Experience
- Loading states with disabled form during submission
- Success redirect to login page with success message
- Error messages displayed prominently
- Form validation with visual feedback (red borders on invalid fields)

## API Routes Created

### `/api/v1/auth/register`
- Proxies registration requests to the auth service
- Handles CORS and error responses
- Returns appropriate HTTP status codes

### `/api/v1/auth/login`
- Proxies login requests to the auth service
- Handles cookie management for authentication
- Ready for future login form implementation

## Usage

1. Start the auth service on port 8080
2. Start the frontend development server: `npm run dev`
3. Navigate to `/register` to test the form
4. Fill out the form and submit to create a new account

## Error Handling

The form handles various error scenarios:
- Network errors
- Validation errors from the server
- User input validation errors
- Server-side validation failures

All errors are displayed to the user with clear, actionable messages.
