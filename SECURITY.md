# Security Policy

## Security Updates

- Security updates will be released as soon as possible
- All security fixes will be released as a new minor version
- Critical vulnerabilities will be addressed within 48 hours

## Reporting a Vulnerability

We take security seriously. Please follow these steps to report a vulnerability:

1. **DO NOT** create a public GitHub issue for security vulnerabilities
2. Email security concerns to [nanodns@mguptahub.com]
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Your recommendation for fixing (if any)

## What to Expect

After reporting a vulnerability:

1. **Initial Response**: Within 48 hours
2. **Status Update**: Within 5 business days
3. **Resolution Timeline**: Based on severity
   - Critical: 48 hours
   - High: 1 week
   - Medium: 2 weeks
   - Low: Next release

## Security Best Practices

When using NanoDNS:

1. **Port Configuration**
   - Avoid running on port 53 in production
   - Use non-privileged ports (>1024)
   - Restrict port access to necessary networks

2. **Access Control**
   - Run with minimal privileges
   - Use Docker's network isolation features
   - Limit exposure to trusted networks only

3. **Record Configuration**
   - Validate DNS records before deployment
   - Use appropriate TTL values
   - Monitor for unexpected record changes

4. **Monitoring**
   - Monitor DNS query logs
   - Watch for unusual traffic patterns
   - Set up alerts for configuration changes

## Known Security Considerations

1. **Service Resolution**
   - Service discovery is limited to Docker network
   - External network access should be restricted
   - Use firewall rules when exposing the service

2. **Environment Variables**
   - Sensitive data in environment variables
   - Secure your environment file
   - Use Docker secrets when possible

## Version Verification

Verify the authenticity of releases:

1. **Docker Images**
   ```bash
   # Check image digest
   docker pull ghcr.io/mguptahub/nanodns:latest
   docker image inspect ghcr.io/mguptahub/nanodns:latest
   ```

2. **Binary Releases**
   - All releases are signed
   - Verify signatures using GPG
   ```bash
   # Example verification
   gpg --verify nanodns_linux_amd64.sig nanodns_linux_amd64
   ```



## Community Security Guidelines

### For Contributors

1. **Code Submissions**
   - No hardcoded credentials or secrets
   - Use environment variables for configuration
   - Follow secure coding practices:
     ```go
     // Do not expose sensitive info in logs
     log.Printf("Processing request from %s", sanitizeInput(source))
     
     // Use strong random number generation
     crypto/rand instead of math/rand
     
     // Validate all user inputs
     validateDNSRecord(record)
     ```

2. **Pull Request Security Checklist**
   - [ ] No sensitive information in code/comments
   - [ ] Input validation for new features
   - [ ] Error handling follows security best practices
   - [ ] Dependencies are from trusted sources
   - [ ] New features don't compromise existing security
   - [ ] Tests don't expose sensitive information

3. **Code Review Guidelines**
   - Check for potential security issues
   - Verify input validation
   - Review error handling
   - Examine logging practices
   - Validate configuration handling

4. **Documentation Contributions**
   - Don't include real domains/IPs in examples
   - Use example.com, example.net for demonstrations
   - Avoid exposing internal infrastructure details
   - Include security warnings where appropriate

### For Community Members

1. **Reporting Issues**
   - Use private reporting for security issues
   - Don't share exploit details publicly
   - Follow responsible disclosure
   - Wait for fixes before discussing publicly

2. **Discussing Security**
   - Use GitHub Security Advisories
   - Don't share vulnerability details in issues
   - Avoid posting sensitive configurations
   - Help others follow security best practices

3. **Testing and Feedback**
   - Report suspicious behavior
   - Test security fixes when requested
   - Provide feedback on security features
   - Share security enhancement ideas safely

### Security Best Practices for Development

1. **Local Development**
   ```bash
   # Use non-privileged ports
   export DNS_PORT=5353

   # Keep environment files secure
   chmod 600 .env
   
   # Use Docker's security features
   docker-compose up --build --force-recreate
   ```

2. **Testing Security Features**
   ```bash
   # Test with restricted permissions
   sudo -u nobody ./nanodns

   # Verify network isolation
   docker network inspect nanodns_network
   ```

3. **Code Analysis**
   ```bash
   # Run security linters
   gosec ./...
   
   # Check dependencies
   go mod verify
   govulncheck ./...
   ```

### Secure Integration Examples

1. **Docker Compose**
   ```yaml
   services:
     nanodns:
       image: ghcr.io/mguptahub/nanodns:latest
       security_opt:
         - no-new-privileges:true
       read_only: true
       environment:
         - DNS_PORT=5353
   ```

2. **Kubernetes**
   ```yaml
   securityContext:
     runAsNonRoot: true
     readOnlyRootFilesystem: true
     capabilities:
       drop:
         - ALL
   ```

### Security Communication Channels

1. **Official Channels**
   - GitHub Security Advisories
   - Security-related issues
   - Official releases

2. **Community Channels**
   - GitHub Discussions for general security topics
   - Release announcements for security updates
   - Documentation updates

### Security Enhancement Process

1. **Proposing Security Improvements**
   - Create a GitHub Discussion
   - Use security advisory if sensitive
   - Follow the security template
   - Wait for maintainer review

2. **Implementing Security Features**
   - Create a draft PR
   - Add tests for security features
   - Update documentation
   - Request security review

3. **Review Process**
   - Security-focused code review
   - Integration testing
   - Documentation review
   - Final security assessment

### Recognition Program

Contributors who help improve security can be recognized through:
1. Security acknowledgments in releases
2. Addition to CONTRIBUTORS.md
3. Special mention in security advisories
4. Community recognition badges


## Bug Bounty Program

Currently, we do not operate a bug bounty program. However, we deeply appreciate security researchers who:
1. Follow responsible disclosure
2. Provide detailed reports
3. Help improve NanoDNS security

## Acknowledgments

We maintain a list of security researchers who have helped improve NanoDNS security. Submit a PR to be added to this list.

## Updates to This Policy

This security policy may be updated. Check the commit history for changes.

Last updated: 2024-10-26