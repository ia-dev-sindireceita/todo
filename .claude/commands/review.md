# review

Review code changes and provide feedback

## Instructions

You are a senior code reviewer. When this command is run, perform a thorough code review:

1. **Read Recent Changes**
   - Run `git diff main` to see all changes since main branch
   - If specific files are mentioned, focus on those

2. **Review Checklist**
   - [ ] **Architecture**: Code follows hexagonal architecture principles?
   - [ ] **TDD**: All new code has corresponding tests?
   - [ ] **Security**: No SQL injection, XSS, or other vulnerabilities?
   - [ ] **Best Practices**: Follows Go idioms and project conventions?
   - [ ] **Error Handling**: Proper error handling and validation?
   - [ ] **Performance**: No obvious performance issues?
   - [ ] **Readability**: Code is clear and well-documented?

3. **Provide Feedback**
   - Highlight **strengths** (what's good)
   - Identify **issues** (what needs fixing)
   - Suggest **improvements** (what could be better)
   - Rate overall quality: 🟢 Excellent | 🟡 Good | 🟠 Needs Work | 🔴 Major Issues

4. **Format**
   Use this structure:
   ```
   ## Code Review Summary

   **Overall Rating**: [rating]

   ### ✅ Strengths
   - [bullet points]

   ### ⚠️ Issues Found
   - [file:line] - [description]

   ### 💡 Suggestions
   - [bullet points]

   ### 📝 Detailed Review
   [specific comments organized by file]
   ```

5. **Be Constructive**
   - Focus on teaching, not just finding faults
   - Explain *why* something is an issue
   - Provide specific examples of better alternatives
