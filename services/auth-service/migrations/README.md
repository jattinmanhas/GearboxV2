# Database Migrations

This directory contains database migrations for the Auth Service. Migrations are designed to be **idempotent** and **reversible**.

## Migration Files

### Current Migrations

| Migration | Description | Status |
|-----------|-------------|---------|
| `000001_initial_schema` | Core authentication tables (users, roles, refresh_tokens) | ✅ **Active** |

### Migration 000001: Initial Schema

**Tables Created:**
- `users` - User accounts and profiles
- `roles` - User roles and permissions
- `refresh_tokens` - JWT refresh token storage

**Features:**
- ✅ Default role assignment (role_id = 1 for new users)
- ✅ Automatic `updated_at` timestamps via triggers
- ✅ Proper foreign key constraints
- ✅ Optimized indexes for performance
- ✅ Soft delete support (`is_deleted` flags)
- ✅ **Optional fields**: `last_name`, `avatar`, `gender`, `date_of_birth`

**Default Roles:**
1. **User** (ID: 1) - Basic authenticated user
2. **Editor** (ID: 2) - Content editor with moderate permissions
3. **Admin** (ID: 3) - Full system administrator

**Optional User Fields:**
- `last_name` - Can be empty string or NULL
- `avatar` - Can be empty string or NULL  
- `gender` - Can be empty string, NULL, or valid values: 'male', 'female', 'other', 'prefer_not_to_say'
- `date_of_birth` - Can be NULL (zero time value)

## Migration Commands

### Using Go Migrate

```bash
# Install go-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" down 1

# Rollback all migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" down

# Check migration status
migrate -path migrations -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" version
```

### Using Docker

```bash
# Run migrations in Docker container
docker run --rm -v $(pwd)/migrations:/migrations --network host migrate/migrate \
  -path=/migrations \
  -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" up
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100), -- Optional
    avatar TEXT, -- Optional
    gender VARCHAR(20) CHECK (gender = '' OR gender IN ('male', 'female', 'other', 'prefer_not_to_say') OR gender IS NULL), -- Optional
    date_of_birth DATE, -- Optional
    role_id BIGINT DEFAULT 1 REFERENCES roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Roles Table
```sql
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Best Practices

### Writing Migrations

1. **Always use `IF NOT EXISTS`** for table creation
2. **Always use `IF EXISTS`** for table deletion
3. **Include proper rollback** in down migrations
4. **Use transactions** when possible
5. **Test both up and down migrations**
6. **Document complex migrations** with comments

### Migration Naming

- Use descriptive names: `000001_initial_schema.sql`
- Include version numbers for ordering
- Use underscores for separators
- Be specific about what the migration does

### Data Types

- **IDs**: Use `BIGSERIAL` for future scalability
- **Timestamps**: Use `TIMESTAMP WITH TIME ZONE` for consistency
- **Strings**: Use appropriate lengths (VARCHAR(50) vs VARCHAR(255))
- **Booleans**: Use `BOOLEAN` with `DEFAULT FALSE`

### Indexes

- Create indexes on frequently queried columns
- Use partial indexes for soft-deleted records
- Consider composite indexes for common query patterns
- Monitor index usage in production

## Troubleshooting

### Common Issues

1. **Migration already applied**: Check migration status with `migrate version`
2. **Foreign key constraints**: Ensure tables are created in correct order
3. **Permission errors**: Verify database user has necessary privileges
4. **Lock issues**: Check for long-running transactions

### Rollback Strategy

If a migration fails:
1. Check the error logs
2. Fix the issue in the migration file
3. Rollback to the previous version
4. Re-apply the corrected migration

### Emergency Rollback

```bash
# Force rollback (use with caution)
migrate -path migrations -database "postgres://..." force VERSION
```

## Future Migrations

When adding new migrations:

1. **Copy existing migration template**
2. **Update this README** with new migration details
3. **Test thoroughly** in development environment
4. **Document breaking changes** clearly
5. **Consider data migration** for existing records

## Security Considerations

- **Never commit sensitive data** in migrations
- **Use environment variables** for database credentials
- **Validate all inputs** in migration scripts
- **Test rollback scenarios** thoroughly
- **Backup database** before applying migrations in production
