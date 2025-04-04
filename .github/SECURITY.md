# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions are:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

Please report security vulnerabilities to security@digitnxt.com.

You will receive a response from us within 48 hours. If the issue is confirmed, we will release a patch as soon as possible depending on complexity.

## Security Update Process

1. Dependabot automatically creates pull requests for security updates
2. Our team reviews these updates within 2 business days
3. Critical updates are fast-tracked and merged as soon as CI passes
4. Non-critical updates are batched and released weekly

## Security Best Practices

When contributing to this repository, please ensure:

1. All dependencies are from trusted sources
2. No sensitive information (API keys, credentials) is committed
3. All input is properly validated and sanitized
4. Security-related changes are reviewed by at least two maintainers

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine affected versions
2. Audit code to find any similar problems
3. Prepare fixes for all supported versions
4. Release new versions and notify users 