## Guidelins for a Great PR:
- Submit small, focused, PRs that solve 1 issue
- Summarize what you did and provide context and guidance
- The developer is responsible for manual testing before reviews are requested
- Automated tests should pass before reviews are requested
- Conflicts must be resolved before reviews are requested
- Always tag the code owner as the first reviewer
- Another, 2nd developer, should also be tagged for review
- Complete any Pull Request Template items or provide explicit justification for not
- If testing (unit, integration, E2E) is requested, testing is required
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


