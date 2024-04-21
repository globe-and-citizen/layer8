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
  
    it('allows users to login with valid credentials', () => {
      cy.get('input[placeholder="Username"]').type('testuser')
      cy.get('input[placeholder="Password"]').type('password123')
      cy.get('button').click()
      cy.url().should('include', '/user')
    })
  
    it('displays an error message for invalid credentials', () => {
      cy.get('input[placeholder="Username"]').type('invaliduser')
      cy.get('input[placeholder="Password"]').type('invalidpassword')
      cy.get('button').click()
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Login failed!')
      })
    })
  
    it('redirects to registration page when "Register" link is clicked', () => {
      cy.contains('Register').click()
      cy.url().should('include', 'http://localhost:5001/user-register-page')
    })
  })
  