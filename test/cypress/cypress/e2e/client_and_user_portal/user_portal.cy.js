const username = 'testuser';
const firstname = 'John';
const lastname = 'Doe';
const displayname = 'JohnDoe';
const country = 'United States';
const password = 'password123';

describe('Authentication Page', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5001/user-register-page')
  })
  
  it('displays the registration form', () => {
    cy.get('h1').should('contain', 'Register')
    cy.get('input[placeholder="Username"]').should('exist')
    cy.get('input[placeholder="First Name"]').should('exist')
    cy.get('input[placeholder="Last Name"]').should('exist')
    cy.get('input[placeholder="Display Name"]').should('exist')
    cy.get('input[placeholder="Country"]').should('exist')
    cy.get('input[type="password"]').should('exist')
    cy.get('button').should('contain', 'Register')
  })

  it('allows users to register with valid data', () => {
    cy.get('input[placeholder="Username"]').type(username)
    cy.get('input[placeholder="First Name"]').type(firstname)
    cy.get('input[placeholder="Last Name"]').type(lastname)
    cy.get('input[placeholder="Display Name"]').type(displayname)
    cy.get('input[placeholder="Country"]').type(country)
    cy.get('input[type="password"]').type(password)
    cy.get('button').click()
    cy.url().should('include', 'http://localhost:5001/user-login-page')
  })

  it('displays an error message for incomplete registration data', () => {
    cy.get('button').click()
    cy.on('window:alert', (message) => {
      expect(message).to.equal('Please enter a username and password!')
    })
  })
})

describe('Authentication Page', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5001/user-login-page')
  })

  it('displays the login form', () => {
    cy.get('h1').should('contain', 'Login')
    cy.get('input[placeholder="Username"]').should('exist')
    cy.get('input[placeholder="Password"]').should('exist')
    cy.get('button').should('contain', 'Login')
  })

  it('displays an error message for invalid credentials', () => {
    cy.get('input[placeholder="Username"]').type('invaliduser')
    cy.get('input[placeholder="Password"]').type('invalidpassword')
    cy.get('button').click()

    cy.get('.bg-red-500').should('be.visible')
    cy.contains('.bg-red-500', 'Failed to get client profile').should('be.visible')
  })

  it('allows users to login with valid credentials', () => {
    cy.get('input[placeholder="Username"]').type(username)
    cy.get('input[placeholder="Password"]').type(password)
    cy.get('button').click()
    cy.url().should('include', '/user')

    cy.get('input[placeholder="Username"]').should('have.value', username);
    cy.get('input[placeholder="First Name"]').should('have.value', firstname);
    cy.get('input[placeholder="Last Name"]').should('have.value', lastname);
    cy.get('input[placeholder="Country"]').should('have.value', country);
  })

  it('allows users to update display name and verify email', () => {
    cy.get('input[placeholder="Username"]').type(username)
    cy.get('input[placeholder="Password"]').type(password)
    cy.get('button').click()
    cy.url().should('include', '/user')
    cy.get('input[placeholder="Display Name"]').type(displayname+'edited')
    cy.get('button').should('contain', 'Save change')
    cy.get('button').should('contain', 'Verify Email')
  })

  it('Logging out', () => {
    cy.get('div').contains('Log out').click();
  });
})

describe('Delete User', () => {
  const tableName = 'users';

  it('deletes the registered client', () => {
    cy.deleteRegisteredUser(username, tableName).then((result) => {
      expect(result).to.be.true; 
    });
  });
});