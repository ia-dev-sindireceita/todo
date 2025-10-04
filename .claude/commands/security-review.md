# security-review

Perform security-focused code review

## Instructions

You are a security expert reviewing code for vulnerabilities.

1. **Read Changes**
   - Run `git diff main` or read specified files

2. **Security Checklist**
   - [ ] **SQL Injection**: All queries use prepared statements?
   - [ ] **XSS**: Proper output encoding in templates?
   - [ ] **Authentication**: JWT properly validated?
   - [ ] **Authorization**: Access control checks in place?
   - [ ] **Sensitive Data**: No secrets in code? Passwords properly hashed?
   - [ ] **Input Validation**: All user input validated?
   - [ ] **CSRF**: Protection against cross-site requests?
   - [ ] **Session Management**: Secure cookie settings?
   - [ ] **Error Messages**: No sensitive info leaked in errors?
   - [ ] **Dependencies**: No known vulnerabilities?

3. **Report Format**
   ```
   ## Security Review

   **Risk Level**: 🔴 Critical | 🟠 High | 🟡 Medium | 🟢 Low

   ### 🔴 Critical Issues
   - [findings]

   ### 🟠 High Priority
   - [findings]

   ### 🟡 Medium Priority
   - [findings]

   ### ✅ Security Strengths
   - [good practices found]

   ### 📋 Recommendations
   - [actionable items]
   ```

4. **OWASP Top 10 Focus**
   - Check for OWASP Top 10 vulnerabilities
   - Provide specific remediation steps
