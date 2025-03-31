## Guidelines for a Great PR:
- Submit small, focused, PRs that solve 1 issue
- Prefix each PR name with a related issue that it resolves (for example, `[ISSUE-123] Name of my great PR`)
- Name of the PR should be descriptive yet not too long
- Write a short high-level description of the changes made. For example:
  - Modified x-y-z on the frontend in xyz.html file
  - Created a new endpoint `xyz` on the backend which does 1-2-3 
  (outline exact steps from the technical perspective, for example: inserted a new record into the `users` table, 
  fetched data about cats from the remote API `https://somename.com`)
- Target to make the least possible number of commits:
  - When developing on your local, create 1 commit with initial implementation, 
  and amend all the subsequent changes to it via `git commit --amend`
  - Then push your branch (having a single commit with all your implementation) to the remote and create a PR
  - After that, you will only have to create a new commit in order to either address a PR revision, 
  or to resolve failures on the CI/CD pipeline
  - Example:
    - commit 1: initial implementation
    - commit 2: addressing first PR revision
    - commit 3: addressing second PR revision
    - commit 4: addressing failing Cypress tests on CI/CD
- The developer is responsible for manual testing before reviews are requested
- Automated tests should pass before reviews are requested
- Conflicts must be resolved before reviews are requested
- Always tag the code owner as the first reviewer
- Another, 2nd developer, should also be tagged for review
- If testing (unit, integration, E2E) is requested, testing is required
- Complete any Pull Request Template items or provide explicit justification for not
- Update any affected documentation
- The developer merges and deletes their own branch
- ! NEVER take feedback & change requests personally !

## Testing Guidelines

### E2E
Automate the "happy path" or the standard success case of your user story.

### Integration Tests
- Mock Sparingly
- Design for stubbing (i.e., use interfaces )
- API abuse?
 
### Structural Testing
Write at least one unit test for each branch of the function

### Bad Data
Include test cases that ensure the following classic parameter abuses don't crash the system

#### Strings
- Empty string
- Long string
- Special characters
- multi-line string
 
#### Numeric
- Zero value
- Min/Max value
- Positive / negative values
 
#### Collections
- Zero elements
- One element
- Nil value
- Duplicates
- Very large number of elements

#### Custom Structs
- Zero value
- Nil value
- Extra properties
- Abuses of sub-types


