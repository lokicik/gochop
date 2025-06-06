# GoChop Security Features

This document describes the security features implemented in the GoChop backend.

## JWT Authentication

The application includes JWT-based authentication to protect sensitive endpoints.

### Environment Variables

```env
JWT_SECRET=your-secret-key-here
BASE_URL=https://yourdomain.com
```

If `JWT_SECRET` is not set, a random secret will be generated (not recommended for production).

### Authentication Endpoints

- `POST /api/auth/login` - Login with username/password
- `GET /api/auth/dev-token` - Generate admin token (development only)
- `GET /api/profile` - Get user profile (requires JWT)

### Demo Credentials

For development/testing purposes, the following demo credentials are available:

**Admin User:**

- Username: `admin`
- Password: `admin123`

**Regular User:**

- Username: `user`
- Password: `user123`

### Protected Endpoints

The following endpoints now require JWT authentication:

- `GET /api/profile` - User profile (any authenticated user)
- `GET /api/admin/links` - View all links (admin only)
- `GET /api/admin/analytics/:shortCode` - View analytics (admin only)

### Usage Example

1. Login to get a token:

```bash
curl -X POST http://localhost:3001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

2. Use the token to access protected endpoints:

```bash
curl -X GET http://localhost:3001/api/admin/links \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## IP Filtering

The application supports IP whitelisting and blacklisting for enhanced security.

### Environment Variables

#### Global IP Filtering

```env
IP_FILTER_MODE=whitelist|blacklist|both
IP_WHITELIST=192.168.1.0/24,10.0.0.1,127.0.0.1
IP_BLACKLIST=192.168.100.0/24,203.0.113.0/24
```

#### Admin-Specific IP Filtering

```env
ADMIN_IP_FILTER_MODE=whitelist|blacklist|both
ADMIN_IP_WHITELIST=192.168.1.100,10.0.0.5
ADMIN_IP_BLACKLIST=
```

### Configuration Options

#### Filter Modes

- `whitelist` - Only allow specified IPs/ranges
- `blacklist` - Block specified IPs/ranges
- `both` - Apply both whitelist and blacklist (blacklist checked first)

#### IP Format Support

- Single IPs: `192.168.1.100`
- CIDR ranges: `192.168.1.0/24`
- Multiple entries: Comma-separated list

### IP Filtering Levels

1. **Global Filtering**: Applied to all endpoints
2. **Admin Filtering**: Additional filtering for admin endpoints only

### Example Configurations

#### Development (Allow localhost only for admin)

```env
ADMIN_IP_FILTER_MODE=whitelist
ADMIN_IP_WHITELIST=127.0.0.1,::1
```

#### Production (Block known bad IPs globally)

```env
IP_FILTER_MODE=blacklist
IP_BLACKLIST=192.168.100.0/24,203.0.113.0/24
ADMIN_IP_FILTER_MODE=whitelist
ADMIN_IP_WHITELIST=10.0.1.0/24
```

## Security Best Practices

### For Production Deployment

1. **Set a strong JWT secret:**

   ```env
   JWT_SECRET=your-very-long-random-secret-key-here
   ```

2. **Remove development endpoints:**

   - Remove or protect the `/api/auth/dev-token` endpoint

3. **Configure IP filtering:**

   - Whitelist known good IP ranges for admin access
   - Blacklist known malicious IP ranges

4. **Use HTTPS:**

   - Always use HTTPS in production
   - Configure proper TLS certificates

5. **Database security:**

   - Use strong database passwords
   - Restrict database access to application servers only

6. **Redis security:**
   - Configure Redis authentication
   - Restrict Redis access to application servers only

## Monitoring and Logging

The application includes analytics logging that tracks:

- IP addresses
- User agents
- Referrers
- Timestamps

This data can be used for:

- Security monitoring
- Abuse detection
- Performance analysis

## Future Enhancements

Consider implementing:

- Rate limiting per IP
- Account lockout after failed login attempts
- Two-factor authentication
- API key-based authentication for external integrations
- Audit logging for administrative actions
