# test-review

Review test coverage and quality

## Instructions

Review the test suite for completeness and quality.

1. **Analyze Tests**
   - Find all `*_test.go` files
   - Check test coverage with `go test -cover ./...`
   - Identify untested code paths

2. **Test Quality Checklist**
   - [ ] **Coverage**: Critical paths have tests?
   - [ ] **TDD**: Tests follow Red-Green-Refactor?
   - [ ] **Assertions**: Tests verify behavior, not implementation?
   - [ ] **Edge Cases**: Boundary conditions tested?
   - [ ] **Error Cases**: Failure scenarios covered?
   - [ ] **Mocking**: Dependencies properly mocked?
   - [ ] **Readability**: Test names describe what they test?
   - [ ] **Independence**: Tests don't depend on each other?

3. **Report**
   ```
   ## Test Review

   **Coverage**: XX%
   **Quality**: 🟢 Excellent | 🟡 Good | 🟠 Needs Work

   ### ✅ Well-Tested Areas
   - [list]

   ### ⚠️ Missing Tests
   - [file:function] - [what needs testing]

   ### 💡 Test Improvements
   - [suggestions]
   ```
