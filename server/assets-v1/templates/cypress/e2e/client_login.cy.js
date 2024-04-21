describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5001/client-login-page')
  })

  it('displays the login form', () => {
    cy.get('h1').should('contain', 'Login')
    cy.get('input#username').should('exist')
    cy.get('input#password').should('exist')
    cy.get('button').should('contain', 'Login')
  })

  it('allows users to login with valid credentials', () => {
    cy.get('input#username').type('hydrolife')
    cy.get('input#password').type('1234')
    cy.get('button').click()
    cy.url().should('include', 'http://localhost:5001/client-profile')
  })

  it('displays an error message for invalid credentials', () => {
    cy.get('input#username').type('invaliduser')
    cy.get('input#password').type('invalidpassword')
    cy.get('button').click()
    cy.on('window:alert', (message) => {
      expect(message).to.equal('Login failed!')
    })
  })

  it('redirects to registration page when "Register" link is clicked', () => {
    cy.contains('Register').click()
    cy.url().should('include', 'http://localhost:5001/client-register-page')
  })
})
